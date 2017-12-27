package user

import (
	"sync"

	"openauth/api/exception"
	"openauth/api/logger"
	"openauth/store/domain"
	"openauth/store/project"
	"openauth/store/user"
)

var (
	controller *Controller
	once       sync.Once
)

// GetController use to new an controller
func GetController() (*Controller, error) {
	if controller == nil {
		return nil, exception.NewInternalServerError("user controller not initial")
	}
	return controller, nil
}

// InitController use to initial user controller
func InitController(log logger.OpenAuthLogger, us user.Store, ds domain.Store, ps project.Store) {
	once.Do(func() {
		controller = &Controller{ds: ds, ps: ps, us: us, log: log}
		controller.log.Debug("user controller initial successful")
	})
	controller.log.Info("user controller already initialed")
}

// Controller is domain pkg
type Controller struct {
	log logger.OpenAuthLogger
	ps  project.Store
	ds  domain.Store
	us  user.Store
}
