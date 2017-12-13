package handler

import (
	"net/http"

	"openauth/api/exception"
	"openauth/api/http/context"
	"openauth/api/http/request"
	"openauth/api/http/response"
)

// CreateDomain use to create domain
func CreateDomain(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 校验传入参数的合法性
	name := val.Get("name").ToString()
	desc := val.Get("description").ToString()
	disp := val.Get("display_name").ToString()

	if name == "" || disp == "" {
		response.Failed(w, exception.NewBadRequest("name or display_name missed"))
		return
	}

	// 交给业务控制层处理
	dom, err := domainsrv.CreateDomain(name, desc, disp, true)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, dom)
	return
}

// GetDomain use to get domain
func GetDomain(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	did := ps.ByName("did")

	dom, err := domainsrv.GetDomain(did)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, dom)
	return
}

// ListDomain domain list
func ListDomain(w http.ResponseWriter, r *http.Request) {
	doms, err := domainsrv.ListDomain()
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, doms)
	return
}

// DeleteDomain destory an domain
func DeleteDomain(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	did := ps.ByName("did")

	// TODO: get token from context, and check permission
	if err := domainsrv.DeleteDomain(did); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}
