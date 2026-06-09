package main

import (
	"bufio"
	"cinemahub-backend/internal/config"
	"cinemahub-backend/internal/handler"
	"log"
	"os"
	"strings"

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

	router := gin.Default()

	healthHandler := handler.NewHealthHandler(db, rdb)
	router.GET("/health", healthHandler.Check)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
