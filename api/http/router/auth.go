package router

import "openauth/api/http/handler"

// RouteToAuth use to deal message
func (r *MyRouter) RouteToAuth() {
	// Token Endpoint https://tools.ietf.org/html/rfc6749#section-3.2
	r.HandlerFunc("POST", "/oauth2/tokens/", handler.CreateProject)
	r.HandlerFunc("GET", "/oauth2/token/", handler.ListDomain)
	r.HandlerFunc("GET", "/oauth2/tokens/", handler.ListProject)
	r.HandlerFunc("DELETE", "/oauth2/tokens/", handler.DeleteProject)

	// Authorization Endpoint https://tools.ietf.org/html/rfc6749#section-3.1
	r.HandlerFunc("GET", "/oauth2/authorize/", handler.ListProjectUsers)
	r.HandlerFunc("POST", "/projects/:pid/users/", handler.AddUsersToProject)
	r.HandlerFunc("DELETE", "/projects/:pid/users/", handler.RemoveUsersFromProject)
}
