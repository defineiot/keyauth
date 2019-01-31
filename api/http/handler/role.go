package handler

import (
	"net/http"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/http/context"
	"github.com/defineiot/keyauth/api/http/request"
	"github.com/defineiot/keyauth/api/http/response"
	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateRole use to create an role
func CreateRole(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	name := val.Get("name").ToString()
	desc := val.Get("description").ToString()

	if name == "" {
		response.Failed(w, exception.NewBadRequest("role name missed"))
		return
	}

	ro := &role.Role{
		Name:        name,
		Description: desc,
	}

	// 交给业务控制层处理
	if err := global.Store.CreateRole(ro); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, ro)
	return
}

// ListRoles todo
func ListRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := global.Store.ListRoles()
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, roles)
	return
}

// GetRole todo
func GetRole(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	ri := ps.ByName("ri")

	role, err := global.Store.GetRole(ri)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, role)
	return
}

// DeleteRole todo
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	ri := ps.ByName("ri")

	err := global.Store.DeleteRole(ri)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}

// AddFeaturesToRole todo
func AddFeaturesToRole(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	rn := ps.ByName("ri")

	iter, err := request.CheckArrayBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	fids := make([]string, 0)
	for iter.ReadArray() {
		fids = append(fids, iter.ReadString())
	}
	if iter.Error != nil {
		response.Failed(w, exception.NewBadRequest("get feature id from body array error, %s", iter.Error))
		return
	}

	if len(fids) == 0 {
		response.Failed(w, exception.NewBadRequest(`body array hase no feature id`))
		return
	}

	if err = global.Store.AddFeaturesToRole(rn, fids...); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, nil)
	return
}

// RemoveFeaturesFromRole todo
func RemoveFeaturesFromRole(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	rn := ps.ByName("ri")

	iter, err := request.CheckArrayBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	fids := make([]string, 0)
	for iter.ReadArray() {
		fids = append(fids, iter.ReadString())
	}
	if iter.Error != nil {
		response.Failed(w, exception.NewBadRequest("get feature id from body array error, %s", iter.Error))
		return
	}

	if len(fids) == 0 {
		response.Failed(w, exception.NewBadRequest(`body array hase no feature id`))
		return
	}

	if err = global.Store.RemoveFeaturesFromRole(rn, fids...); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, nil)
	return
}
