package auth

import (
	"context"
	"linkshare_api/contextual"
	"linkshare_api/database"
	"linkshare_api/graph/model"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

var testFirstname = "first"
var testLastname = "last"
var testEmail = "test@gmail.com"
var testGoogleID = "abc123"

var testUser *model.User

func init() {
	err := godotenv.Load("../../.env", "../../.secrets")
	if err != nil {
		log.Fatalf("Failed to load .env file: %s", err)
	}
	rand.Seed(time.Now().UnixNano())

	testUser = &model.User{
		FirstName: &testFirstname,
		LastName:  &testLastname,
		Email:     &testEmail,
		GoogleID:  &testGoogleID,
		PageURLs:  []string{},
		Schema:    1,
	}
}

func cleanupTestUser(t *testing.T) func() {
	return func() {
		db, err := database.NewLinkShareDB(context.TODO())
		if err != nil {
			t.Fatal()
		}
		defer db.Disconnect(context.TODO())
		_, err = db.Users.DeleteOne(context.TODO(), bson.M{"email": testUser.Email})
		if err != nil {
			t.Fatalf("failed to cleanup testuser: %s", err)
		}
	}
}

// func cleanupTestSession(t *testing.T, session *contextual.Session) func() {
// 	return func() {
// 		db, err := database.NewLinkShareDB(context.TODO())
// 		if err != nil {
// 			t.Fatal()
// 		}
// 		defer db.Disconnect(context.TODO())
// 		_, err = db.Sessions.DeleteOne(context.TODO(), bson.M{"_id": session.ID})
// 		if err != nil {
// 			t.Fatalf("failed to cleanup session: %s", err)
// 		}
// 	}
// }

func upsertTestUserByGoogleID(t *testing.T) (*model.User, error) {
	db, err := database.NewLinkShareDB(context.TODO())
	if err != nil {
		t.Fatal()
	}
	defer db.Disconnect(context.TODO())
	return db.UpsertUserByGoogleID(context.TODO(), testUser, db.Users.FindOneAndUpdate)
}

func TestHandleJWTLogin(t *testing.T) {
	db, err := database.NewLinkShareDB(context.TODO())
	if err != nil {
		t.Fatal()
	}
	defer db.Disconnect(context.TODO())
	respWriter := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Authorization", "Bearer ASDASDASDA")
	handleJWTLogin(respWriter, req, func(bearerToken string, db *database.LinkShareDB, w http.ResponseWriter,
		r *http.Request) (user *model.User, err error) {
		user, err = upsertTestUserByGoogleID(t)
		t.Cleanup(cleanupTestUser(t))
		return
	})

	res := respWriter.Result()

	if res.StatusCode != http.StatusSeeOther {
		t.Errorf("Expected http code %d, got http code: %d", http.StatusSeeOther, res.StatusCode)
	}

	sessionID := ""
	for _, cookie := range res.Cookies() {
		if cookie.Name == contextual.Session_cookie_key {
			sessionID = cookie.Value
		}
	}
	if len(sessionID) == 0 {
		t.Errorf("Cookie not found in response")
	}

	session := &contextual.Session{}
	err = db.Sessions.FindOneAndDelete(context.TODO(), bson.M{"_id": sessionID}).Decode(session)
	if (err != nil) || (session == nil) {
		t.Errorf("Finding / deleting session error: %s", err)
	}
	user := &model.User{}
	err = db.Users.FindOne(context.TODO(), bson.M{"_id": session.UserID}).Decode(user)
	if (err != nil) || (user == nil) {
		t.Fatalf("Failed to find user created: %s", err)
	}

	if (*user.FirstName != *testUser.FirstName) || user.Username != nil {
		t.Errorf("testUser not correct: %#v", user)
	}

	if !time.Now().After(session.Modified.Time()) {
		t.Errorf("Modified time doesn't make sense")
	}

}
