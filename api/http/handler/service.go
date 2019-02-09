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

// CreateService use to create an service
func CreateService(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	name := val.Get("name").ToString()
	desc := val.Get("description").ToString()

	var stype models.ServiceType
	switch val.Get("type").ToString() {
	case "controller_pannel":
		stype = models.Public
	case "data_pannel":
		stype = models.Agent
	case "internal_rpc":
		stype = models.Internal
	default:
		response.Failed(w, exception.NewBadRequest("unknown service type, support type (controller_pannel,data_pannel,internal_rpc)"))
		return
	}

	if name == "" {
		response.Failed(w, exception.NewBadRequest("service name missed"))
		return
	}

	svr := &models.Service{
		Name:        name,
		Description: desc,
		Type:        stype,
		Enabled:     true,
	}

	// 交给业务控制层处理
	if err := global.Store.CreateService(svr); err != nil {
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
	sid := ps.ByName("sid")

	svr, err := global.Store.GetService(sid)
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
	sid := ps.ByName("sid")

	err := global.Store.DeleteService(sid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}

// ListServiceFeatures todo
func ListServiceFeatures(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	sid := ps.ByName("sid")

	features, err := global.Store.ListServiceFeatures(sid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, features)
	return
}

// FeatureRegistryReq 服务功能注册接口
type FeatureRegistryReq struct {
	Version  string            `json:"version"`
	Features []*models.Feature `json:"features"`
}

// RegistryServiceFeatures todo
func RegistryServiceFeatures(w http.ResponseWriter, r *http.Request) {
	tk := context.GetTokenFromContext(r)
	sid := tk.ServiceID

	if sid == "" {
		response.Failed(w, exception.NewBadRequest("service id not found in token"))
		return
	}

	body, err := request.CheckBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	frReq := new(FeatureRegistryReq)
	if err := json.Unmarshal(body, frReq); err != nil {
		response.Failed(w, err)
		return
	}

	if frReq.Version == "" {
		response.Failed(w, exception.NewBadRequest("service version required"))
		return
	}

	if len(frReq.Features) == 0 {
		response.Failed(w, exception.NewBadRequest(`body array parse not features`))
		return
	}

	if err := global.Store.RegistryServiceFeatures(sid, frReq.Version, frReq.Features...); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, nil)
	return
}
