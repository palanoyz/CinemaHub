package repository

import (
	"context"
	"github.com/palanoyz/cinemahub/internal/model"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MovieRepository struct {
	collection *mongo.Collection
}

func NewMovieRepository(db *mongo.Database) *MovieRepository {
	return &MovieRepository{
		collection: db.Collection("movies"),
	}
}

func (r *MovieRepository) GetAll(ctx context.Context) ([]model.Movie, error) {
	cursor, err := r.collection.Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []model.Movie
	if err := cursor.All(ctx, &movies); err != nil {
		return nil, err
	}
	// Return empty slice instead of nil for better JSON response
	if movies == nil {
		movies = []model.Movie{}
	}
	return movies, nil
}

func (r *MovieRepository) GetByID(ctx context.Context, id any) (*model.Movie, error) {
	var movie model.Movie
	err := r.collection.FindOne(ctx, map[string]any{"_id": id}).Decode(&movie)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}
