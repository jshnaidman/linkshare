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

// TODO right now we will redirect to playground, but later we will redirect to user page
var loginRedirect string = "http://localhost:8080/"

type loginValidationMethod func(bearerToken string, db *database.LinkShareDB, w http.ResponseWriter,
	r *http.Request) (user *model.User, err error)

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
		Name:     database.Session_cookie_key,
		Value:    sessionID,
		MaxAge:   utils.GetConf().SessionLifetimeSeconds,
		HttpOnly: true,
		Secure:   utils.GetConf().IsProduction,
		SameSite: http.SameSiteLaxMode,
	}
}

// Middleware decodes the share session cookie and packs the session into context
func AuthMiddleware() utils.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: It's important to cache the session and user to avoid doing 2 db queries on every request.
			db, err := database.NewLinkShareDB(r.Context())
			defer db.Disconnect(r.Context())
			if err != nil {
				utils.LogError("AuthMiddleware - Failed to retrieve db: %s", err)
				// database is down so just return 500
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			sessionIDCookie, err := r.Cookie(database.Session_cookie_key)

			// No session, create one and allow unauthenticated users in
			if err != nil || sessionIDCookie == nil || len(sessionIDCookie.Value) == 0 {
				session := database.NewSession()
				err = session.Persist(r.Context(), db.Sessions.InsertOne)
				if err != nil {
					utils.LogError("AuthMiddleware - Failed to create session: %s", err)
				}
				next.ServeHTTP(w, r)
				return
			}

			session, err := database.FindSessionFromID(r.Context(), sessionIDCookie.Value, db.Sessions.FindOne)
			if err != nil || session == nil {
				// failed to find session for this cookie, delete it and continue unauthenticated
				sessionIDCookie.MaxAge = -1
				http.SetCookie(w, sessionIDCookie)
				next.ServeHTTP(w, r)
				return
			}

			// get the user from the database
			user := &model.User{
				ID: session.UserID,
			}
			err = user.Update(r.Context(), db.Users.UpdateByID)
			if err != nil {
				// session exists but no user, continue unauthenticated.
				next.ServeHTTP(w, r)
				return
			}

			// put it in context
			ctx := contextual.AddUserToContext(r.Context(), user)
			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

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

func GoogleLoginHandleFunc(w http.ResponseWriter, r *http.Request) {
	handleJWTLogin(w, r, validateGoogleLogin)
}

// Receives JWT in auth header and creates authenticated session. Returns session object if successful
func handleJWTLogin(w http.ResponseWriter, r *http.Request, loginValidation loginValidationMethod) *database.Session {
	// If the user is logged in and tries to access the login api, then we just redirect to the homepage.
	user := contextual.UserForContext(r.Context())
	if user != nil {
		http.Redirect(w, r, loginRedirect, http.StatusFound)
		return nil
	}

	// grab the bearer token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.LogDebug("No auth header in login request")
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	regMatch := bearerTokenRegex.FindStringSubmatch(authHeader)
	var bearerToken string
	var err error
	if len(regMatch) == 2 {
		bearerToken = regMatch[1]
	} else {
		err = errors.New("loginJWTHandler - regex didn't match")
		utils.LogDebug("loginJWTHandler - failed to validate jwt for request: %s.\n Auth header: \n%s", err, authHeader)
		w.WriteHeader(http.StatusUnauthorized)
	}

	db, err := database.NewLinkShareDB(r.Context())
	defer db.Disconnect(r.Context())
	if err != nil {
		utils.LogError("loginJWTHandler - Failed to retrieve db: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	user, err = loginValidation(bearerToken, db, w, r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.LogError("loginJWTHandler - Failed to update user: %s\n%#v", err, user)
		return nil
	}

	// get session from context
	session := contextual.SessionForContext(r.Context())
	if session == nil {
		// if there is no session because there was some kind of error creating session we create one
		session = database.NewSession()
		session.UserID = user.ID
	}
	err = session.Persist(r.Context(), db.Sessions.InsertOne)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.LogError("loginJWTHandler - Failed to create session: %s", err)
		return nil
	}
	http.SetCookie(w, NewSessionCookie(session.ID, session.Modified.Time()))
	http.Redirect(w, r, loginRedirect, http.StatusSeeOther)
	return session
}

func validateGoogleLogin(bearerToken string, db *database.LinkShareDB, w http.ResponseWriter, r *http.Request) (user *model.User, err error) {
	payload, err := ValidateGoogleJWT(r.Context(), bearerToken)
	if err != nil {
		utils.LogDebug("loginJWTHandler - failed to validate jwt for request: %s", err)
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

	user, err = user.UpsertUserByGoogleID(r.Context(), db.Users.FindOneAndUpdate)
	return
}
