package handler

import (
	"net/http"
	"strings"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/http/context"
	"github.com/defineiot/keyauth/api/http/request"
	"github.com/defineiot/keyauth/api/http/response"
	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/store"
)

// IssueToken use to issue access token
func IssueToken(w http.ResponseWriter, r *http.Request) {
	var grantType string

	tokenReq := new(store.TokenRequest)

	contentT := strings.Split(r.Header.Get("content-type"), ";")
	if len(contentT) == 0 {
		response.Failed(w, exception.NewBadRequest("content-type missed, your must be choice one [application/json, application/x-www-form-urlencoded]"))
		return
	}

	switch strings.TrimSpace(contentT[0]) {
	case "application/json":
		val, err := request.CheckObjectBody(r)
		if err != nil {
			response.Failed(w, err)
			return
		}
		tokenReq.ClientID = val.Get("client_id").ToString()
		tokenReq.ClientSecret = val.Get("client_secret").ToString()
		tokenReq.Username = val.Get("username").ToString()
		tokenReq.Password = val.Get("password").ToString()
		tokenReq.RefreshToken = val.Get("refresh_token").ToString()
		tokenReq.AccessToken = val.Get("access_token").ToString()
		tokenReq.Scope = val.Get("scope").ToString()
		grantType = val.Get("grant_type").ToString()

	case "", "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			response.Failed(w, exception.NewBadRequest("parse x-www-form-urlencoded data error, %s", err))
			return
		}
		tokenReq.ClientID = r.FormValue("client_id")
		tokenReq.ClientSecret = r.FormValue("client_secret")
		tokenReq.Username = r.FormValue("username")
		tokenReq.Password = r.FormValue("password")
		tokenReq.RefreshToken = r.FormValue("refresh_token")
		tokenReq.AccessToken = r.FormValue("access_token")
		grantType = r.FormValue("grant_type")

		tokenReq.Scope = r.FormValue("scope")
	default:
		response.Failed(w, exception.NewBadRequest("content-type only support for application/json and application/x-www-form-urlencoded others(%s) don't supported", r.Header.Get("content-type")))
		return
	}

	clientID, clientSecret, ok := r.BasicAuth()
	if ok {
		tokenReq.ClientID = clientID
		tokenReq.ClientSecret = clientSecret
	}

	switch grantType {
	case "authorization_code":
		tokenReq.GrantType = models.AUTHCODE
	case "implicit":
		tokenReq.GrantType = models.IMPLICIT
	case "password":
		tokenReq.GrantType = models.PASSWORD
	case "client_credentials":
		tokenReq.GrantType = models.CLIENT
	case "refresh_token":
		tokenReq.GrantType = models.REFRESH
	case "upgrade_scope":
		tokenReq.GrantType = models.UPSCOPE
	case "":
		response.Failed(w, exception.NewBadRequest("grant_type missed"))
		return
	default:
		response.Failed(w, exception.NewBadRequest("grant_type suport list [authorization_code, implicit, password, client_credentials, refresh_token, upgrade_scope]"))
		return
	}

	t, err := global.Store.IssueToken(tokenReq)
	if err != nil {
		response.Failed(w, err)
		return
	}

	w.Header().Set("x-oauth-access-token", t.AccessToken)
	w.Header().Set("x-oauth-refresh-token", t.RefreshToken)

	response.Success(w, http.StatusCreated, t)
	return

}

// ValidateToken use to validate token information
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	var (
		clientID     string
		clientSecret string
	)

	qs := r.URL.Query()
	ps := context.GetParamsFromContext(r)

	contentT := strings.Split(r.Header.Get("content-type"), ";")
	if len(contentT) == 0 {
		response.Failed(w, exception.NewBadRequest("content-type missed, your must be choice one [application/json, application/x-www-form-urlencoded]"))
		return
	}

	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		response.Failed(w, exception.NewForbidden("client_credentials missed"))
		return
	}

	vreq := &store.ValidateTokenReq{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AccessToken:  ps.ByName("tk"),
		FeatureName:  qs.Get("feature"),
	}

	// 2. validate token
	t, err := global.Store.ValidateTokenWithClient(vreq)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, t)
	return
}

// RevolkToken use to revolk token
func RevolkToken(w http.ResponseWriter, r *http.Request) {
	var (
		clientID     string
		clientSecret string
	)

	ps := context.GetParamsFromContext(r)
	tk := ps.ByName("tk")

	contentT := strings.Split(r.Header.Get("content-type"), ";")
	if len(contentT) == 0 {
		response.Failed(w, exception.NewBadRequest("content-type missed, your must be choice one [application/json, application/x-www-form-urlencoded]"))
		return
	}

	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		response.Failed(w, exception.NewForbidden("client_credentials missed"))
		return
	}

	revokeReq := &store.RevokeReq{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AccessToken:  tk,
	}

	if err := global.Store.RevokeToken(revokeReq); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}

// IssueCode use to issue oauth2 auth code
func IssueCode(w http.ResponseWriter, r *http.Request) {

}
