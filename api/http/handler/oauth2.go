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
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	tokenReq := new(oauth2.TokenRequest)
	tokenReq.ClientID = val.Get("client_id").ToString()
	tokenReq.ClientSecret = val.Get("client_secret").ToString()
	switch val.Get("grant_type").ToString() {
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
	default:
		response.Failed(w, exception.NewBadRequest("grant_type suport list [authorization_code, implicit, password, client_credentials, refresh_token"))
		return
	}

	tokenReq.DomainName = val.Get("domian_name").ToString()
	tokenReq.Username = val.Get("username").ToString()
	tokenReq.Password = val.Get("password").ToString()

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
