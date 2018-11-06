package router

import (
	"net/http"

	"github.com/defineiot/keyauth/api/http/response"
	"github.com/defineiot/keyauth/internal/exception"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.Failed(w, exception.NewNotFound("here is no handler for this endpoint"))
	return
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	response.Failed(w, exception.NewMethodNotAllowed("this endpoint not support %s method", r.Method))
	return
}
