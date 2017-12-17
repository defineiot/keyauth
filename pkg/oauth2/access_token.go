package oauth2

import (
	"time"

	"openauth/api/exception"
	"openauth/storage/token"
)

// TokenRequest use to request access token
type TokenRequest struct {
	Scope        *token.Scope
	GrantType    token.GrantType
	clientID     string
	clientSecret string
}

// IssueToken use to issue access token
func (c *Controller) IssueToken(req *TokenRequest) (*token.Token, error) {
	cli, err := c.as.GetClient(req.clientID)
	if err != nil {
		return nil, err
	}
	if req.GrantType != token.IMPLICIT {
		if req.clientSecret != cli.ClientSecret {
			return nil, exception.NewUnauthorized("unauthorized_client")
		}
	}

	switch req.GrantType {
	case token.AUTHCODE:
	case token.IMPLICIT:
	case token.PASSWORD:
	case token.CLIENT:
	case token.REFRESH:
	default:
		return nil, exception.NewBadRequest("invalid_grant")
	}

	return nil, nil
}

// ValidateToken use to valdiate token
func (c *Controller) ValidateToken() {

}

// RevolkToken refresh token
func (c *Controller) RevolkToken() {

}

// issueTokenByAuthCode implement Authorization Code Grant
// https://tools.ietf.org/html/rfc6749#section-4.1.3
func (c *Controller) issueTokenByAuthCode(clientID, clientSecret, code, redirectURI string) (*token.Token, error) {
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
func (c *Controller) issuteTokenByClient(clientID, clientSecret, scope string) (*token.Token, error) {
	return nil, nil
}

// issueTokenByRefresh implement Refreshing an Access Token
// https://tools.ietf.org/html/rfc6749#section-6
func (c *Controller) issueTokenByRefresh(cientID, clientSecret, refreshToken string) (*token.Token, error) {
	return nil, nil
}
