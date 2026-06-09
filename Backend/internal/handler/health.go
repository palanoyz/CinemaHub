package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type HealthHandler struct {
	db    *mongo.Client
	redis *redis.Client
}

func NewHealthHandler(db *mongo.Client, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, redis: rdb}
}

func (h *HealthHandler) Check(c *gin.Context) {
	dbStatus := "UP"
	redisStatus := "UP"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx, nil); err != nil {
		dbStatus = "DOWN"
	}

	if _, err := h.redis.Ping(ctx).Result(); err != nil {
		redisStatus = "DOWN"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
		"components": gin.H{
			"database": dbStatus,
			"redis":    redisStatus,
		},
		"message": "CinemaHub API is running",
	})
}
