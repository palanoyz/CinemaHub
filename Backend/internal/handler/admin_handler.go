package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/palanoyz/cinemahub/internal/model"
	"github.com/palanoyz/cinemahub/internal/repository"
	"github.com/palanoyz/cinemahub/internal/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AdminHandler struct {
	bookingRepo  *repository.BookingRepository
	showtimeRepo *repository.ShowtimeRepository
	hub          *websocket.Hub
}

func NewAdminHandler(bookingRepo *repository.BookingRepository, showtimeRepo *repository.ShowtimeRepository, hub *websocket.Hub) *AdminHandler {
	return &AdminHandler{
		bookingRepo:  bookingRepo,
		showtimeRepo: showtimeRepo,
		hub:          hub,
	}
}

func (h *AdminHandler) ListBookings(c *gin.Context) {
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

func (h *AdminHandler) CancelBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	ctx := c.Request.Context()
	booking, err := h.bookingRepo.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// 1. Release seats in showtime
	err = h.showtimeRepo.UpdateSeatStatus(ctx, booking.ShowtimeID, booking.SeatIDs, model.SeatAvailable, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to release seats"})
		return
	}

	// 2. Delete booking
	err = h.bookingRepo.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}

	// 3. Broadcast
	h.hub.Broadcast(websocket.SeatUpdate{
		ShowtimeID: booking.ShowtimeID.Hex(),
		SeatIDs:    booking.SeatIDs,
		Status:     string(model.SeatAvailable),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
}
