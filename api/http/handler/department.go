package handler

import (
	"net/http"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/http/context"
	"github.com/defineiot/keyauth/api/http/request"
	"github.com/defineiot/keyauth/api/http/response"
	"github.com/defineiot/keyauth/dao/department"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateDepartment use to create user's application
func CreateDepartment(w http.ResponseWriter, r *http.Request) {
	tk := context.GetTokenFromContext(r)
	uid := tk.UserID

	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 校验传入参数的合法性
	name := val.Get("name").ToString()
	parentID := val.Get("parent_id").ToString()

	if name == "" {
		response.Failed(w, exception.NewBadRequest("name  missed"))
		return
	}

	dep := &department.Department{
		Name:      name,
		DomainID:  tk.DomainID,
		ParentID:  parentID,
		ManagerID: uid,
	}

	// 交给业务控制层处理
	if err := global.Store.CreateDepartment(dep); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, dep)
	return
}

// DeleteDepartment delete an application
func DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	did := ps.ByName("did")

	tk := context.GetTokenFromContext(r)
	app, err := global.Store.GetUserApp(did)
	if err != nil {
		response.Failed(w, err)
		return
	}
	if app.UserID != tk.UserID {
		response.Failed(w, exception.NewForbidden("application: %s is not belone to you", did))
		return
	}

	if app.ID == tk.ApplicationID {
		response.Failed(w, exception.NewForbidden("the application: %s your has used now, can't be deleted", did))
		return
	}

	// TODO: get token from context, and check permission
	if err := global.Store.DeleteApplication(did); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}
