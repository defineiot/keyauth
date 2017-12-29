package oauth2_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"openauth/pkg/oauth2"
	"openauth/store/application"
	appmockstr "openauth/store/application/mock"
	"openauth/store/domain"
	domainrmockstr "openauth/store/domain/mock"
	"openauth/store/token"
	tokenmockstr "openauth/store/token/mock"
	usermockstr "openauth/store/user/mock"
)

func NewOAuth2Ctroller(expirted int64) *oauth2.Controller {
	log := logrus.New()

	ts := new(tokenmockstr.TokenStore)
	us := new(usermockstr.UserStore)
	ds := new(domainrmockstr.DomainStore)
	as := new(appmockstr.AppStore)

	t := token.Token{
		UserID:       "user01",
		ClientID:     "client01",
		GrantType:    token.PASSWORD,
		AccessToken:  "validated-token-string",
		RefreshToken: "refresh token",
		TokenType:    "bearer",
		CreatedAt:    time.Now().Unix(),
		ExpiresIn:    expirted,
	}
	ts.GetTokenFn = func(accessToken string) (*token.Token, error) {
		if accessToken == "validated-token-string" {
			return &t, nil
		}
		return nil, fmt.Errorf("token %s not found", accessToken)
	}
	ts.SaveTokenFn = func(t *token.Token) (*token.Token, error) {
		return t, nil
	}

	d := domain.Domain{
		ID:      "domain-id",
		Name:    "validated-domain",
		Enabled: true,
	}
	ds.GetDomainByNameFn = func(name string) (*domain.Domain, error) {
		if name == "validated-domain" {
			return &d, nil
		}
		return nil, fmt.Errorf("domain %s not found", name)
	}

	us.ValidateUserFn = func(domainID, userName, password string) (string, error) {
		if !(domainID == "validated-domain" && userName == "validated-user" && password == "validated-pass") {
			return "validated-user-id", nil
		}
		return "", errors.New("user not validated")
	}

	cli := application.Client{
		ClientID:     "validated-client-id",
		ClientSecret: "validated-client-secret",
		ClientType:   "public",
	}
	as.GetClientFn = func(clientID string) (*application.Client, error) {
		if clientID == "validated-client-id" {
			return &cli, nil
		}
		return nil, fmt.Errorf("client %s not found", clientID)
	}

	return oauth2.NewController(ts, us, ds, as, log, "bearer", 3600)
}
