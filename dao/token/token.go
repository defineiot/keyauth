package token

import (
	"time"

	"github.com/defineiot/keyauth/internal/exception"
)

// GrantType is the type for OAuth2 param `grant_type`
type GrantType string

const (
	// oauth2 Authorization Grant: https://tools.ietf.org/html/rfc6749#section-1.3

	// AUTHCODE oauth2 Authorization Code Grant
	AUTHCODE GrantType = "authorization_code"
	// IMPLICIT oauth2 Implicit Grant
	IMPLICIT GrantType = "implicit"
	// PASSWORD oauth2 Resource Owner Password Credentials Grant
	PASSWORD GrantType = "password"
	// CLIENT oauth2 Client Credentials Grant
	CLIENT GrantType = "client_credentials"
	// REFRESH oauth2 Refreshing an Access Token
	REFRESH GrantType = "refresh_token"
	// UPSCOPE is an custom grant for use unscope token acquire scope token
	UPSCOPE GrantType = "upgrade_scope"
)

// Code is oauth2 auth code https://tools.ietf.org/html/rfc6749#section-4.1.2
type Code struct {
	Code  string
	State string
}

// Token is user's access resource token
type Token struct {
	UserID        string    `json:"user_id,omitempty"`
	DomainID      string    `json:"domain_id,omitempty"`
	ServiceName   string    `json:"service_name,omitempty"`
	ClientID      string    `json:"client_id"`
	GrantType     GrantType `json:"grant_type"`
	AccessToken   string    `json:"access_token"`
	RefreshToken  string    `json:"refresh_token,omitempty"`
	TokenType     string    `json:"token_type"`
	CreatedAt     int64     `json:"create_at"`
	ExpiresIn     int64     `json:"expires_in"`
	IsSystemAdmin bool      `json:"is_system_admin,omitempty"`
	Roles         []string  `json:"roles"`
	IsDomainAdmin bool      `json:"is_domain_admin,omitempty"`
	Scope         *Scope    `json:"scope"`
}

// Scope token scope detail https://tools.ietf.org/html/rfc6749#section-3.3
// here use struct, but when http headler while use rfc format
type Scope struct {
	AvaliableProjects []string `json:"available_projects"`
	WorkProject       string   `json:"current_project,omitempty"`
	WorkDomain        string   `json:"current_domain,omitempty"`
}

// Store is auth service
type Store interface {
	StoreReader
	StoreWriter
	Close() error
}

// StoreReader read information from store
type StoreReader interface {
	GetToken(accessToken string) (*Token, error)
	GetTokenByRefresh(refreshToken string) (*Token, error)
}

// StoreWriter write information to store
type StoreWriter interface {
	SaveToken(t *Token) (*Token, error)
	DeleteTokenByRefresh(refreshToken string) error
	SetTokenProject(accessToken, projectID string) error
	DeleteToken(accessToken string) error
}

// Validate use to validate token to save
func (t *Token) validateSave() error {
	if t.ClientID == "" || t.UserID == "" {
		return exception.NewBadRequest("token's client_id or user_id is missed")
	}
	if t.AccessToken == "" {
		return exception.NewInternalServerError("token's access token must'nt be \"\"")
	}
	if t.Scope == nil {
		return exception.NewBadRequest("token's scope must'nt be null")
	}
	if t.DomainID == "" && t.Scope.WorkProject == "" {
		return exception.NewBadRequest("token's scope domain or project must choice one")
	}
	if t.TokenType != "bearer" && t.TokenType != "jwt" {
		return exception.NewInternalServerError("token's type must one of bearer or jwt")
	}
	if t.GrantType != AUTHCODE && t.GrantType != IMPLICIT && t.GrantType != PASSWORD && t.GrantType != CLIENT && t.GrantType != REFRESH {
		return exception.NewBadRequest("grant_type must one of authorization_code,implicit,password,client_credentials,refresh_token")
	}

	if t.CreatedAt == 0 {
		t.CreatedAt = time.Now().Unix()
	}
	if t.ExpiresIn == 0 {
		t.ExpiresIn = 3600
	}

	return nil
}

// IsExpired use to validate the token is expired
func (t *Token) IsExpired() (bool, int64) {
	now := time.Now().Unix()
	allow := t.CreatedAt + t.ExpiresIn

	if now < allow {
		return true, allow - now
	}

	return false, 0
}
