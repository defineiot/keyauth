package handler

import (
	"net/http"

	"openauth/api/exception"
	"openauth/api/http/context"
	"openauth/api/http/request"
	"openauth/api/http/response"
)

// CreateUser handler
func CreateUser(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// get did from token

	name := val.Get("name").ToString()
	pass := val.Get("password").ToString()
	desc := val.Get("description").ToString()

	if name == "" || pass == "" {
		response.Failed(w, exception.NewBadRequest("name or password is missed"))
		return
	}

	// 交给业务控制层处理
	u, err := userctl.CreateUser(TestDomainID, name, pass, desc)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, u)
	return
}

// RetreveUser use to get user
func RetreveUser(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")

	// TODO: get token from context, and check permission

	u, err := userctl.GetUser(uid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, u)

	return
}

// ListUser use to list domain users
func ListUser(w http.ResponseWriter, r *http.Request) {
	users, err := userctl.ListUser(TestDomainID)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, users)
	return
}

// DeleteUser delete an user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")

	if err := userctl.DeleteUser(uid); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")

	return
}
