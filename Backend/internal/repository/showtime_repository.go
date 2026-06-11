package repository

import (
	"context"
	"github.com/palanoyz/cinemahub/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (r *ShowtimeRepository) UpdateSeatStatus(ctx context.Context, showtimeID bson.ObjectID, seatIDs []string, status model.SeatStatus, userID string) error {
	filter := bson.M{
		"_id":      showtimeID,
		"seats.id": bson.M{"$in": seatIDs},
	}

	var update bson.M
	if status == model.SeatAvailable {
		// Clearing lock
		update = bson.M{
			"$set": bson.M{
				"seats.$[elem].status": status,
			},
			"$unset": bson.M{
				"seats.$[elem].locked_by": "",
			},
		}
	} else {
		// Setting lock/booking
		update = bson.M{
			"$set": bson.M{
				"seats.$[elem].status":    status,
				"seats.$[elem].locked_by": userID,
			},
		}
	}

	opts := options.UpdateMany().SetArrayFilters([]any{
		bson.M{"elem.id": bson.M{"$in": seatIDs}},
	})

	_, err := r.collection.UpdateMany(ctx, filter, update, opts)
	return err
}
