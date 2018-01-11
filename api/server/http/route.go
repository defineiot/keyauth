package http

import (
	"openauth/api/server/http/handler"
	"openauth/api/server/http/router"
)

// RouteToV1 route to openauth v1 endpoint
func RouteToV1(r *router.MyRouter) {
	// validate message
	r.HandlerFunc("POST", "/v1/domains/", handler.CreateDomain)
	r.HandlerFunc("GET", "/v1/domains/", handler.ListDomain)
	r.HandlerFunc("GET", "/v1/domains/:did/", handler.GetDomain)
	r.HandlerFunc("DELETE", "/v1/domains/:did/", handler.DeleteDomain)

	// validate message
	r.HandlerFunc("POST", "/v1/projects/", handler.CreateProject)
	r.HandlerFunc("GET", "/v1/projects/", handler.ListProject)
	r.HandlerFunc("GET", "/v1/projects/:pid/", handler.GetProject)
	r.HandlerFunc("DELETE", "/v1/projects/:pid/", handler.DeleteProject)
	r.HandlerFunc("GET", "/v1/projects/:pid/users/", handler.ListProjectUsers)
	r.HandlerFunc("POST", "/v1/projects/:pid/users/", handler.AddUsersToProject)
	r.HandlerFunc("DELETE", "/v1/projects/:pid/users/", handler.RemoveUsersFromProject)

	// validate message
	r.HandlerFunc("POST", "/v1/users/", handler.CreateUser)
	r.HandlerFunc("GET", "/v1/users/", handler.ListUser)
	r.HandlerFunc("GET", "/v1/users/:uid/", handler.RetreveUser)
	r.HandlerFunc("DELETE", "/v1/users/:uid/", handler.DeleteUser)
	r.HandlerFunc("POST", "/v1/users/:uid/default/project/:pid/", handler.SetUserDefaultProject)
	r.HandlerFunc("GET", "/v1/users/:uid/projects/", handler.ListUserProjects)
	r.HandlerFunc("POST", "/v1/users/:uid/projects/", handler.AddProjectsToUser)
	r.HandlerFunc("DELETE", "/v1/users/:uid/projects/", handler.RemoveProjectsFromUser)
	r.HandlerFunc("POST", "/v1/users/:uid/applications/", handler.RegisteApplication)
	r.HandlerFunc("DELETE", "/v1/users/:uid/applications/:aid/", handler.UnRegisteApplication)
	r.HandlerFunc("GET", "/v1/users/:uid/applications/", handler.GetUserApplications)

	// Token Endpoint https://tools.ietf.org/html/rfc6749#section-3.2
	r.HandlerFunc("POST", "/v1/oauth2/tokens/", handler.IssueToken)
	r.HandlerFunc("GET", "/v1/oauth2/tokens/", handler.ValidateToken)
	// r.HandlerFunc("GET", "/v1/oauth2/tokens/", handler.ListProject)
	// r.HandlerFunc("DELETE", "/v1/oauth2/tokens/", handler.DeleteProject)
	// // Authorization Endpoint https://tools.ietf.org/html/rfc6749#section-3.1
	// r.HandlerFunc("GET", "/v1/oauth2/authorize/", handler.ListProjectUsers)

	r.HandlerFunc("POST", "/v1/services/", handler.CreateService)
	r.HandlerFunc("DELETE", "/v1/services/:sid/", handler.DeleteService)
	r.HandlerFunc("GET", "/v1/services/", handler.ListService)
	r.HandlerFunc("GET", "/v1/services/:sid/", handler.GetService)

	r.AddV1Root()
}
