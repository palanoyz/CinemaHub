package repository

import (
	"github.com/palanoyz/cinemahub/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ShowtimeRepository struct {
	collection *mongo.Collection
}

func NewShowtimeRepository(db *mongo.Database) *ShowtimeRepository {
	return &ShowtimeRepository{
		collection: db.Collection("showtimes"),
	}
}

func (r *ShowtimeRepository) GetByMovieID(ctx context.Context, movieID bson.ObjectID) ([]model.Showtime, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"movie_id": movieID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var showtimes []model.Showtime
	if err := cursor.All(ctx, &showtimes); err != nil {
		return nil, err
	}
	if showtimes == nil {
		showtimes = []model.Showtime{}
	}
	return showtimes, nil
}

func (r *ShowtimeRepository) GetByID(ctx context.Context, id bson.ObjectID) (*model.Showtime, error) {
	var showtime model.Showtime
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&showtime)
	if err != nil {
		return nil, err
	}
	return &showtime, nil
}

func (r *ShowtimeRepository) Create(ctx context.Context, showtime *model.Showtime) error {
	res, err := r.collection.InsertOne(ctx, showtime)
	if err != nil {
		return err
	}
	showtime.ID = res.InsertedID.(bson.ObjectID)
	return nil
}
