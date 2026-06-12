package auth

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/palanoyz/cinemahub/internal/config"
	"google.golang.org/api/option"
)

var AuthClient *auth.Client

func InitFirebase() {
	projectID := config.AppConfig.FirebaseProjectID
	serviceAccountJSON := config.AppConfig.FirebaseServiceAccountJSON

	if projectID == "" || serviceAccountJSON == "" {
		log.Fatal("FIREBASE_PROJECT_ID and FIREBASE_SERVICE_ACCOUNT_JSON must be set in .env")
	}

	opt := option.WithAuthCredentialsJSON(option.ServiceAccount, []byte(serviceAccountJSON))

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
	log.Println("Firebase Admin SDK initialized successfully using JSON string")
}
