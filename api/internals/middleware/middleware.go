package middleware

import (
	"net/http"

	"github.com/codemaverick07/api/internals/store"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

func SetUser(r *http.Request, user *store.User) *http.Request {
	return nil
}
