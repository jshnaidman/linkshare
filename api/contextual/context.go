package contextual

import (
	"context"
	"linkshare_api/graph/model"
)

// private keys for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

// var sessionCtxKey = &contextKey{"session"}

type contextKey struct {
	name string
}

// UserForContext finds the user from the context. REQUIRES auth Middleware to have run.
// If not logged in, then user will be null.
func UserForContext(ctx context.Context) *model.User {
	user, _ := ctx.Value(userCtxKey).(*model.User)
	return user
}

// SessionForContext finds the session from the context. REQUIRES auth Middleware to have run.
// func SessionForContext(ctx context.Context) *Session {
// 	session, _ := ctx.Value(sessionCtxKey).(*Session)
// 	return session
// }
