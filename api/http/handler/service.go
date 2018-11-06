package handler

import (
	"net/http"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/http/context"
	"github.com/defineiot/keyauth/api/http/request"
	"github.com/defineiot/keyauth/api/http/response"
	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateService use to create an service
func CreateService(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	name := val.Get("name").ToString()
	desc := val.Get("description").ToString()

	if name == "" {
		response.Failed(w, exception.NewBadRequest("service name missed"))
		return
	}

	// 交给业务控制层处理
	svr, err := global.Store.CreateService(name, desc)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, svr)
	return
}

// ListServices todo
func ListServices(w http.ResponseWriter, r *http.Request) {
	svrs, err := global.Store.ListServices()
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, svrs)
	return
}

// GetService todo
func GetService(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	sn := ps.ByName("sn")

	svr, err := global.Store.GetService(sn)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, svr)
	return
}

// DeleteService todo
func DeleteService(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	sn := ps.ByName("sn")

	err := global.Store.DeleteService(sn)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}

// RegistryServiceFeatures todo
func RegistryServiceFeatures(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	sn := ps.ByName("sn")

	iter, err := request.CheckArrayBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	features := []service.Feature{}
	for iter.ReadArray() {
		f := service.Feature{}
		for l1Field := iter.ReadObject(); l1Field != ""; l1Field = iter.ReadObject() {
			switch l1Field {
			case "name":
				f.Name = iter.ReadString()
			case "method":
				f.Method = iter.ReadString()
			case "endpoint":
				f.Endpoint = iter.ReadString()
			default:
				iter.Skip()
			}
		}

		features = append(features, f)
	}

	if iter.Error != nil {
		response.Failed(w, exception.NewBadRequest("get service features from body array error, %s", iter.Error))
		return
	}

	if len(features) == 0 {
		response.Failed(w, exception.NewBadRequest(`body array parse not features`))
		return
	}

	if err := global.Store.RegistryServiceFeatures(sn, features...); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, nil)
	return
}

// ListServiceFeatures todo
func ListServiceFeatures(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	sn := ps.ByName("sn")

	features, err := global.Store.ListServiceFeatures(sn)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, features)
	return
}
