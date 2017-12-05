package router

import "openauth/api/http/handler"

// RouteToUser use to deal message
func (r *MyRouter) RouteToUser() {
	// validate message
	r.HandlerFunc("POST", "/users/", handler.CreateUser)
	r.HandlerFunc("GET", "/users/", handler.ListUser)
	r.HandlerFunc("GET", "/users/:uid/", handler.RetreveUser)
	r.HandlerFunc("DELETE", "/users/:uid/", handler.DeleteUser)

	r.HandlerFunc("POST", "/users/:uid/default/project/:pid/", handler.SetUserDefaultProject)
	r.HandlerFunc("GET", "/users/:uid/projects/", handler.ListUserProjects)
	r.HandlerFunc("POST", "/users/:uid/projects/", handler.AddProjectsToUser)
	r.HandlerFunc("DELETE", "/users/:uid/projects/", handler.RemoveProjectsFromUser)
}
