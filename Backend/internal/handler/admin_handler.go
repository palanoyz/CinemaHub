package handler

import (
	"github.com/palanoyz/cinemahub/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AdminHandler struct {
	bookingRepo *repository.BookingRepository
}

func NewAdminHandler(repo *repository.BookingRepository) *AdminHandler {
	return &AdminHandler{bookingRepo: repo}
}

func (h *AdminHandler) ListBookings(c *gin.Context) {
	// Build filters
	filter := bson.M{}
	
	movieTitle := c.Query("movie")
	if movieTitle != "" {
		filter["movie_title"] = bson.M{"$regex": movieTitle, "$options": "i"}
	}

	userID := c.Query("user_id")
	if userID != "" {
		filter["user_id"] = bson.M{"$regex": userID, "$options": "i"}
	}

	bookings, err := h.bookingRepo.GetAll(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}
