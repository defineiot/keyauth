package service

import (
	"openauth/store/service"
)

// CreateService service
func (c *Controller) CreateService(name, description string) (*service.Service, error) {
	svr, err := c.ss.SaveService(name, description)
	if err != nil {
		return nil, err
	}
	return svr, nil
}

// DeleteService delete service
func (c *Controller) DeleteService(sid string) error {
	return nil
}

// ListService list all service
func (c *Controller) ListService() ([]*service.Service, error) {
	return nil, nil
}

// GetService get on service by id
func (c *Controller) GetService(sid string) (*service.Service, error) {
	return nil, nil
}
