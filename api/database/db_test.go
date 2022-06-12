package database

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	err := godotenv.Load("../../.env", ".testenv")
	if err != nil {
		log.Fatalf("Failed to load .env file: %s", err)
	}
	// override dev environment variables with .testenv
	err = godotenv.Overload(".testenv")
	if err != nil {
		log.Fatalf("Failed to load .env file: %s", err)
	}
}

func TestPingDB(t *testing.T) {
	client, err := GetClient()
	if err != nil {
		t.Errorf("Client error: %s", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Errorf("Client ping error: %s", err)
	}
}
