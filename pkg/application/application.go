package application

import (
	"openauth/api/exception"
	"openauth/api/logger"
	"openauth/storage/application"
	"openauth/storage/user"
)

// NewController use to new an controller
func NewController(logger logger.OpenAuthLogger, as application.Storage, us user.Storage) *Controller {
	return &Controller{logger: logger, as: as, us: us}
}

// Controller is domain pkg
type Controller struct {
	logger logger.OpenAuthLogger
	as     application.Storage
	us     user.Storage
}

// RegisteApplication use to regist application
func (c *Controller) RegisteApplication(userID, name, redirectURI, clientType, description, website string) (*application.Application, error) {
	ok, err := c.us.CheckUserIsExistByID(userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("user %s not exist", userID)
	}

	app, err := c.as.Registration(userID, name, redirectURI, clientType, description, website)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// UnregisteApplication delete application
func (c *Controller) UnregisteApplication(id string) error {
	if err := c.as.Unregistration(id); err != nil {
		return err
	}

	return nil
}

// GetUserApplications get user's applications
func (c *Controller) GetUserApplications(userID string) ([]*application.Application, error) {
	ok, err := c.us.CheckUserIsExistByID(userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("user %s not exist", userID)
	}

	apps, err := c.as.GetUserApps(userID)
	if err != nil {
		return nil, err
	}

	return apps, nil
}
