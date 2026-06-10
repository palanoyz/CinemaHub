package main

import (
	"github.com/palanoyz/cinemahub/internal/auth"
	"github.com/palanoyz/cinemahub/internal/config"
	"github.com/palanoyz/cinemahub/internal/handler"
	"github.com/palanoyz/cinemahub/internal/middleware"
	"github.com/palanoyz/cinemahub/internal/repository"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	// Initialize Connections
	db := config.ConnectMongoDB()
	rdb := config.ConnectRedis()
	auth.InitFirebase()

	// Initialize Repositories
	movieRepo := repository.NewMovieRepository(db.Database("cinemahub"))
	showtimeRepo := repository.NewShowtimeRepository(db.Database("cinemahub"))

	// Initialize Handlers
	healthHandler := handler.NewHealthHandler(db, rdb)
	movieHandler := handler.NewMovieHandler(movieRepo)
	showtimeHandler := handler.NewShowtimeHandler(showtimeRepo)

	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.AppConfig.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	router.GET("/health", healthHandler.Check)

	api := router.Group("/api")
	{
		api.GET("/movies", movieHandler.List)
		api.GET("/movies/:id/showtimes", showtimeHandler.GetByMovie)
		api.GET("/showtimes/:id", showtimeHandler.GetByID)
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

	port := config.AppConfig.Port
	log.Println("Server starting on :" + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
