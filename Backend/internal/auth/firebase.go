package auth

import (
	"cinemahub-backend/internal/config"
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var AuthClient *auth.Client

func InitFirebase() {
	projectID := config.AppConfig.FirebaseProjectID
	serviceAccountPath := config.AppConfig.FirebaseServiceAccountPath

	opt := option.WithAuthCredentialsFile(option.ServiceAccount, serviceAccountPath)
	config := &firebase.Config{ProjectID: projectID}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting firebase auth client: %v\n", err)
	}

	AuthClient = client
	log.Println("Firebase Admin SDK initialized successfully")
}
