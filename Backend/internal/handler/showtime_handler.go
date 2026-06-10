package handler

import (
	"github.com/palanoyz/cinemahub/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ShowtimeHandler struct {
	repo *repository.ShowtimeRepository
}

func NewShowtimeHandler(repo *repository.ShowtimeRepository) *ShowtimeHandler {
	return &ShowtimeHandler{repo: repo}
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

	c.JSON(http.StatusOK, showtime)
}
