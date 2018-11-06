package handler

import (
	"net/http"
	"strconv"

	"github.com/json-iterator/go"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/http/context"
	"github.com/defineiot/keyauth/api/http/request"
	"github.com/defineiot/keyauth/api/http/response"
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

	// 交给业务控制层处理
	role, err := global.Store.CreateRole(name, desc)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, role)
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
	rn := ps.ByName("rn")

	role, err := global.Store.GetRole(rn)
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
	rn := ps.ByName("rn")

	err := global.Store.DeleteRole(rn)
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
	rn := ps.ByName("rn")

	iter, err := request.CheckArrayBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	fids := make([]int64, 0)
	for iter.ReadArray() {
		var fid int64
		var err error
		switch iter.WhatIsNext() {
		case jsoniter.NumberValue:
			fid = iter.ReadInt64()
		case jsoniter.StringValue:
			fidStr := iter.ReadString()
			fid, err = strconv.ParseInt(fidStr, 10, 64)
			if err != nil {
				response.Failed(w, exception.NewBadRequest("parse feature id from body array error, %s", iter.Error))
				return
			}
		default:
			response.Failed(w, exception.NewBadRequest("get feature id from body array only support string or number"))
			return
		}

		fids = append(fids, fid)
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
	rn := ps.ByName("rn")

	iter, err := request.CheckArrayBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	fids := make([]int64, 0)
	for iter.ReadArray() {
		fid, err := iter.ReadNumber().Int64()
		if err != nil {
			response.Failed(w, exception.NewBadRequest("get feature id from body array error, %s", iter.Error))
			return
		}
		fids = append(fids, fid)
	}
	if iter.Error != nil {
		response.Failed(w, exception.NewBadRequest("get feature id from body array error, %s", iter.Error))
		return
	}

	if len(fids) == 0 {
		response.Failed(w, exception.NewBadRequest(`body array hase no feature id`))
		return
	}

	if _, err = global.Store.RemoveFeaturesFromRole(rn, fids...); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, nil)
	return
}
