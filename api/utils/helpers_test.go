package utils

import (
	"bytes"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	err := godotenv.Load("../../.env", "../../.secrets")
	if err != nil {
		log.Fatalf("Failed to load .env file: %s", err)
	}
	rand.Seed(time.Now().UnixNano())
}

func TestIsValidURL(t *testing.T) {
	validURLs := []string{"a", GetRandomURL(30), "-_aAZ09"}
	invalidURLs := []string{"ab/c", "😡", "もしもし", "", GetRandomURL(31)}

	for _, url := range validURLs {
		if !IsValidURL(url) {
			t.Errorf("Expected string to pass: %s", url)
		}
	}
	for _, url := range invalidURLs {
		if IsValidURL(url) {
			t.Errorf("Expected string to fail: %s", url)
		}
	}
}

func TestMarshalObjectID(t *testing.T) {
	id := primitive.NewObjectIDFromTimestamp(time.Now())
	marshaler := MarshalObjectID(id)
	buffer := bytes.NewBuffer([]byte{})
	marshaler.MarshalGQL(buffer)

	unmarshaled, err := UnmarshalObjectID(buffer.String())
	if err != nil {
		t.Fatalf("Failed to unmarshal")
	}
	if unmarshaled != id {
		t.Errorf("Unmarshaeled %s not equal to id %s", id.Hex(), unmarshaled.Hex())
	}
}
