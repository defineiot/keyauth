package handler

import (
	"net/http"
	"strings"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/http/request"
	"github.com/defineiot/keyauth/api/http/response"
	"github.com/defineiot/keyauth/dao/token"
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
		tokenReq.GrantType = token.AUTHCODE
	case "implicit":
		tokenReq.GrantType = token.IMPLICIT
	case "password":
		tokenReq.GrantType = token.PASSWORD
	case "client_credentials":
		tokenReq.GrantType = token.CLIENT
	case "refresh_token":
		tokenReq.GrantType = token.REFRESH
	case "upgrade_scope":
		tokenReq.GrantType = token.UPSCOPE
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
		cleintSecret string
	)

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

	// 从认证头中获取token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		response.Failed(w, exception.NewUnauthorized("Authorization missed in header"))
		return
	}
	headerSlice := strings.Split(authHeader, " ")
	if len(headerSlice) != 2 {
		response.Failed(w, exception.NewUnauthorized("Authorization header value is not validated, must be: {token_type} {token}"))
		return
	}
	access := headerSlice[1]

	// 2. validate token
	t, err := global.Store.ValidateToken(access)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// 3. validate client token only for registry service
	qs := r.URL.Query()
	serviceName := qs.Get("sn")
	featureName := qs.Get("feature")

	switch t.GrantType {
	case token.CLIENT:
		if featureName != "" && featureName != "RegistryServiceFeatures" {
			response.Failed(w, exception.NewForbidden("client_credentials only can acess RegistryServiceFeatures"))
			return
		}
	case token.PASSWORD, token.UPSCOPE, token.REFRESH:
		if featureName != "" {
			if featureName == "RegistryServiceFeatures" {
				response.Failed(w, exception.NewForbidden("RegistryServiceFeatures only for client_credentials access"))
				return
			}
			if serviceName == "" {
				response.Failed(w, exception.NewForbidden("sn and feature needed when check permission"))
				return
			}
		}

		if serviceName != "" && featureName == "" {
			response.Failed(w, exception.NewForbidden("sn and feature needed when check permission"))
			return
		}

		// 4. check user role
		u, err := global.Store.GetUser(t.DomainID, t.UserID)
		if err != nil {
			response.Failed(w, err)
			return
		}
		for i := range u.Roles {
			switch u.Roles[i].Name {
			case "system_admin":
				t.IsSystemAdmin = true
			case "domain_admin":
				t.IsDomainAdmin = true
			default:
			}
		}

		// 4. if your has feature check user permission
		if t.UserID != "" && t.GrantType != token.CLIENT && featureName != "" && serviceName != "" {
			var hasPerm bool

			ok, err := global.Store.CheckServiceHasFeature(serviceName, featureName)
			if err != nil {
				response.Failed(w, err)
				return
			}
			if !ok {
				response.Failed(w, exception.NewBadRequest("not find service feature (%s:%s)", serviceName, featureName))
				return
			}

			for i := range u.Roles {
				rn := u.Roles[i].Name
				if rn == "system_admin" {
					hasPerm = true
					break
				}

				role, err := global.Store.GetRole(rn)
				if err != nil {
					response.Failed(w, exception.NewUnauthorized(err.Error()))
					return
				}

				for _, f := range role.Features {
					if f.Name == featureName {
						hasPerm = true
					}
				}
			}

			if !hasPerm {
				response.Failed(w, exception.NewForbidden("user: %s has no permisson for access feature: %s", u.Account, featureName))
				return
			}
		}
	default:
		response.Failed(w, exception.NewBadRequest("other grant type not support"))
		return
	}

	response.Success(w, http.StatusOK, t)
	return
}

// RevolkToken use to revolk token
func RevolkToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		response.Failed(w, exception.NewUnauthorized("Authorization missed in header"))
		return
	}

	headerSlice := strings.Split(authHeader, " ")
	if len(headerSlice) != 2 {
		response.Failed(w, exception.NewUnauthorized("Authorization header value is not validated, must be: {token_type} {token}"))
		return
	}

	access := headerSlice[1]

	if err := global.Store.RevokeToken(access); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}

// IssueCode use to issue oauth2 auth code
func IssueCode(w http.ResponseWriter, r *http.Request) {

}
