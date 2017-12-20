package handler

import (
	"net/http"
	"strconv"

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
	var (
		ps  int64
		pn  int64
		err error
	)

	pageNumber := r.URL.Query().Get("page_number")
	pageSize := r.URL.Query().Get("page_size")

	if pageNumber != "" {
		pn, err = strconv.ParseInt(pageNumber, 10, 64)
		if err != nil {
			response.Failed(w, exception.NewBadRequest("page_size must be number"))
			return
		}
	}
	if pageSize != "" {
		ps, err = strconv.ParseInt(pageSize, 10, 64)
		if err != nil {
			response.Failed(w, exception.NewBadRequest("page_number must be number"))
			return
		}
	}

	doms, totalP, err := domainsrv.ListDomain(pn, ps)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.SuccessWithPage(w, http.StatusOK, doms, totalP)
	return
}

// DeleteDomain destory an domain
func DeleteDomain(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	did := ps.ByName("did")

	// TODO: get token from context, and check permission
	if err := domainsrv.DestroyDomain(did); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}
