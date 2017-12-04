package router

import "openauth/api/http/handler"

// RouteToUser use to deal message
func (r *MyRouter) RouteToUser() {
	// validate message
	r.HandlerFunc("POST", "/users/", handler.CreateUser)
	// r.HandlerFunc("GET", "/users/", handler.)
	r.HandlerFunc("GET", "/users/:uid/", handler.RetreveUser)
	// r.HandlerFunc("DELETE", "/users/:uid/", handler.DeleteProject)
}
