package websocket

import (
	"context"
	"log"
	"strings"

	"github.com/palanoyz/cinemahub/internal/model"
	"github.com/palanoyz/cinemahub/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func StartExpirationWorker(rdb *redis.Client, hub *Hub, showtimeRepo *repository.ShowtimeRepository) {
	ctx := context.Background()
	// Subscribe to the expiration event channel
	pubsub := rdb.PSubscribe(ctx, "__keyevent@0__:expired")

	log.Println("Redis Expiration Worker started")

	go func() {
		defer pubsub.Close()

		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Redis PubSub error: %v", err)
				continue
			}

			// Key format: lock:showtime:SHOWTIME_ID:seat:SEAT_ID
			key := msg.Payload
			if strings.HasPrefix(key, "lock:showtime:") {
				parts := strings.Split(key, ":")
				if len(parts) == 5 {
					showtimeIDStr := parts[2]
					seatID := parts[4]

					showtimeID, err := bson.ObjectIDFromHex(showtimeIDStr)
					if err != nil {
						continue
					}

					log.Printf("Auto-releasing expired seat %s for showtime %s", seatID, showtimeIDStr)

					// 1. Update MongoDB to AVAILABLE
					err = showtimeRepo.UpdateSeatStatus(ctx, showtimeID, []string{seatID}, model.SeatAvailable, "")
					if err != nil {
						log.Printf("Failed to auto-release seat in DB: %v", err)
						continue
					}

					// 2. Broadcast via WebSocket
					hub.Broadcast(SeatUpdate{
						ShowtimeID: showtimeIDStr,
						SeatIDs:    []string{seatID},
						Status:     string(model.SeatAvailable),
					})
				}
			}
		}
	}()
}
