package main

import (
	"github.com/palanoyz/cinemahub/internal/config"
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	config.LoadEnv()
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI not set")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("cinemahub")
	moviesColl := db.Collection("movies")

	// Clean existing movies
	moviesColl.DeleteMany(context.Background(), bson.M{})

	movies := []interface{}{
		bson.M{
			"title":       "The Super Mario Galaxy Movie",
			"description": "Having thwarted Bowser's previous plot to marry Princess Peach, Mario and Luigi now face a fresh threat in Bowser Jr., who is determined to liberate his father from captivity and restore the family legacy. Alongside companions new and old, the brothers travel across the stars to stop the young heir's crusade.",
			"poster_url":  "https://www.themoviedb.org/t/p/w1280/eJGWx219ZcEMVQJhAgMiqo8tYY.jpg",
			"duration":    98,
			"genre":       []string{"Family", "Adventure", "Animation", "Comedy", "Fantasy"},
			"rating":      8.2,
		},
		bson.M{
			"title":       "Demon Slayer: Kimetsu no Yaiba Infinity Castle",
			"description": "The Demon Slayer Corps are drawn into the Infinity Castle, where Tanjiro, Nezuko, and the Hashira face terrifying Upper Rank demons in a desperate fight as the final battle against Muzan Kibutsuji begins.",
			"poster_url":  "https://www.themoviedb.org/t/p/w1280/fWVSwgjpT2D78VUh6X8UBd2rorW.jpg",
			"duration":    156,
			"genre":       []string{"Animation", "Action", "Fantasy"},
			"rating":      7.7,
		},
		bson.M{
			"title":       "Zootopia 2",
			"description": "After cracking the biggest case in Zootopia's history, rookie cops Judy Hopps and Nick Wilde find themselves on the twisting trail of a great mystery when Gary De'Snake arrives and turns the animal metropolis upside down. To crack the case, Judy and Nick must go undercover to unexpected new parts of town, where their growing partnership is tested like never before.",
			"poster_url":  "https://www.themoviedb.org/t/p/w1280/oJ7g2CifqpStmoYQyaLQgEU32qO.jpg",
			"duration":    108,
			"genre":       []string{"Animation", "Adventure", "Comedy", "Family", "Mystery"},
			"rating":      7.7,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = moviesColl.InsertMany(ctx, movies)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully seeded movies into MongoDB Atlas!")
}
