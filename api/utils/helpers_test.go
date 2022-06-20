package utils

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/joho/godotenv"
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
