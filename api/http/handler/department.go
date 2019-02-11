package handler

import (
	"encoding/json"
	"net/http"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/http/context"
	"github.com/defineiot/keyauth/api/http/request"
	"github.com/defineiot/keyauth/api/http/response"
	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/internal/exception"
)

type createDepReq struct {
	Name     string   `json:"name,omitempty"`
	ParentID string   `json:"parent_id,omitempty"`
	Projects []string `json:"projects,omitempty"`
	Roles    []string `json:"roles,omitempty"`
}

// CreateDepartment use to create domain department
func CreateDepartment(w http.ResponseWriter, r *http.Request) {
	tk := context.GetTokenFromContext(r)
	uid := tk.UserID

	body, err := request.CheckBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 校验传入参数的合法性
	reqData := new(createDepReq)
	if err := json.Unmarshal(body, reqData); err != nil {
		response.Failed(w, err)
		return
	}

	if reqData.Name == "" {
		response.Failed(w, exception.NewBadRequest("department name missed"))
		return
	}

	dep := &models.Department{
		Name:       reqData.Name,
		DomainID:   tk.DomainID,
		ParentID:   reqData.ParentID,
		ManagerID:  uid,
		ProjectIDs: reqData.Projects,
		RoleIDs:    reqData.Roles,
	}

	// 交给业务控制层处理
	if err := global.Store.CreateDepartment(dep); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, dep)
	return
}

// DeleteDepartment delete an department
func DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	did := ps.ByName("did")

	tk := context.GetTokenFromContext(r)
	dep, err := global.Store.GetDepartment(did)
	if err != nil {
		response.Failed(w, err)
		return
	}
	if dep.DomainID != tk.DomainID {
		response.Failed(w, exception.NewForbidden("department: %s is not belone to you domain", did))
		return
	}

	// TODO: get token from context, and check permission
	if err := global.Store.DeleteDepartment(did); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}

// ListSubDepartments list sub department
func ListSubDepartments(w http.ResponseWriter, r *http.Request) {
	tk := context.GetTokenFromContext(r)

	qs := r.URL.Query()
	parentID := qs.Get("parent_id")

	if parentID == "" {
		parentID = "/"
	}

	deps, err := global.Store.ListSubDepartments(tk.DomainID, parentID)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, deps)
	return
}

// GetDepartment get  department
func GetDepartment(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	did := ps.ByName("did")

	tk := context.GetTokenFromContext(r)
	dep, err := global.Store.GetDepartment(did)
	if err != nil {
		response.Failed(w, err)
		return
	}

	if dep.DomainID != tk.DomainID {
		response.Failed(w, exception.NewForbidden("department: %s is not belone to you", did))
		return
	}

	response.Success(w, http.StatusOK, dep)
	return
}
