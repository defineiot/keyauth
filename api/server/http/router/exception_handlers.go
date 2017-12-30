package router

import (
	"net/http"

	"openauth/api/exception"
	"openauth/api/server/http/response"
)

type notFoundHandler struct{}
type methodNotAllowedHandler struct{}

func (n *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response.Failed(w, exception.NewNotFound("here is no handler for this endpoint"))
	return
}

func (m *methodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response.Failed(w, exception.NewMethodNotAllowed("this endpoint not support %s method", r.Method))
	return
}
