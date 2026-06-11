package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Booking struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        string        `bson:"user_id" json:"user_id"`
	ShowtimeID    bson.ObjectID `bson:"showtime_id" json:"showtime_id"`
	MovieTitle    string        `bson:"movie_title" json:"movie_title"`
	HallName      string        `bson:"hall_name" json:"hall_name"`
	StartTime     time.Time     `bson:"start_time" json:"start_time"`
	SeatIDs       []string      `bson:"seat_ids" json:"seat_ids"`
	TotalPrice    float64       `bson:"total_price" json:"total_price"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
}
