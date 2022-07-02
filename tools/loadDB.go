package main

import (
	"context"
	"linkshare_api/database"
	"linkshare_api/graph/model"
	"linkshare_api/utils"
	"log"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
)

func LoadUsersSessions() {
	db, err := database.NewLinkShareDB(context.TODO())
	defer db.Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}

	for j := 0; j < 3; j++ {
		users := make([]interface{}, 1e6)
		sessions := make([]interface{}, 1e6)

		for i := 0; i < 1e6; i++ {
			username := "loadUsername"
			firstName := "loadFirstName"
			lastName := "loadLastName"
			email := "loadEmail"
			var googleID string
			user := &model.User{
				Username:  &username,
				FirstName: &firstName,
				LastName:  &lastName,
				Email:     &email,
				GoogleID:  &googleID,
				PageURLs:  []string{},
				Schema:    1,
			}
			googleID = utils.GetRandomURL(10)
			user.GoogleID = &googleID
			username = utils.GetRandomURL(10)
			user.Username = &username
			session := database.NewSession()
			session.UserID = user.ID
			users[i] = *user
			sessions[i] = *session
		}
		_, err = db.Users.InsertMany(context.TODO(), users)
		if err != nil {
			panic(err)
		}
		_, err = db.Sessions.InsertMany(context.TODO(), sessions)
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	err := godotenv.Load("../.env", "../.secrets")
	if err != nil {
		log.Fatalf("Failed to load .env file: %s", err)
	}
	conf := utils.GetConf()
	conf.DBPort = "28017" // I use a different port externally because I don't want it to conflict with my local mongodb instance.
	conf.SetConnectionURL()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	LoadUsersSessions()
}
