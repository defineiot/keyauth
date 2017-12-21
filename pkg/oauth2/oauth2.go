package oauth2

import (
	"sync"

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
func InitController(ts token.Storage, us user.Storage, ds domain.Store, as application.Storage, log logger.OpenAuthLogger, tokenType string, expiresIn int64) {
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
	ds        domain.Store
	as        application.Storage
	log       logger.OpenAuthLogger
	tokenType string
	expiresIn int64
}
