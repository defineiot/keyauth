package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/json-iterator/go"

	"openauth/api/controller/domain"
	"openauth/api/exception"
	"openauth/api/http/context"
	"openauth/api/http/response"
	"openauth/pkg/user"
)

// CreateDomain use to create domain
func CreateDomain(w http.ResponseWriter, r *http.Request) {
	// 请求检测
	if r.ContentLength == 0 {
		response.Failed(w, http.StatusBadRequest, "request body is empty")
		return
	}
	if r.ContentLength > 20971520 {
		response.Failed(w, http.StatusBadRequest, "the body exceeding the maximum limit, max size 20M")
		return
	}

	// 读取body数据
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Failed(w, http.StatusInternalServerError, err.Error())
		return
	}

	iter := jsoniter.ParseBytes(jsoniter.ConfigDefault, body)
	val := iter.ReadAny()
	if val.ValueType() != jsoniter.ObjectValue {
		response.Failed(w, http.StatusBadRequest, "body must be an valid json object")
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
	dm, err := domain.GetController()
	if err != nil {
		response.Failed(w, http.StatusInternalServerError, err.Error())
	}

	dom, err := dm.CreateDomain(name, desc, disp, user.Credential{})
	if err != nil {
		response.Failed(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, dom)
	return
}

// GetDomain use to get domain
func GetDomain(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	did := ps.ByName("did")

	// TODO: get token from context, and check permission
	dm, err := domain.GetController()
	if err != nil {
		response.Failed(w, http.StatusInternalServerError, err.Error())
		return
	}

	dom, err := dm.GetDomain(did)
	if apiErr, ok := err.(*exception.APIException); ok {
		response.Failed(w, apiErr.Code, apiErr.Error())
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
