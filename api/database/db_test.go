package database

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	err := godotenv.Load("../../.env", "../../.secrets")
	if err != nil {
		log.Fatalf("Failed to load .env file: %s", err)
	}
	rand.Seed(time.Now().UnixNano())
}

// quickly fail if there are connection/firewall issues
func TestPingDB(t *testing.T) {
	client, err := GetClient(context.TODO())
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

func TestCreatePage(t *testing.T) {
	// Get 10 free URL IDs. They should all be unique.
	pageURLs := map[string]bool{}
	total := time.Duration(0)
	runs := 10

	linksDB, err := NewLinkShareDB(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < runs; i++ {
		start := time.Now()
		page, err := linksDB.CreatePage(context.TODO(), "", primitive.NewObjectID(), linksDB.Pages.InsertOne)
		elapsed := time.Since(start)
		total += elapsed
		if err != nil {
			t.Fatalf("failed to create a new page: \n%s", err)
		}
		pageURL := page.URL
		if pageURLs[pageURL] {
			t.Fatalf("pageURL is not unique")
		}
		pageURLs[pageURL] = true
	}
	t.Logf("CreatePage took %s\n", total/time.Duration(runs))
}

func TestCreatePageNameTaken(t *testing.T) {

	linksDB, err := NewLinkShareDB(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	test_URL := "my_page"
	linksDB.Pages.DeleteOne(context.TODO(), bson.M{"_id": test_URL})
	page, err := linksDB.CreatePage(context.TODO(), test_URL, primitive.NewObjectID(), linksDB.Pages.InsertOne)
	if page.URL != test_URL {
		t.Fatal("created URL not same as input URL")
	}
	if err != nil {
		t.Fatalf("failed to create %s", test_URL)
	}
	_, err = linksDB.CreatePage(context.TODO(), test_URL, primitive.NewObjectID(), linksDB.Pages.InsertOne)
	_, isURLTakenError := err.(URLTakenError)
	if !isURLTakenError {
		t.Fatalf("expected URLTakenError, got: %s", err)
	}
}

func TestPageCreationLottery(t *testing.T) {
	linksDB, err := NewLinkShareDB(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	missCount := 0
	insertMock := func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (
		*mongo.InsertOneResult, error) {
		missCount++
		return nil, errors.New("E11000 duplicate key error")
	}
	_, err = linksDB.CreatePage(context.TODO(), "", primitive.NewObjectID(), insertMock)

	if missCount != 3 {
		t.Errorf("expected missCount to be 3, got: %d", missCount)
	}
	if !strings.Contains(err.Error(), "lottery") {
		t.Errorf("expected lottery message! Got: %s", err)
	}
}
