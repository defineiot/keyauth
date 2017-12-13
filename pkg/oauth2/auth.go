package oauth2

import (
	"openauth/pkg/application"
	"openauth/pkg/user"
)

// GrantType is the type for OAuth2 param `grant_type`
type GrantType string

const (
	// oauth2 Authorization Grant: https://tools.ietf.org/html/rfc6749#section-1.3
	AuthorizationCode                GrantType = "authorization_code"
	Implicit                         GrantType = "implicit"
	ResourceOwnerPasswordCredentials GrantType = "resource_owner_password_credentials"
	ClientCredentials                GrantType = "client_credentials"
	RefreshToken                     GrantType = "refresh_token"
)

// Request use to request token
type Request struct {
	user         *user.User
	app          *application.Application
	grantType    GrantType
	code         string
	codeVerifier string
	state        string
}

// Token is user's access resource token
type Token struct {
	UserID       string `json:"user_id"`
	GrantType    string `json:"grant_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	CreatedAt    int64  `json:"create_at"`
	ExpiresIn    int32  `json:"expires_in"`
	Scope        *Scope `json:"scope"`
}

// Scope token scope
type Scope struct {
	ProjectID string `json:"project_id,omitempty"`
	DomainID  string `json:"domain_id,omitempty"`
}

// Code https://tools.ietf.org/html/rfc6749#section-4.1.2
type Code struct {
	Code  string
	State string
}

// Service is auth service
type Service interface {
	IssueTokenWithProject(userID, projectID string) (*Token, error)
	IssueTokenWithDomain(userID, domainID string) (*Token, error)
	IssueTokenByCode(code string) (*Token, error)
	IssueAuthCode(app *application.Application) (*Code, error)
	ValidateToken(accessToken string) (*Token, error)
	RefreshToken(refreshToken string) (*Token, error)
	RevokeToken(accessToken string) error
}
