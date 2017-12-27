package application

import (
	"sync"

	"openauth/api/exception"
	"openauth/api/logger"
	"openauth/store/application"
	"openauth/store/user"
)

var (
	controller *Controller
	once       sync.Once
)

// GetController use to new an controller
func GetController() (*Controller, error) {
	if controller == nil {
		return nil, exception.NewInternalServerError("application controller not initial")
	}
	return controller, nil
}

// InitController use to initial controller
func InitController(logger logger.OpenAuthLogger, as application.Store, us user.Store) {
	once.Do(func() {
		controller = &Controller{log: logger, as: as, us: us}
		controller.log.Debug("initial application controller successful")
	})
	controller.log.Info("application controller already initialed")
}

// Controller is domain pkg
type Controller struct {
	log logger.OpenAuthLogger
	as  application.Store
	us  user.Store
}
