package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client

func ConnectMongoDB() *mongo.Client {
	uri := AppConfig.MongoDB_URI

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}

	log.Println("Connected to MongoDB successfully")
	MongoClient = client
	return client
}
