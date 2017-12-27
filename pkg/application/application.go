package application

import (
	"openauth/api/exception"
	"openauth/store/application"
)

// RegisterApplication use to regist application
func (c *Controller) RegisterApplication(userID, name, redirectURI, clientType, description, website string) (*application.Application, error) {
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
