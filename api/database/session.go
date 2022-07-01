package database

import (
	"context"
	"encoding/base32"
	"linkshare_api/utils"

	"github.com/gorilla/securecookie"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID       string             `bson:"_id,omitempty"`
	UserID   primitive.ObjectID `bson:"user_id"`
	Modified primitive.DateTime
	Schema   int
}

const Session_cookie_key string = "linkshare_session"

var base32RawStdEncoding *base32.Encoding = base32.StdEncoding.WithPadding(base32.NoPadding)

// the chances of this not being unique is incredibly small. 2^(32*8) ~= 10^77, so if we have 10k sessions, that's still 1 in 10^73.
func GetNewSessionID() string {
	return base32RawStdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
}

func NewSession() *Session {
	return &Session{
		ID:       GetNewSessionID(),
		Modified: utils.DateTimeNow(),
		Schema:   1,
	}
}

func (session *Session) Persist(ctx context.Context, insertOneSession utils.InsertOneFunc) error {
	_, err := insertOneSession(ctx, *session)
	return err
}

func (session *Session) Delete(ctx context.Context, deleteFunc utils.DeleteOneFunc) error {
	_, err := deleteFunc(ctx, *session)
	return err
}

func FindSessionFromID(ctx context.Context, sessionID string, findOneSession utils.FindOneFunc) (session *Session, err error) {
	session = &Session{}
	err = findOneSession(ctx, bson.M{"_id": sessionID}).Decode(session)
	return
}
