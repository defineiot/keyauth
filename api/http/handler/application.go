package handler

import (
	"net/http"
	"openauth/api/exception"
	"openauth/api/http/context"
	"openauth/api/http/request"
	"openauth/api/http/response"
)

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
	ok, err := usersrv.CheckUserIsExistByID(uid)
	if err != nil {
		response.Failed(w, err)
		return
	}
	if !ok {
		response.Failed(w, exception.NewBadRequest("user %s not exist", uid))
		return
	}

	app, err := appsrv.Registration(uid, name, redirectURI, clientType, desc, website)
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
	if err := appsrv.Unregistration(aid); err != nil {
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

	ok, err := usersrv.CheckUserIsExistByID(uid)
	if err != nil {
		response.Failed(w, err)
		return
	}
	if !ok {
		response.Failed(w, exception.NewBadRequest("user %s not exist", uid))
		return
	}

	apps, err := appsrv.GetUserApps(uid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, apps)
	return
}
