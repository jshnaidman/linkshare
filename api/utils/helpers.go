package utils

import (
	"math/rand"
	"regexp"
	"strings"
)

// taken from base64.encodeURL
const Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

var urlRegex *regexp.Regexp = regexp.MustCompile(`^[A-Za-z0-9_\-]{1,30}$`)

func IsValidURL(URL string) bool {
	return urlRegex.MatchString(URL)
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
