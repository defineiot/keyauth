package handler

import (
	"net/http"

	"openauth/api/exception"
	"openauth/api/http/request"
	"openauth/api/http/response"
	"openauth/pkg/oauth2"
	"openauth/storage/token"
)

// IssueToken use to issue access token
func IssueToken(w http.ResponseWriter, r *http.Request) {
	var grantType string
	tokenReq := new(oauth2.TokenRequest)

	switch r.Header.Get("content-type") {
	case "application/json":
		val, err := request.CheckObjectBody(r)
		if err != nil {
			response.Failed(w, err)
			return
		}
		tokenReq.ClientID = val.Get("client_id").ToString()
		tokenReq.ClientSecret = val.Get("client_secret").ToString()
		tokenReq.DomainName = val.Get("domain_name").ToString()
		tokenReq.Username = val.Get("username").ToString()
		tokenReq.Password = val.Get("password").ToString()
		grantType = val.Get("grant_type").ToString()
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			response.Failed(w, exception.NewBadRequest("parse x-www-form-urlencoded data error, %s", err))
			return
		}
		tokenReq.ClientID = r.FormValue("client_id")
		tokenReq.ClientSecret = r.FormValue("client_secret")
		tokenReq.DomainName = r.FormValue("domain_name")
		tokenReq.Username = r.FormValue("username")
		tokenReq.Password = r.FormValue("password")
		grantType = r.FormValue("grant_type")
	case "":
		response.Failed(w, exception.NewBadRequest("content-type missed, your must be choice one [application/json, application/x-www-form-urlencoded]"))
		return
	default:
		response.Failed(w, exception.NewBadRequest("content-type only support for application/json and application/x-www-form-urlencoded others don't supported"))
		return
	}

	clientID, clientSecret, ok := r.BasicAuth()
	if ok {
		tokenReq.ClientID = clientID
		tokenReq.ClientSecret = clientSecret
	}

	switch grantType {
	case "authorization_code":
		tokenReq.GrantType = token.AUTHCODE
	case "implicit":
		tokenReq.GrantType = token.IMPLICIT
	case "password":
		tokenReq.GrantType = token.PASSWORD
	case "client_credentials":
		tokenReq.GrantType = token.CLIENT
	case "refresh_token":
		tokenReq.GrantType = token.REFRESH
	case "":
		response.Failed(w, exception.NewBadRequest("grant_type missed"))
		return
	default:
		response.Failed(w, exception.NewBadRequest("grant_type suport list [authorization_code, implicit, password, client_credentials, refresh_token]"))
		return
	}

	t, err := authsrc.IssueToken(tokenReq)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, t)
	return

}

// ValidateToken use to validate token information
func ValidateToken(w http.ResponseWriter, r *http.Request) {

}

// RevolkToken use to revolk token
func RevolkToken(w http.ResponseWriter, r *http.Request) {

}

// IssueCode use to issue oauth2 auth code
func IssueCode(w http.ResponseWriter, r *http.Request) {

}
