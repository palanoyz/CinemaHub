package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type HealthHandler struct {
	db *mongo.Client
}

func NewHealthHandler(db *mongo.Client) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *gin.Context) {
	dbStatus := "UP"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx, nil); err != nil {
		dbStatus = "DOWN"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
		"components": gin.H{
			"database": dbStatus,
		},
		"message": "CinemaHub API is running",
	})
}
