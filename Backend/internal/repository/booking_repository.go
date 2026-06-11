package repository

import (
	"context"
	"github.com/palanoyz/cinemahub/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BookingRepository struct {
	collection *mongo.Collection
}

func NewBookingRepository(db *mongo.Database) *BookingRepository {
	return &BookingRepository{
		collection: db.Collection("bookings"),
	}
}

func (r *BookingRepository) Create(ctx context.Context, booking *model.Booking) error {
	res, err := r.collection.InsertOne(ctx, booking)
	if err != nil {
		return err
	}
	booking.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *BookingRepository) GetAll(ctx context.Context, filter bson.M) ([]model.Booking, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []model.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	if bookings == nil {
		bookings = []model.Booking{}
	}
	return bookings, nil
}
