package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/palanoyz/cinemahub/internal/config"
	"github.com/palanoyz/cinemahub/internal/model"
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
	showtimesColl := db.Collection("showtimes")

	// Clean existing data
	moviesColl.DeleteMany(context.Background(), bson.M{})
	showtimesColl.DeleteMany(context.Background(), bson.M{})

	movies := []model.Movie{
		{
			Title:       "Your Name.",
			Description: "High schoolers Mitsuha and Taki are complete strangers living separate lives. But one night, they suddenly switch places. Mitsuha wakes up in Taki's body, and he in hers. This bizarre occurrence continues to happen randomly, and the two must adjust their lives around each other.",
			PosterURL:   "https://www.themoviedb.org/t/p/w1280/q719jXXEzOoYaps6babgKnONONX.jpg",
			Duration:    106,
			Genre:       []string{"Animation", "Romance", "Drama"},
			Rating:      8.5,
		},
		{
			Title:       "The Super Mario Galaxy Movie",
			Description: "Having thwarted Bowser's previous plot to marry Princess Peach, Mario and Luigi now face a fresh threat in Bowser Jr., who is determined to liberate his father from captivity and restore the family legacy. Alongside companions new and old, the brothers travel across the stars to stop the young heir's crusade.",
			PosterURL:   "https://www.themoviedb.org/t/p/w1280/eJGWx219ZcEMVQJhAgMiqo8tYY.jpg",
			Duration:    98,
			Genre:       []string{"Family", "Adventure", "Animation", "Comedy", "Fantasy"},
			Rating:      8.2,
		},
		{
			Title:       "Demon Slayer: Kimetsu no Yaiba Infinity Castle",
			Description: "The Demon Slayer Corps are drawn into the Infinity Castle, where Tanjiro, Nezuko, and the Hashira face terrifying Upper Rank demons in a desperate fight as the final battle against Muzan Kibutsuji begins.",
			PosterURL:   "https://www.themoviedb.org/t/p/w1280/fWVSwgjpT2D78VUh6X8UBd2rorW.jpg",
			Duration:    156,
			Genre:       []string{"Animation", "Action", "Fantasy"},
			Rating:      7.7,
		},
		{
			Title:       "Zootopia 2",
			Description: "After cracking the biggest case in Zootopia's history, rookie cops Judy Hopps and Nick Wilde find themselves on the twisting trail of a great mystery when Gary De'Snake arrives and turns the animal metropolis upside down. To crack the case, Judy and Nick must go undercover to unexpected new parts of town, where their growing partnership is tested like never before.",
			PosterURL:   "https://www.themoviedb.org/t/p/w1280/oJ7g2CifqpStmoYQyaLQgEU32qO.jpg",
			Duration:    108,
			Genre:       []string{"Animation", "Adventure", "Comedy", "Family", "Mystery"},
			Rating:      7.7,
		},
		{
			Title:       "Jujutsu Kaisen 0",
			Description: "Yuta Okkotsu is a nervous high school student who is suffering from a serious problem—his childhood friend Rika has turned into a curse and won't leave him alone. Since Rika is no ordinary curse, his plight is noticed by Satoru Gojo, a teacher at Jujutsu High, a school where fledgling exorcists learn how to combat curses. Gojo convinces Yuta to enroll, but can he learn enough in time to confront the curse that haunts him?",
			PosterURL:   "https://www.themoviedb.org/t/p/w1280/23oJaeBh0FDk2mQ2P240PU9Xxfh.jpg",
			Duration:    105,
			Genre:       []string{"Animation", "Adventure", "Comedy", "Family", "Mystery"},
			Rating:      8.1,
		},
	}

	for _, m := range movies {
		res, err := moviesColl.InsertOne(context.Background(), m)
		if err != nil {
			log.Fatal(err)
		}
		movieID := res.InsertedID.(bson.ObjectID)

		// Generate 2 showtimes for each movie
		for i := 1; i <= 2; i++ {
			showtime := model.Showtime{
				MovieID:   movieID,
				HallName:  fmt.Sprintf("Hall %d", i),
				StartTime: time.Date(2026, time.June, 12, 11, 30, 0, 0, time.Local),
				EndTime:   time.Date(2026, time.June, 12, 14, 0, 0, 0, time.Local),
				Seats:     generateSeats(6, 8),
			}
			_, err := showtimesColl.InsertOne(context.Background(), showtime)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Println("Successfully seeded movies and showtimes into MongoDB Atlas!")
}

func generateSeats(rows, cols int) []model.Seat {
	var seats []model.Seat
	for r := 0; rows > r; r++ {
		rowName := string(rune('A' + r))
		for c := 1; cols >= c; c++ {
			seats = append(seats, model.Seat{
				ID:     fmt.Sprintf("%s%d", rowName, c),
				Row:    rowName,
				Number: c,
				Status: model.SeatAvailable,
				Price:  250.0,
			})
		}
	}
	return seats
}
