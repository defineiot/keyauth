package handler

import (
	"net/http"

	"openauth/api/controller/domain"
	"openauth/api/exception"
	"openauth/api/http/context"
	"openauth/api/http/request"
	"openauth/api/http/response"
	"openauth/pkg/user"
)

// CreateDomain use to create domain
func CreateDomain(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, http.StatusBadRequest, err.Error())
		return
	}

	// 校验传入参数的合法性
	name := val.Get("name").ToString()
	desc := val.Get("description").ToString()
	disp := val.Get("display_name").ToString()

	if name == "" || disp == "" {
		response.Failed(w, http.StatusBadRequest, "name or display_name missed")
		return
	}

	// 交给业务控制层处理
	dc, err := domain.GetController()
	if err != nil {
		response.Failed(w, http.StatusInternalServerError, err.Error())
	}

	dom, err := dc.CreateDomain(name, desc, disp, user.Credential{})
	if err != nil {
		response.Failed(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, dom)
	return
}

// GetDomain use to get domain
func GetDomain(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	did := ps.ByName("did")

	// TODO: get token from context, and check permission
	dc, err := domain.GetController()
	if err != nil {
		response.Failed(w, http.StatusInternalServerError, err.Error())
		return
	}

	dom, err := dc.GetDomain(did)
	if err != nil {
		if apiErr, ok := err.(*exception.APIException); ok {
			response.Failed(w, apiErr.Code, apiErr.Error())
			return
		}

		response.Failed(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, dom)
	return
}

// ListDomain domain list
func ListDomain(w http.ResponseWriter, r *http.Request) {
	// TODO: get token from context, and check permission
	dm, err := domain.GetController()
	if err != nil {
		response.Failed(w, http.StatusInternalServerError, err.Error())
		return
	}

	doms, err := dm.ListDomain()
	if err != nil {
		response.Failed(w, http.StatusInternalServerError, err.Error())
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
	dm, err := domain.GetController()
	if err != nil {
		response.Failed(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := dm.DestoryDomain(did); err != nil {
		response.Failed(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}
