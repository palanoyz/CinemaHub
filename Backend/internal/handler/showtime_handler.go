package handler

import (
	"context"
	"fmt"
	"github.com/palanoyz/cinemahub/internal/model"
	"github.com/palanoyz/cinemahub/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ShowtimeHandler struct {
	repo  *repository.ShowtimeRepository
	redis *redis.Client
}

func NewShowtimeHandler(repo *repository.ShowtimeRepository, rdb *redis.Client) *ShowtimeHandler {
	return &ShowtimeHandler{
		repo:  repo,
		redis: rdb,
	}
}

func (h *ShowtimeHandler) GetByMovie(c *gin.Context) {
	movieIDStr := c.Param("id")
	movieID, err := bson.ObjectIDFromHex(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	showtimes, err := h.repo.GetByMovieID(c.Request.Context(), movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch showtimes"})
		return
	}

	c.JSON(http.StatusOK, showtimes)
}

func (h *ShowtimeHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid showtime ID"})
		return
	}

	showtime, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		return
	}

	// --- LAZY CLEANUP LOGIC ---
	// Check if any LOCKED seats have expired in Redis
	ctx := c.Request.Context()
	updated := false
	var expiredSeatIDs []string

	for i, seat := range showtime.Seats {
		if seat.Status == model.SeatLocked {
			lockKey := fmt.Sprintf("lock:showtime:%s:seat:%s", idStr, seat.ID)
			exists, _ := h.redis.Exists(ctx, lockKey).Result()
			
			if exists == 0 {
				// Lock expired in Redis, make it available in MongoDB
				showtime.Seats[i].Status = model.SeatAvailable
				showtime.Seats[i].LockedBy = ""
				expiredSeatIDs = append(expiredSeatIDs, seat.ID)
				updated = true
			}
		}
	}

	if updated {
		// Async update MongoDB so we don't block the response
		go h.repo.UpdateSeatStatus(context.Background(), id, expiredSeatIDs, model.SeatAvailable, "")
	}
	// ---------------------------

	c.JSON(http.StatusOK, showtime)
}
