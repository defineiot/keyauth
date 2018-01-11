package handler

import (
	"net/http"

	"openauth/api/exception"
	"openauth/api/server/http/context"
	"openauth/api/server/http/request"
	"openauth/api/server/http/response"
)

// CreateService service
func CreateService(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 校验传入参数的合法性
	name := val.Get("name").ToString()
	desc := val.Get("description").ToString()

	if name == "" {
		response.Failed(w, exception.NewBadRequest("name missed"))
		return
	}

	s, err := svr.CreateService(name, desc)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, s)
	return
}

// DeleteService delete
func DeleteService(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	sid := ps.ByName("sid")

	if err := svr.DeleteService(sid); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}

// ListService list
func ListService(w http.ResponseWriter, r *http.Request) {
	svrs, err := svr.ListService()
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, svrs)
	return
}

// GetService get
func GetService(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	sid := ps.ByName("sid")

	s, err := svr.GetService(sid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, s)
	return
}
