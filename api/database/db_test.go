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

// quickly fail if there are connection/firewall issues
func TestPingDB(t *testing.T) {
	client, err := GetClient()
	if err != nil {
		t.Fatalf("Client error: %s", err)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Fatalf("Client ping error: %s", err)
	}
}

func TestGetFreeURL(t *testing.T) {
	// Get 10 free URL IDs. They should all be unique.
	pageIDs := map[uint32]bool{}
	total := time.Duration(0)
	runs := 10
	for i := 0; i < runs; i++ {
		start := time.Now()
		pageID, err := createNewPage("")
		elapsed := time.Since(start)
		total += elapsed
		if err != nil {
			t.Fatalf("Failed to retrieve free URL ID: %s", err)
		}
		if pageIDs[pageID] {
			t.Fatalf("pageID is not unique")
		}
		pageIDs[pageID] = true
	}
	t.Logf("createNewPage took %s", total/time.Duration(runs))

}
