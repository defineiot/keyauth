package domain

import (
	"sync"

	"openauth/api/exception"
	"openauth/api/logger"
	"openauth/store/domain"
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
func InitController(ds domain.Store, log logger.OpenAuthLogger) {
	once.Do(func() {
		controller = &Controller{ds: ds, log: log}
		controller.log.Debug("initial domain controller successful")
	})
	controller.log.Info("domain controller already initialed")
}

// Controller is domain pkg
type Controller struct {
	ds  domain.Store
	log logger.OpenAuthLogger
}

// CreateDomain use to create domain
func (c *Controller) CreateDomain(name, description, displayName string, enabled bool) (*domain.Domain, error) {
	dom, err := c.ds.CreateDomain(name, description, displayName, enabled)
	if err != nil {
		return nil, err
	}

	return dom, nil
}

// ListDomain use to list all domains
func (c *Controller) ListDomain(pageNumber, pageSize int64) ([]*domain.Domain, int64, error) {
	doms, totalPage, err := c.ds.ListDomain(pageNumber, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return doms, totalPage, nil
}

// GetDomain use to get an domain
func (c *Controller) GetDomain(domainID string) (*domain.Domain, error) {
	dom, err := c.ds.GetDomain(domainID)
	if err != nil {
		return nil, err
	}

	return dom, nil

}

// UpdateDomain use to update an domain
func (c *Controller) UpdateDomain() {

}

// DestroyDomain use to delete an domain
func (c *Controller) DestroyDomain(domainID string) error {
	if err := c.ds.DeleteDomain(domainID); err != nil {
		return err
	}

	return nil
}
