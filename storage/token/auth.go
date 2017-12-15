package token

import (
	"openauth/api/exception"
	"time"
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
	//REFRESH oauth2 Refreshing an Access Token
	REFRESH GrantType = "refresh_token"
)

// Code is oauth2 auth code https://tools.ietf.org/html/rfc6749#section-4.1.2
type Code struct {
	Code  string
	State string
}

// Token is user's access resource token
type Token struct {
	UserID       string    `json:"user_id"`
	ClientID     string    `json:"-"`
	GrantType    GrantType `json:"grant_type"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type"`
	CreatedAt    int64     `json:"create_at"`
	ExpiresIn    int32     `json:"expires_in"`
	Scope        *Scope    `json:"scope"`
}

// Scope token scope detail https://tools.ietf.org/html/rfc6749#section-3.3
// here use struct, but when http headler while use rfc format
type Scope struct {
	ProjectID string `json:"project_id,omitempty"`
	DomainID  string `json:"domain_id,omitempty"`
}

// Storage is auth service
type Storage interface {
	SaveToken(t *Token) (*Token, error)
	// IssueTokenWithProject(userID, projectID string) (*Token, error)
	// IssueTokenWithDomain(userID, domainID string) (*Token, error)
	// IssueTokenByCode(code string) (*Token, error)
	// IssueAuthCode(app *application.Application) (*Code, error)
	// ValidateToken(accessToken string) (*Token, error)
	// RefreshToken(refreshToken string) (*Token, error)
	// RevokeToken(accessToken string) error
}

// Validate use to validate token to save
func (t *Token) Validate() error {
	if t.ClientID == "" || t.UserID == "" {
		return exception.NewBadRequest("token's client_id or user_id is missed")
	}
	if t.AccessToken == "" {
		return exception.NewInternalServerError("token's access token must'nt be \"\"")
	}
	if t.Scope == nil {
		return exception.NewBadRequest("token's scope must'nt be null")
	}
	if t.Scope.DomainID == "" && t.Scope.ProjectID == "" {
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
