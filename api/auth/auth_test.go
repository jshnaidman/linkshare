package auth

import (
	"linkshare_api/database"
	"linkshare_api/graph/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleJWTLogin(t *testing.T) {
	respWriter := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", nil)
	handleJWTLogin(respWriter, req, func(bearerToken string, db *database.LinkShareDB, w http.ResponseWriter,
		r *http.Request) (user *model.User, err error) {

		return
	})
}
