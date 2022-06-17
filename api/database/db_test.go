package database

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	err := godotenv.Load("../../.env", "../../.secrets", ".testenv")
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

func TestGetFreeURL(t *testing.T) {
	// Get 10 free URL IDs. They should all be unique.
	urlIDs := map[int32]bool{}
	for i := 0; i < 100; i++ {
		page, err := createNewPage("")
		if err != nil {
			t.Errorf("Failed to retrieve free URL ID: %s", err)
			return
		}
		urlID := page["_id"].(int32)
		fmt.Println(urlID)
		if urlIDs[urlID] {
			t.Error("urlID is not unique")
			return
		}
		urlIDs[urlID] = true
	}

}
