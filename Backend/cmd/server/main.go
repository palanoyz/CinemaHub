package main

import (
	"bufio"
	"cinemahub-backend/internal/auth"
	"cinemahub-backend/internal/config"
	"cinemahub-backend/internal/handler"
	"cinemahub-backend/internal/middleware"
	"cinemahub-backend/internal/repository"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	// Simple .env loader
	file, err := os.Open(".env")
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
	}
}

func main() {
	// Initialize MongoDB
	db := config.ConnectMongoDB()
	// Initialize Redis
	rdb := config.ConnectRedis()
	// Initialize Firebase
	auth.InitFirebase()

	// Initialize Repositories
	movieRepo := repository.NewMovieRepository(db.Database("cinemahub"))

	// Initialize Handlers
	healthHandler := handler.NewHealthHandler(db, rdb)
	movieHandler := handler.NewMovieHandler(movieRepo)

	router := gin.Default()

	// CORS Middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	// Public routes
	router.GET("/health", healthHandler.Check)

	api := router.Group("/api")
	{
		api.GET("/movies", movieHandler.List)
	}

	// Protected routes
	protected := router.Group("/api/protected")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"message": "You are authenticated",
			})
		})
	}

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
