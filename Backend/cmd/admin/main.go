package main

import (
	"context"
	"log"
	"os"

	"github.com/palanoyz/cinemahub/internal/auth"
	"github.com/palanoyz/cinemahub/internal/config"
)

func main() {
	config.LoadEnv()
	auth.InitFirebase()

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/admin/main.go <user_id>")
	}

	uid := os.Args[1]

	claims := map[string]any{"admin": true}
	err := auth.AuthClient.SetCustomUserClaims(context.Background(), uid, claims)
	if err != nil {
		log.Fatalf("error setting custom claims: %v\n", err)
	}

	log.Printf("Successfully promoted user %s to ADMIN\n", uid)
}
