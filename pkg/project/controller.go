package project

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
		return nil, exception.NewInternalServerError("project controller is not initial")
	}
	return controller, nil
}

// InitController use to initial controller
func InitController(logger logger.OpenAuthLogger, ds domain.Store, ps project.Store, us user.Store) {
	once.Do(func() {
		controller = &Controller{log: logger, ds: ds, ps: ps, us: us}
		controller.log.Debug("initial project controller successful")
	})
	controller.log.Info("project controller already initialed")
}

// Controller is domain pkg
type Controller struct {
	log logger.OpenAuthLogger
	ps  project.Store
	ds  domain.Store
	us  user.Store
}
