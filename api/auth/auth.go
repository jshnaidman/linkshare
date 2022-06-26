package auth

import (
	"context"
	"errors"
	"linkshare_api/contextual"
	"linkshare_api/database"
	"linkshare_api/graph/model"
	"linkshare_api/utils"
	"net/http"
	"regexp"
	"time"

	"google.golang.org/api/idtoken"
)

// TODO: Testing
func ValidateGoogleJWT(ctx context.Context, JWT string) (payload *idtoken.Payload, err error) {
	validator, err := idtoken.NewValidator(ctx)
	if err != nil {
		return
	}
	payload, err = validator.Validate(ctx, JWT, utils.GetConf().GoogleClientID)
	if err != nil {
		return
	}
	// validator doesn't validate issuer
	// validating issuer is recommended: https://developers.google.com/identity/gsi/web/guides/verify-google-id-token
	if payload.Issuer != "accounts.google.com" && payload.Issuer != "https://accounts.google.com" {
		return nil, errors.New("invalid Issuer")
	}
	return
}

func NewSessionCookie(sessionID string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     contextual.Session_cookie_key,
		Value:    sessionID,
		MaxAge:   utils.GetConf().SessionLifetimeSeconds,
		HttpOnly: true,
		Secure:   utils.GetConf().IsProduction,
		SameSite: http.SameSiteLaxMode,
	}
}

// // Middleware decodes the share session cookie and packs the session into context
// func AuthMiddleware() utils.Middleware {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// Allow unauthenticated users in
// 			if err != nil || c == nil {
// 				next.ServeHTTP(w, r)
// 				return
// 			}

// 			userId, err := validateAndGetUserID(c)
// 			if err != nil {
// 				http.Error(w, "Invalid cookie", http.StatusForbidden)
// 				return
// 			}

// 			// get the user from the database
// 			user := getUserByID(db, userId)

// 			// put it in context
// 			ctx := context.WithValue(r.Context(), userCtxKey, user)

// 			// and call the next with our new context
// 			r = r.WithContext(ctx)
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

var bearerTokenRegex *regexp.Regexp = regexp.MustCompile(`Bearer ([a-zA-Z0-9\-_]+)$`)

// JWT format for google credential
// header
// {
//   "alg": "RS256",
//   "kid": "f05415b13acb9590f70df862765c655f5a7a019e", // JWT signature
//   "typ": "JWT"
// }
// payload
// {
//   "iss": "https://accounts.google.com", // The JWT's issuer
//   "nbf":  161803398874,
//   "aud": "314159265-pi.apps.googleusercontent.com", // Your server's client ID
//   "sub": "3141592653589793238", // The unique ID of the user's Google Account
//   "hd": "gmail.com", // If present, the host domain of the user's GSuite email address
//   "email": "elisa.g.beckett@gmail.com", // The user's email address
//   "email_verified": true, // true, if Google has verified the email address
//   "azp": "314159265-pi.apps.googleusercontent.com",
//   "name": "Elisa Beckett",
//                             // If present, a URL to user's profile picture
//   "picture": "https://lh3.googleusercontent.com/a-/e2718281828459045235360uler",
//   "given_name": "Elisa",
//   "family_name": "Beckett",
//   "iat": 1596474000, // Unix timestamp of the assertion's creation time
//   "exp": 1596477600, // Unix timestamp of the assertion's expiration time
//   "jti": "abc161803398874def"
// }

func LoginJWTHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO right now we will redirect to playground, but later we will redirect to user page
		loginRedirect := "http://localhost:8080/"
		// If the user is logged in and tries to access the login api, then we just redirect to the homepage.
		user := contextual.UserForContext(r.Context())
		if user != nil {
			http.Redirect(w, r, loginRedirect, http.StatusFound)
			return
		}

		// grab the bearer token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.LogDebug("No auth header in login request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		regMatch := bearerTokenRegex.FindStringSubmatch(authHeader)
		var bearerToken string
		var payload *idtoken.Payload
		var err error
		if len(regMatch) > 1 {
			bearerToken = regMatch[1]
			payload, err = ValidateGoogleJWT(r.Context(), bearerToken)
		} else {
			err = errors.New("loginJWTHandler - regex didn't match")
		}

		if err != nil {
			utils.LogDebug("loginJWTHandler - failed to validate jwt for request: %s.\n Auth header: \n%s", err, authHeader)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Create user object from JWT to update / create user
		email := payload.Claims["email"].(string)
		firstName := payload.Claims["given_name"].(string)
		lastName := payload.Claims["family_name"].(string)
		id := payload.Claims["sub"].(string)

		// user will choose username later if the user profile hasn't been created yet.
		user = &model.User{
			FirstName: &firstName,
			LastName:  &lastName,
			GoogleID:  &id,
			Email:     &email,
			Schema:    utils.GetConf().SchemaVersion,
		}

		db, err := database.NewLinkShareDB(r.Context())
		defer db.Disconnect(r.Context())
		if err != nil {
			utils.LogError("loginJWTHandler - Failed to retrieve db: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		user, err = db.UpsertUserByGoogleID(r.Context(), user, db.Users.FindOneAndUpdate)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.LogError("loginJWTHandler - Failed to update user: %s\n%#v", err, user)
			return
		}

		session := contextual.NewSessionForUser(user.Id)
		err = db.CreateSession(r.Context(), session, db.Sessions.InsertOne)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.LogError("loginJWTHandler - Failed to create session: %s", err)
			return
		}
		http.SetCookie(w, NewSessionCookie(session.Id, session.Modified.Time()))
		http.Redirect(w, r, loginRedirect, http.StatusSeeOther)
	})
}
