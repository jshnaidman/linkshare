package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
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

// generate a random n character string
func GetRandomURL(n int) string {
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

func MarshalObjectID(objectID primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(objectID.String()))
	})
}

func UnmarshalObjectID(v interface{}) (primitive.ObjectID, error) {
	objIDStr, ok := v.(string)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("%T is not a string", v)
	}
	return primitive.ObjectIDFromHex(objIDStr)
}
