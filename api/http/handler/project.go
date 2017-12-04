package handler

import (
	"net/http"

	"openauth/api/http/context"
	"openauth/api/http/request"
	"openauth/api/http/response"
	"openauth/storage/user"
)

// TestDomainID Just for test
const TestDomainID = "fa735972-b059-44f3-b95f-78f0aaa1306e"

// CreateProject use to create an project
func CreateProject(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// get did from token
	did := TestDomainID

	name := val.Get("name").ToString()
	desc := val.Get("description").ToString()

	if name == "" {
		response.Failed(w, err)
		return
	}

	// 交给业务控制层处理

	proj, err := projectctl.CreateProject(did, name, desc, user.Credential{})
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, proj)
	return
}

// GetProject use to get one project
func GetProject(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	pid := ps.ByName("pid")

	// TODO: get token from context, and check permission

	proj, err := projectctl.GetProject(pid, user.Credential{})
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, proj)
	return
}

// ListProject use to list all project
func ListProject(w http.ResponseWriter, r *http.Request) {

	projects, err := projectctl.ListProject(TestDomainID)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, projects)
	return
}

// DeleteProject use to delete an project
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	pid := ps.ByName("pid")

	// TODO: get token from context, and check permission

	if err := projectctl.DestroyProject(pid, user.Credential{}); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}
