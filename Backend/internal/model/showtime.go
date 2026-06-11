package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SeatStatus string

const (
	SeatAvailable SeatStatus = "AVAILABLE"
	SeatLocked    SeatStatus = "LOCKED"
	SeatBooked    SeatStatus = "BOOKED"
)

type Seat struct {
	ID       string     `bson:"id" json:"id"`         // e.g., "A1", "A2"
	Row      string     `bson:"row" json:"row"`       // e.g., "A"
	Number   int        `bson:"number" json:"number"` // e.g., 1
	Status   SeatStatus `bson:"status" json:"status"`
	Price    float64    `bson:"price" json:"price"`
	LockedBy string     `bson:"locked_by,omitempty" json:"locked_by,omitempty"`
}

type Showtime struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	MovieID   bson.ObjectID `bson:"movie_id" json:"movie_id"`
	HallName  string        `bson:"hall_name" json:"hall_name"`
	StartTime time.Time     `bson:"start_time" json:"start_time"`
	EndTime   time.Time     `bson:"end_time" json:"end_time"`
	Seats     []Seat        `bson:"seats" json:"seats"`
}
