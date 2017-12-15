package oauth2

import (
	"sync"
	"time"

	"openauth/api/exception"
	"openauth/api/logger"
	"openauth/storage/application"
	"openauth/storage/domain"
	"openauth/storage/token"
	"openauth/storage/user"
)

var (
	controller *Controller
	once       sync.Once
)

// GetController use to new an controller
func GetController() (*Controller, error) {
	if controller == nil {
		return nil, exception.NewInternalServerError("domain controller not initial")
	}
	return controller, nil
}

// InitController use to init controller
func InitController(ts token.Storage, us user.Storage, ds domain.Storage, as application.Storage, log logger.OpenAuthLogger, tokenType string, expiresIn int32) {
	once.Do(func() {
		controller = &Controller{ts: ts, us: us, ds: ds, as: as, log: log, tokenType: tokenType, expiresIn: expiresIn}
		controller.log.Debug("initial token controller successful")
	})
	controller.log.Info("token contoller aready initialed")
}

// Controller is domain pkg
type Controller struct {
	ts        token.Storage
	us        user.Storage
	ds        domain.Storage
	as        application.Storage
	log       logger.OpenAuthLogger
	tokenType string
	expiresIn int32
}

// IssueToken use to issue access token
func (c *Controller) IssueToken(grantType token.GrantType, clientID string, clientSecret string) (*token.Token, error) {
	cli, err := c.as.GetClient(clientID)
	if err != nil {
		return nil, err
	}
	if grantType != token.IMPLICIT {
		if clientSecret != cli.ClientSecret {
			return nil, exception.NewUnauthorized("unauthorized_client")
		}
	}

	switch grantType {
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

// IssueCode use to issue auth code
func (c *Controller) IssueCode() {

}

// ValidateToken use to valdiate token
func (c *Controller) ValidateToken() {

}

// RefreshToken refresh token
func (c *Controller) RefreshToken() {

}

// Authorization Code Grant
// https://tools.ietf.org/html/rfc6749#section-4.1.3
func (c *Controller) issueTokenByAuthCode(clientID, clientSecret, code, redirectURI string) (*token.Token, error) {
	return nil, nil
}

// Implicit Grant
// https://tools.ietf.org/html/rfc6749#section-4.2
func (c *Controller) issueTokenByImplicit(clientID, redirectURI string) (*token.Token, error) {
	return nil, nil
}

// Resource Owner Password Credentials Grant
// https://tools.ietf.org/html/rfc6749#section-4.3
func (c *Controller) issueTokenByPassword(domainname, username, password, scope string) (*token.Token, error) {
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

	t := new(token.Token)
	t.CreatedAt = time.Now().Unix()
	t.ExpiresIn = c.expiresIn
	t.GrantType = token.PASSWORD
	t.Scope = &token.Scope{DomainID: dm.ID}
	t.TokenType = c.tokenType
	t.UserID = uid

	switch c.tokenType {
	case "bearer":
		t.AccessToken = makeBearerToken(24)
		t.RefreshToken = makeBearerToken(32)
	case "jwt":
	default:
		return nil, exception.NewInternalServerError("unknow token type, %s", c.tokenType)
	}

	if err := c.ts.SaveToken(t); err != nil {
		return nil, err
	}

	return t, nil
}

// Client Credentials Grant
// https://tools.ietf.org/html/rfc6749#section-4.4.2
func (c *Controller) issuteTokenByClient(clientID, clientSecret, scope string) (*token.Token, error) {
	return nil, nil
}

// Refreshing an Access Token
// https://tools.ietf.org/html/rfc6749#section-6
func (c *Controller) issueTokenByRefresh(cientID, clientSecret, refreshToken string) (*token.Token, error) {
	return nil, nil
}
