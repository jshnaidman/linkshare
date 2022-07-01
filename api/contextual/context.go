package contextual

import (
	"context"
	"linkshare_api/database"
	"linkshare_api/graph/model"
)

// This package tracks all context keys

type contextKey struct {
	name string
}

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}
var sessionCtxKey = &contextKey{"session"}

// This should only be called in the auth module during login.
func AddUserToContext(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

// UserForContext finds the user from the context. REQUIRES auth Middleware to have run.
// If not logged in, then user will be null.
func UserForContext(ctx context.Context) *model.User {
	user, _ := ctx.Value(userCtxKey).(*model.User)
	return user
}

// This should only be called in the auth module in auth middleware.
func AddSessionToContext(ctx context.Context, session *database.Session) context.Context {
	return context.WithValue(ctx, userCtxKey, session)
}

// SessionForContext finds the session from the context. REQUIRES auth Middleware to have run.
func SessionForContext(ctx context.Context) *database.Session {
	session, _ := ctx.Value(sessionCtxKey).(*database.Session)
	return session
}
