package handler

import (
	"net/http"

	"github.com/defineiot/keyauth/dao/application"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/http/context"
	"github.com/defineiot/keyauth/api/http/request"
	"github.com/defineiot/keyauth/api/http/response"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateApplication use to create user's application
func CreateApplication(w http.ResponseWriter, r *http.Request) {
	tk := context.GetTokenFromContext(r)
	uid := tk.UserID

	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 校验传入参数的合法性
	name := val.Get("name").ToString()
	desc := val.Get("description").ToString()
	redirectURI := val.Get("redirect_uri").ToString()
	website := val.Get("website").ToString()

	if name == "" {
		response.Failed(w, exception.NewBadRequest("name  missed"))
		return
	}

	app := &application.Application{
		UserID:      uid,
		Name:        name,
		RedirectURI: redirectURI,
		Description: desc,
		Website:     website,
	}

	// 交给业务控制层处理
	if err := global.Store.CreateApplication(app); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, app)
	return
}

// DeleteApplication delete an application
func DeleteApplication(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	aid := ps.ByName("aid")

	tk := context.GetTokenFromContext(r)
	app, err := global.Store.GetUserApp(aid)
	if err != nil {
		response.Failed(w, err)
		return
	}
	if app.UserID != tk.UserID {
		response.Failed(w, exception.NewForbidden("application: %s is not belone to you", aid))
		return
	}

	if app.ID == tk.ApplicationID {
		response.Failed(w, exception.NewForbidden("the application: %s your has used now, can't be deleted", aid))
		return
	}

	// TODO: get token from context, and check permission
	if err := global.Store.DeleteApplication(aid); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}

// ListUserApplications get user's applications
func ListUserApplications(w http.ResponseWriter, r *http.Request) {
	tk := context.GetTokenFromContext(r)
	uid := tk.UserID

	apps, err := global.Store.ListUserApps(uid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, apps)
	return
}

// GetApplication get user's applications
func GetApplication(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	aid := ps.ByName("aid")

	tk := context.GetTokenFromContext(r)
	app, err := global.Store.GetUserApp(aid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	if app.UserID != tk.UserID {
		response.Failed(w, exception.NewForbidden("application: %s is not belone to you", aid))
		return
	}

	response.Success(w, http.StatusOK, app)
	return
}
