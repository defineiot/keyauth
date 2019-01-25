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

	// Token
	r.HandlerFunc("POST", "/oauth2/tokens/", "IssueToken", handler.IssueToken)
	r.HandlerFunc("GET", "/oauth2/tokens/", "ValidateToken", handler.ValidateToken)
	r.HandlerFunc("DELETE", "/oauth2/tokens/", "RevolkToken", handler.RevolkToken)

	// Project
	r.HandlerFunc("POST", "/projects/", "CreateProject", handler.CreateProject)
	r.HandlerFunc("GET", "/projects/", "ListDomainProjects", handler.ListDomainProjects)
	r.HandlerFunc("GET", "/self/projects/", "ListUserProjects", handler.ListUserProjects)
	r.HandlerFunc("GET", "/projects/:pid/", "GetProject", handler.GetProject)
	r.HandlerFunc("DELETE", "/projects/:pid/", "DeleteProject", handler.DeleteProject)
	// r.HandlerFunc("GET", "/v1/projects/:pid/users/", "ListProjectUser", handler.ListProjectUser)
	// r.HandlerFunc("POST", "/v1/projects/:pid/users/", "AddUsersToProject", handler.AddUsersToProject)
	// r.HandlerFunc("DELETE", "/v1/projects/:pid/users/", "RemoveUsersFromProject", handler.RemoveUsersFromProject)
	// // r.HandlerFunc("PUT", "/v1/projects/:pid/", handler.UpdateProject)

	// // Application
	// r.HandlerFunc("POST", "/v1/applications/", "CreateApplication", handler.CreateApplication)
	// r.HandlerFunc("GET", "/v1/applications/", "ListApplications", handler.ListApplications)
	// r.HandlerFunc("GET", "/v1/applications/:aid/", "GetApplication", handler.GetApplication)
	// r.HandlerFunc("DELETE", "/v1/applications/:aid/", "DeleteApplication", handler.DeleteApplication)
	// // r.HandlerFunc("PUT", "/v1/users/:uid/applications/:aid/", handler.UpdateApplication)

	// // Service
	// r.HandlerFunc("POST", "/v1/services/", "CreateService", handler.CreateService)
	// r.HandlerFunc("GET", "/v1/services/", "ListServices", handler.ListServices)
	// r.HandlerFunc("GET", "/v1/services/:sn/", "GetService", handler.GetService)
	// r.HandlerFunc("DELETE", "/v1/services/:sn/", "DeleteService", handler.DeleteService)
	// r.HandlerFunc("POST", "/v1/services/:sn/features/", "RegistryServiceFeatures", handler.RegistryServiceFeatures)
	// r.HandlerFunc("GET", "/v1/services/:sn/features/", "ListServiceFeatures", handler.ListServiceFeatures)

	// // Role
	// r.HandlerFunc("POST", "/v1/roles/", "CreateRole", handler.CreateRole)
	// r.HandlerFunc("GET", "/v1/roles/", "ListRoles", handler.ListRoles)
	// r.HandlerFunc("GET", "/v1/roles/:rn/", "GetRole", handler.GetRole)
	// r.HandlerFunc("DELETE", "/v1/roles/:rn/", "DeleteRole", handler.DeleteRole)
	// r.HandlerFunc("POST", "/v1/roles/:rn/features/", "AddFeaturesToRole", handler.AddFeaturesToRole)
	// r.HandlerFunc("DELETE", "/v1/roles/:rn/features/", "RemoveFeaturesFromRole", handler.RemoveFeaturesFromRole)

	r.AddV1Root()
}
