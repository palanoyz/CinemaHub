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

	if len(os.Args) < 3 {
		log.Fatal("Usage: go run cmd/admin/main.go <promote|demote> <user_id>")
	}

	action := os.Args[1]
	uid := os.Args[2]

	isAdmin := false
	if action == "promote" {
		isAdmin = true
	} else if action != "demote" {
		log.Fatal("Action must be 'promote' or 'demote'")
	}

	claims := map[string]any{"admin": isAdmin}
	err := auth.AuthClient.SetCustomUserClaims(context.Background(), uid, claims)
	if err != nil {
		log.Fatalf("error setting custom claims: %v\n", err)
	}

	log.Printf("Successfully set admin status to %v for user %s\n", isAdmin, uid)
}
