package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/palanoyz/cinemahub/internal/model"
	"github.com/palanoyz/cinemahub/internal/repository"
	"github.com/palanoyz/cinemahub/internal/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookingHandler struct {
	showtimeRepo *repository.ShowtimeRepository
	bookingRepo  *repository.BookingRepository
	movieRepo    *repository.MovieRepository
	redis        *redis.Client
	hub          *websocket.Hub
	rabbit       *amqp.Connection
}

func NewBookingHandler(
	showtimeRepo *repository.ShowtimeRepository,
	bookingRepo *repository.BookingRepository,
	movieRepo *repository.MovieRepository,
	rdb *redis.Client,
	hub *websocket.Hub,
	rabbit *amqp.Connection,
) *BookingHandler {
	return &BookingHandler{
		showtimeRepo: showtimeRepo,
		bookingRepo:  bookingRepo,
		movieRepo:    movieRepo,
		redis:        rdb,
		hub:          hub,
		rabbit:       rabbit,
	}
}

type BookingRequest struct {
	SeatIDs []string `json:"seat_ids" binding:"required"`
}

func (h *BookingHandler) LockSeats(c *gin.Context) {
	showtimeIDStr := c.Param("id")
	showtimeID, err := bson.ObjectIDFromHex(showtimeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid showtime ID"})
		return
	}

	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	ctx := c.Request.Context()

	// 1. Try to acquire/refresh Redis Distributed Locks
	for _, seatID := range req.SeatIDs {
		lockKey := fmt.Sprintf("lock:showtime:%s:seat:%s", showtimeIDStr, seatID)

		// SET NX EX
		success, err := h.redis.SetNX(ctx, lockKey, userID, 5*time.Minute).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Locking system error"})
			return
		}

		if !success {
			currentOwner, _ := h.redis.Get(ctx, lockKey).Result()
			if currentOwner != userID {
				c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Seat %s is already being selected by someone else", seatID)})
				return
			}
			h.redis.Expire(ctx, lockKey, 5*time.Minute)
		}
	}

	// 2. Update MongoDB status to LOCKED
	err = h.showtimeRepo.UpdateSeatStatus(ctx, showtimeID, req.SeatIDs, model.SeatLocked, userID)
	if err != nil {
		// Rollback Redis locks if MongoDB fails
		for _, seatID := range req.SeatIDs {
			lockKey := fmt.Sprintf("lock:showtime:%s:seat:%s", showtimeIDStr, seatID)
			h.redis.Del(ctx, lockKey)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seat status"})
		return
	}

	// 3. Broadcast to all clients
	h.hub.Broadcast(websocket.SeatUpdate{
		ShowtimeID: showtimeIDStr,
		SeatIDs:    req.SeatIDs,
		Status:     string(model.SeatLocked),
		LockedBy:   userID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Seats locked successfully"})
}

func (h *BookingHandler) UnlockSeats(c *gin.Context) {
	showtimeIDStr := c.Param("id")
	showtimeID, err := bson.ObjectIDFromHex(showtimeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid showtime ID"})
		return
	}

	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	ctx := c.Request.Context()

	// 1. Verify ownership and remove from Redis
	for _, seatID := range req.SeatIDs {
		lockKey := fmt.Sprintf("lock:showtime:%s:seat:%s", showtimeIDStr, seatID)

		owner, err := h.redis.Get(ctx, lockKey).Result()
		if err == nil && owner == userID {
			h.redis.Del(ctx, lockKey)
		}
	}

	// 2. Update MongoDB status to AVAILABLE
	err = h.showtimeRepo.UpdateSeatStatus(ctx, showtimeID, req.SeatIDs, model.SeatAvailable, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to release seats"})
		return
	}

	// 3. Broadcast
	h.hub.Broadcast(websocket.SeatUpdate{
		ShowtimeID: showtimeIDStr,
		SeatIDs:    req.SeatIDs,
		Status:     string(model.SeatAvailable),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Seats released successfully"})
}

func (h *BookingHandler) ConfirmBooking(c *gin.Context) {
	showtimeIDStr := c.Param("id")
	showtimeID, err := bson.ObjectIDFromHex(showtimeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid showtime ID"})
		return
	}

	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	ctx := c.Request.Context()

	// 1. Verify that user still owns the locks in Redis
	for _, seatID := range req.SeatIDs {
		lockKey := fmt.Sprintf("lock:showtime:%s:seat:%s", showtimeIDStr, seatID)
		owner, _ := h.redis.Get(ctx, lockKey).Result()
		if owner != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Seat %s lock has expired or belongs to someone else", seatID)})
			return
		}
	}

	// 2. Update MongoDB status to BOOKED
	err = h.showtimeRepo.UpdateSeatStatus(ctx, showtimeID, req.SeatIDs, model.SeatBooked, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm booking"})
		return
	}

	// Get Showtime and Movie info for the booking record
	showtime, _ := h.showtimeRepo.GetByID(ctx, showtimeID)
	movie, _ := h.movieRepo.GetByID(ctx, showtime.MovieID)

	// Calculate price (for now fixed at 250 per seat)
	totalPrice := float64(len(req.SeatIDs)) * 250.0

	// Save Booking Record
	booking := &model.Booking{
		UserID:     userID,
		ShowtimeID: showtimeID,
		MovieTitle: movie.Title,
		HallName:   showtime.HallName,
		StartTime:  showtime.StartTime,
		SeatIDs:    req.SeatIDs,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now(),
	}
	h.bookingRepo.Create(ctx, booking)

	// 3. Clear Redis Locks
	for _, seatID := range req.SeatIDs {
		lockKey := fmt.Sprintf("lock:showtime:%s:seat:%s", showtimeIDStr, seatID)
		h.redis.Del(ctx, lockKey)
	}

	// 4. Broadcast via WebSockets
	h.hub.Broadcast(websocket.SeatUpdate{
		ShowtimeID: showtimeIDStr,
		SeatIDs:    req.SeatIDs,
		Status:     string(model.SeatBooked),
	})

	// 5. ASYNC: Publish to RabbitMQ
	go h.publishToRabbit(userID, showtimeIDStr, req.SeatIDs)

	c.JSON(http.StatusOK, gin.H{"message": "Booking confirmed! Your tickets are being generated."})
}

func (h *BookingHandler) ListUserBookings(c *gin.Context) {
	userID := c.GetString("user_id")

	bookings, err := h.bookingRepo.GetAll(c.Request.Context(), bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) publishToRabbit(userID, showtimeID string, seatIDs []string) {
	ch, err := h.rabbit.Channel()
	if err != nil {
		return
	}
	defer ch.Close()

	msg := websocket.BookingMessage{
		UserID:     userID,
		ShowtimeID: showtimeID,
		SeatIDs:    seatIDs,
	}

	body, _ := json.Marshal(msg)

	ch.Publish(
		"",         // exchange
		"bookings", // routing key (queue name)
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
