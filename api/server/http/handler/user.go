package handler

import (
	"net/http"

	"openauth/api/exception"
	"openauth/api/server/http/context"
	"openauth/api/server/http/request"
	"openauth/api/server/http/response"
)

// CreateUser handler
func CreateUser(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// get did from token
	did := TestDomainID
	name := val.Get("name").ToString()
	pass := val.Get("password").ToString()
	// desc := val.Get("description").ToString()

	if name == "" || pass == "" {
		response.Failed(w, exception.NewBadRequest("name or password is missed"))
		return
	}

	// 交给业务控制层处理
	user, err := usersrv.CreateUser(did, name, pass, "")
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, user)
	return
}

// RetreveUser use to get user
func RetreveUser(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")

	// TODO: get token from context, and check permission
	u, err := usersrv.GetUser(uid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, u)
	return
}

// ListUser use to list domain users
func ListUser(w http.ResponseWriter, r *http.Request) {

	users, err := usersrv.ListUser(TestDomainID)
	if err != nil {
		response.Failed(w, err)
		return
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

	if err := usersrv.SetUserDefaultProject(uid, pid); err != nil {
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

	projects, err := usersrv.ListUserProjects(uid)
	if err != nil {
		response.Failed(w, err)
		return
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
	if err := usersrv.RemoveProjectsFromUser(uid, pids...); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "")
	return
}

// RegisteApplication use to registe app
func RegisteApplication(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")

	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 校验传入参数的合法性
	name := val.Get("name").ToString()
	desc := val.Get("description").ToString()
	redirectURI := val.Get("redirect_uri").ToString()
	clientType := val.Get("client_type").ToString()
	website := val.Get("website").ToString()

	if name == "" {
		response.Failed(w, exception.NewBadRequest("name  missed"))
		return
	}

	// 交给业务控制层处理
	app, err := appsrv.RegisterApplication(uid, name, redirectURI, clientType, desc, website)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, app)
	return
}

// UnRegisteApplication delete application
func UnRegisteApplication(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	aid := ps.ByName("aid")

	// TODO: get token from context, and check permission
	if err := appsrv.UnregisteApplication(aid); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}

// GetUserApplications get user's application
func GetUserApplications(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")

	apps, err := appsrv.GetUserApplications(uid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, apps)
	return
}
