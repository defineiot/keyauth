package oauth2

import (
	"openauth/api/logger"
	"openauth/store/application"
	"openauth/store/domain"
	"openauth/store/token"
	"openauth/store/user"
)

// NewController use to init controller
func NewController(ts token.Store, us user.Store, ds domain.Store, as application.Store, log logger.OpenAuthLogger, tokenType string, expiresIn int64) *Controller {
	return &Controller{ts: ts, us: us, ds: ds, as: as, log: log, tokenType: tokenType, expiresIn: expiresIn}
}

// Controller is domain pkg
type Controller struct {
	ts        token.Store
	us        user.Store
	ds        domain.Store
	as        application.Store
	log       logger.OpenAuthLogger
	tokenType string
	expiresIn int64
}
