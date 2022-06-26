package utils

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// taken from base64.encodeURL
const Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

var urlRegex *regexp.Regexp = regexp.MustCompile(`^[A-Za-z0-9_\-]{1,30}$`)

type Middleware func(http.Handler) http.Handler

func IsValidURL(URL string) bool {
	return urlRegex.MatchString(URL)
}

func DateTimeNow() primitive.DateTime {
	return primitive.NewDateTimeFromTime(time.Now())
}

func GetRandomURL(n int) string {
	// generate a random 6 character string
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(Charset[rand.Intn(len(Charset))])
	}
	return sb.String()
}

func GetRandomKeyString() string {
	value := securecookie.GenerateRandomKey(32)
	return base64.StdEncoding.EncodeToString(value)
}

// meant to be used on string we got from GetRandom32ByteB64EncodedString
func GetBytesFromKeyString(encodedString string) []byte {
	val, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		panic(err)
	}
	return val
}
