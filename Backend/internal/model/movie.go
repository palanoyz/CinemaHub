package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Movie struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	PosterURL   string        `bson:"poster_url" json:"poster_url"`
	Duration    int           `bson:"duration" json:"duration"` // in minutes
	Genre       []string      `bson:"genre" json:"genre"`
	Rating      float64       `bson:"rating" json:"rating"`
}
