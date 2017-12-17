package oauth2

import (
	"time"

	"openauth/api/exception"
	"openauth/storage/token"
)

// TokenRequest use to request access token
type TokenRequest struct {
	Scope               *token.Scope
	GrantType           token.GrantType
	ClientID            string
	ClientSecret        string
	AuthorizationHeader string
	DomainName          string
	Username            string
	Password            string
	Code                string
	RedirectURI         string
	RefreshToken        string
}

// Validate validate the request
func (t *TokenRequest) Validate() error {
	if t.Scope == nil {
		return exception.NewBadRequest("scope must'nt be nil")
	}
	if t.Scope.DomainID == "" && t.Scope.ProjectID == "" {
		return exception.NewBadRequest("scope's domain id or project id must choice one")
	}

	switch t.GrantType {
	case token.AUTHCODE:
		if t.Code == "" {
			return exception.NewBadRequest("if use %s grant type, code is needed", t.GrantType)
		}
		if t.RedirectURI == "" {
			return exception.NewBadRequest("if use %s grant type, redirect uri is need", t.GrantType)
		}
		goto CHECK_CLIENT

	case token.IMPLICIT:
		if t.ClientID == "" {
			return exception.NewBadRequest("if use %s grant type, client id is needed", t.GrantType)
		}
		if t.RedirectURI == "" {
			return exception.NewBadRequest("if use %s grant type, redirect uri is need", t.GrantType)
		}

	case token.PASSWORD:
		if t.DomainName == "" || t.Username == "" || t.Password == "" {
			return exception.NewBadRequest("if use %s grant type, domainname, username, password is needed", t.GrantType)
		}
		goto CHECK_CLIENT

	case token.CLIENT:
		goto CHECK_CLIENT

	case token.REFRESH:
		if t.RefreshToken == "" {
			return exception.NewBadRequest("if use %s grant type, refresh token is needed", t.GrantType)
		}
		goto CHECK_CLIENT

	default:
		return exception.NewBadRequest("invalid_grant")
	}

CHECK_CLIENT:
	if t.ClientID == "" || t.ClientSecret == "" {
		return exception.NewBadRequest("if use %s grant type, client id and client secret is needed", t.GrantType)
	}

	return nil
}

// IssueToken use to issue access token
func (c *Controller) IssueToken(req *TokenRequest) (*token.Token, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	cli, err := c.as.GetClient(req.ClientID)
	if err != nil {
		return nil, err
	}

	if req.GrantType != token.IMPLICIT {
		if req.ClientSecret != cli.ClientSecret {
			return nil, exception.NewUnauthorized("unauthorized_client")
		}
	}

	var t *token.Token
	switch req.GrantType {
	case token.AUTHCODE:
		t, err = c.issueTokenByAuthCode(req.ClientID, req.Code, req.RedirectURI)
		goto DEAL_ERROR
	case token.IMPLICIT:
		t, err = c.issueTokenByImplicit(req.ClientID, req.RedirectURI)
		goto DEAL_ERROR
	case token.PASSWORD:
		t, err = c.issueTokenByPassword(req.Scope, cli.ClientID, req.DomainName, req.Username, req.Password)
		goto DEAL_ERROR
	case token.CLIENT:
		t, err = c.issuteTokenByClient(req.ClientID, req.Scope)
		goto DEAL_ERROR
	case token.REFRESH:
		t, err = c.issueTokenByRefresh(req.ClientID, req.RefreshToken)
		goto DEAL_ERROR
	default:
		return nil, exception.NewBadRequest("invalid_grant")
	}

DEAL_ERROR:
	if err != nil {
		return nil, err
	}

	return t, nil
}

// ValidateToken use to valdiate token
func (c *Controller) ValidateToken() {

}

// RevolkToken refresh token
func (c *Controller) RevolkToken() {

}

// issueTokenByAuthCode implement Authorization Code Grant
// https://tools.ietf.org/html/rfc6749#section-4.1.3
func (c *Controller) issueTokenByAuthCode(clientID, code, redirectURI string) (*token.Token, error) {
	return nil, nil
}

// issueTokenByImplicit implement Implicit Grant
// https://tools.ietf.org/html/rfc6749#section-4.2
func (c *Controller) issueTokenByImplicit(clientID, redirectURI string) (*token.Token, error) {
	return nil, nil
}

// issueTokenByPassword implement Resource Owner Password Credentials Grant
// https://tools.ietf.org/html/rfc6749#section-4.3
func (c *Controller) issueTokenByPassword(scope *token.Scope, clientID, domainname, username, password string) (*token.Token, error) {
	dm, err := c.ds.GetDomainByName(domainname)
	if err != nil {
		return nil, err
	}

	uid, err := c.us.ValidateUser(dm.ID, username, password)
	if err != nil {
		return nil, err
	}
	if uid == "" {
		return nil, exception.NewUnauthorized("user password error")
	}

	scope.DomainID = dm.ID

	t := new(token.Token)
	t.Scope = scope
	t.CreatedAt = time.Now().Unix()
	t.ExpiresIn = c.expiresIn
	t.GrantType = token.PASSWORD
	t.TokenType = c.tokenType
	t.UserID = uid
	t.ClientID = clientID

	switch c.tokenType {
	case "bearer":
		t.AccessToken = makeBearerToken(24)
		t.RefreshToken = makeBearerToken(32)
	case "jwt":
	default:
		return nil, exception.NewInternalServerError("unknow token type, %s", c.tokenType)
	}

	retToken, err := c.ts.SaveToken(t)
	if err != nil {
		return nil, err
	}

	return retToken, nil
}

// issuteTokenByClient implement Client Credentials Grant
// https://tools.ietf.org/html/rfc6749#section-4.4.2
func (c *Controller) issuteTokenByClient(clientID string, scope *token.Scope) (*token.Token, error) {
	return nil, nil
}

// issueTokenByRefresh implement Refreshing an Access Token
// https://tools.ietf.org/html/rfc6749#section-6
func (c *Controller) issueTokenByRefresh(cientID, refreshToken string) (*token.Token, error) {
	return nil, nil
}
