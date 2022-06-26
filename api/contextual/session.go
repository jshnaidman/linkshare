package contextual

import (
	"encoding/base32"
	"linkshare_api/utils"

	"github.com/gorilla/securecookie"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	Id       string
	UserId   primitive.ObjectID
	Modified primitive.DateTime
	Schema   int
}

const Session_cookie_key string = "linkshare_session"

var base32RawStdEncoding *base32.Encoding = base32.StdEncoding.WithPadding(base32.NoPadding)

// the chances of this not being unique is incredibly small. 2^(32*8) ~= 10^77, so if we have 10k sessions, that's still 1 in 10^73.
func GetNewSessionID() string {
	return base32RawStdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
}

func NewSessionForUser(userId primitive.ObjectID) *Session {
	return &Session{
		Id:       GetNewSessionID(),
		UserId:   userId,
		Modified: utils.DateTimeNow(),
	}
}
