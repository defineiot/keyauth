package handler

import (
	"net/http"

	"openauth/api/exception"
	"openauth/api/http/context"
	"openauth/api/http/request"
	"openauth/api/http/response"
	"openauth/pkg/project"
)

// CreateUser handler
func CreateUser(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// get did from token

	name := val.Get("name").ToString()
	pass := val.Get("password").ToString()
	// desc := val.Get("description").ToString()

	if name == "" || pass == "" {
		response.Failed(w, exception.NewBadRequest("name or password is missed"))
		return
	}

	// 交给业务控制层处理
	ok, err := domainsrv.CheckDomainIsExistByID(TestDomainID)
	if err != nil {
		response.Failed(w, err)
		return
	}
	if !ok {
		response.Failed(w, exception.NewBadRequest("domain %s not exist", TestDomainID))
		return
	}

	u, err := usersrv.CreateUser(TestDomainID, name, pass, true, 4096, 4096)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, u)
	return
}

// RetreveUser use to get user
func RetreveUser(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")

	// TODO: get token from context, and check permission
	u, err := usersrv.GetUserByID(uid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, u)
	return
}

// ListUser use to list domain users
func ListUser(w http.ResponseWriter, r *http.Request) {
	ok, err := domainsrv.CheckDomainIsExistByID(TestDomainID)
	if err != nil {
		response.Failed(w, err)
		return
	}
	if !ok {
		response.Failed(w, exception.NewBadRequest("domain %s not exist", TestDomainID))
		return
	}

	users, err := usersrv.ListUser(TestDomainID)
	if err != nil {
		response.Failed(w, err)
		return
	}
	for _, u := range users {
		if u.DefaultProjectID != "" {
			dp, err := projectsrv.GetProject(u.DefaultProjectID)
			if err != nil {
				response.Failed(w, exception.NewInternalServerError("get user %s project error, %s", u.Name, err))
				return
			}

			u.DefaultProject = dp
		}
	}

	response.Success(w, http.StatusOK, users)
	return
}

// DeleteUser delete an user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")

	if err := usersrv.DeleteUser(uid); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}

// SetUserDefaultProject set default project
func SetUserDefaultProject(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")
	pid := ps.ByName("pid")

	ok, err := projectsrv.CheckProjectIsExistByID(pid)
	if err != nil {
		response.Failed(w, err)
		return
	}
	if !ok {
		response.Failed(w, exception.NewBadRequest("project %s not exist", pid))
		return
	}

	if err := usersrv.SetDefaultProject(uid, pid); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "")
	return
}

// ListUserProjects list user's projects
func ListUserProjects(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")

	projectIDs, err := usersrv.ListUserProjects(uid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	projects := []*project.Project{}
	for _, pid := range projectIDs {
		pj, err := projectsrv.GetProject(pid)
		if err != nil {
			response.Failed(w, err)
			return
		}
		projects = append(projects, pj)
	}

	response.Success(w, http.StatusOK, projects)
	return
}

// AddProjectsToUser add projects
func AddProjectsToUser(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")

	iter, err := request.CheckArrayBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// get did from token
	pids := []string{}
	for iter.ReadArray() {
		pids = append(pids, iter.ReadString())
	}
	if iter.Error != nil {
		response.Failed(w, exception.NewBadRequest("json format decode error, %s", iter.Error))
		return
	}

	if len(pids) == 0 {
		response.Failed(w, exception.NewBadRequest("pids not find"))
		return
	}

	// 业务层
	// 1. 检测用户传入的pid是否属于这个用户
	upids, err := usersrv.ListUserProjects(uid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	invalidatePIDs := []string{}
	for _, pid := range pids {
		pidExist := false
		for _, upid := range upids {
			if pid == upid {
				pidExist = true
			}
		}
		if !pidExist {
			invalidatePIDs = append(invalidatePIDs, pid)
		}
	}

	if len(invalidatePIDs) != 0 {
		response.Failed(w, exception.NewBadRequest("the projects %s not owned by this user %s", invalidatePIDs, uid))
		return
	}
	// 2. 处理
	if err := usersrv.AddProjectsToUser(uid, pids...); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "")
	return
}

// RemoveProjectsFromUser remove projects
func RemoveProjectsFromUser(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")

	iter, err := request.CheckArrayBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// get did from token
	pids := []string{}
	for iter.ReadArray() {
		pids = append(pids, iter.ReadString())
	}
	if iter.Error != nil {
		response.Failed(w, exception.NewBadRequest("json format decode error, %s", iter.Error))
		return
	}

	if len(pids) == 0 {
		response.Failed(w, exception.NewBadRequest("pids not find"))
		return
	}

	// 业务层
	// 1. 检测用户传入的pid是否属于这个用户
	upids, err := usersrv.ListUserProjects(uid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	invalidatePIDs := []string{}
	for _, pid := range pids {
		pidExist := false
		for _, upid := range upids {
			if pid == upid {
				pidExist = true
			}
		}
		if !pidExist {
			invalidatePIDs = append(invalidatePIDs, pid)
		}
	}
	if len(invalidatePIDs) != 0 {
		response.Failed(w, exception.NewBadRequest("the projects %s not owned by this user %s", invalidatePIDs, uid))
		return
	}
	// 2. 处理
	if err := usersrv.RemoveProjectsFromUser(uid, pids...); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "")
	return
}
