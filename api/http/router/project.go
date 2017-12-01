package router

import (
	"openauth/api/http/handler"
)

// RouteToProject use to deal message
func (r *MyRouter) RouteToProject() {
	// validate message
	r.HandlerFunc("POST", "/projects/", handler.CreateProject)
	r.HandlerFunc("GET", "/projects/", handler.ListProject)
	r.HandlerFunc("GET", "/projects/:pid/", handler.GetProject)
	r.HandlerFunc("DELETE", "/projects/:pid/", handler.DeleteProject)
}
