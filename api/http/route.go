package http

import (
	"github.com/defineiot/keyauth/api/http/handler"
	"github.com/defineiot/keyauth/api/http/router"
)

// RouteToV1 use to route to ep
func RouteToV1(r *router.MyRouter) {
	// User
	r.HandlerFunc("POST", "/members/", "CreateMemberUser", handler.CreateMemberUser)
	r.HandlerFunc("GET", "/members/", "ListMemberUsers", handler.ListMemberUsers)
	r.HandlerFunc("GET", "/members/:mid/", "GetMemberUser", handler.GetMemberUser)
	r.HandlerFunc("DELETE", "/members/:mid/", "DeleteMemberUser", handler.DeleteMemberUser)

	// Token
	r.HandlerFunc("POST", "/oauth2/tokens/", "IssueToken", handler.IssueToken)
	r.HandlerFunc("GET", "/oauth2/tokens/:tk/", "ValidateToken", handler.ValidateToken)
	r.HandlerFunc("DELETE", "/oauth2/tokens/:tk/", "RevolkToken", handler.RevolkToken)

	// Project
	r.HandlerFunc("POST", "/projects/", "CreateProject", handler.CreateProject)
	r.HandlerFunc("GET", "/projects/", "ListDomainProjects", handler.ListDomainProjects)
	r.HandlerFunc("GET", "/self/projects/", "ListUserProjects", handler.ListUserProjects)
	r.HandlerFunc("GET", "/projects/:pid/", "GetProject", handler.GetProject)
	r.HandlerFunc("DELETE", "/projects/:pid/", "DeleteProject", handler.DeleteProject)
	r.HandlerFunc("GET", "/projects/:pid/members/", "ListProjectUser", handler.ListProjectUser)
	r.HandlerFunc("POST", "/projects/:pid/members/", "AddUsersToProject", handler.AddUsersToProject)
	r.HandlerFunc("DELETE", "/projects/:pid/members/", "RemoveUsersFromProject", handler.RemoveUsersFromProject)

	// Application
	r.HandlerFunc("POST", "/applications/", "CreateApplication", handler.CreateApplication)
	r.HandlerFunc("GET", "/applications/", "ListUserApplications", handler.ListUserApplications)
	r.HandlerFunc("GET", "/applications/:aid/", "GetApplication", handler.GetApplication)
	r.HandlerFunc("DELETE", "/applications/:aid/", "DeleteApplication", handler.DeleteApplication)
	// // r.HandlerFunc("PUT", "/v1/users/:uid/applications/:aid/", handler.UpdateApplication)

	// Service
	r.HandlerFunc("POST", "/services/", "CreateService", handler.CreateService)
	r.HandlerFunc("GET", "/services/", "ListServices", handler.ListServices)
	r.HandlerFunc("GET", "/services/:sid/", "GetService", handler.GetService)
	r.HandlerFunc("DELETE", "/services/:sid/", "DeleteService", handler.DeleteService)
	r.HandlerFunc("POST", "/features/", "RegistryServiceFeatures", handler.RegistryServiceFeatures)
	r.HandlerFunc("GET", "/services/:sid/features/", "ListServiceFeatures", handler.ListServiceFeatures)

	// Role
	r.HandlerFunc("POST", "/roles/", "CreateRole", handler.CreateRole)
	r.HandlerFunc("GET", "/roles/", "ListRoles", handler.ListRoles)
	r.HandlerFunc("GET", "/roles/:ri/", "GetRole", handler.GetRole)
	r.HandlerFunc("DELETE", "/roles/:ri/", "DeleteRole", handler.DeleteRole)
	r.HandlerFunc("POST", "/roles/:ri/features/", "AddFeaturesToRole", handler.AddFeaturesToRole)
	r.HandlerFunc("DELETE", "/roles/:ri/features/", "RemoveFeaturesFromRole", handler.RemoveFeaturesFromRole)

	// r.HandlerFunc("POST", "/v1/domains/users/", "CreateDomainUser", handler.CreateDomainUser)
	// r.HandlerFunc("GET", "/v1/users/:uid/domains/", "ListUserDomain", handler.ListUserDomain)
	// r.HandlerFunc("PUT", "/v1/users/:uid/password/", "SetUserPassword", handler.SetUserPassword)
	// r.HandlerFunc("DELETE", "/v1/unregistry/", "UnRegistry", handler.UnRegistry)
	// r.HandlerFunc("POST", "/v1/users/:uid/projects/", "AddProjectsToUser", handler.AddProjectsToUser)
	// r.HandlerFunc("DELETE", "/v1/users/:uid/projects/", "RemoveProjectsFromUser", handler.RemoveProjectsFromUser)
	// r.HandlerFunc("POST", "/v1/users/:uid/bind/roles/:rn/", "BindRole", handler.BindRole)
	// r.HandlerFunc("POST", "/v1/users/:uid/unbind/roles/:rn/", "UnBindRole", handler.UnBindRole)
	// r.HandlerFunc("POST", "/v1/invitations/", "InvitationsUser", handler.InvitationsUser)
	// r.HandlerFunc("DELETE", "/v1/invitations/:code/", "RevolkInvitation", handler.RevolkInvitation)
	// r.HandlerFunc("GET", "/v1/invitations/", "ListInvitationsRecords", handler.ListInvitationsRecords)
	// r.HandlerFunc("GET", "/v1/invitations/:code/", "GetInvitationsRecord", handler.GetInvitationsRecord)
	// r.HandlerFunc("POST", "/v1/registry/", "RegistryUser", handler.RegistryUser)
	// r.HandlerFunc("POST", "/v1/verifycode/", "IssueVerifyCode", handler.IssueVerifyCode)
	// r.HandlerFunc("POST", "/v1/invitations/users/:uid/code/:code/", "AcceptInvitation", handler.AcceptInvitation)
	//  r.HandlerFunc("PUT", "/v1/users/:uid/", handler.UpdateUser)
	// r.HandlerFunc("POST", "/v1/default/projects/:pid/", "SetUserDefaultProject", handler.SetUserDefaultProject)

	r.AddV1Root()
}
