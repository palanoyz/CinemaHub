package main

import (
	"cinemahub-backend/internal/auth"
	"cinemahub-backend/internal/config"
	"cinemahub-backend/internal/handler"
	"cinemahub-backend/internal/middleware"
	"cinemahub-backend/internal/repository"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func main() {
	db := config.ConnectMongoDB()
	rdb := config.ConnectRedis()
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
