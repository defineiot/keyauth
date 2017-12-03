package handler

import (
	"net/http"

	"openauth/api/controller/project"
	"openauth/api/http/context"
	"openauth/api/http/request"
	"openauth/api/http/response"
	"openauth/pkg/user"
)

// CreateProject use to create an project
func CreateProject(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// get did from token
	did := "08b6d234-c06f-4acb-8ca4-fd3bd841f607"

	name := val.Get("name").ToString()
	desc := val.Get("description").ToString()

	if name == "" {
		response.Failed(w, err)
		return
	}

	// 交给业务控制层处理
	pc, err := project.GetController()
	if err != nil {
		response.Failed(w, err)
	}

	proj, err := pc.CreateProject(did, name, desc, user.Credential{})
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
	pc, err := project.GetController()
	if err != nil {
		response.Failed(w, err)
		return
	}

	proj, err := pc.GetProject(pid, user.Credential{})
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, proj)
	return
}

// ListProject use to list all project
func ListProject(w http.ResponseWriter, r *http.Request) {
	// TODO: get token from context, and check permission
	pc, err := project.GetController()
	if err != nil {
		response.Failed(w, err)
		return
	}

	projects, err := pc.ListProject("08b6d234-c06f-4acb-8ca4-fd3bd841f607")
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
	pc, err := project.GetController()
	if err != nil {
		response.Failed(w, err)
		return
	}

	if err := pc.DestroyProject(pid, user.Credential{}); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}
