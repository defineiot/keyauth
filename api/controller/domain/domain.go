package domain

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/sirupsen/logrus"

	"openauth/pkg/domain"
	"openauth/pkg/user"
)

var (
	controller *Controller
	once       sync.Once
)

// GetController use to use an domain controller
func GetController() (*Controller, error) {
	if controller == nil {
		return nil, errors.New("domain controller isn't initial")
	}

	return controller, nil
}

// InitController use to initial an domain controller instance
func InitController(db *sql.DB, logger *logrus.Logger, domain domain.Manager) error {
	once.Do(func() {
		controller = &Controller{db: db, logger: logger, dm: domain}
	})

	return nil
}

// Controller is domain controller
type Controller struct {
	db     *sql.DB
	logger *logrus.Logger
	dm     domain.Manager
}

// CreateDomain use to create domain
func (c *Controller) CreateDomain(name, description, verbose string, cert user.Credential) (*domain.Domain, error) {
	return nil, nil
}

// ListDomain use to list all domains
func (c *Controller) ListDomain() {

}

// GetDomain use to get an domain
func (c *Controller) GetDomain() {

}

// UpdateDomain use to update an domain
func (c *Controller) UpdateDomain() {

}

// DestoryDomain use to delete an domain
func (c *Controller) DestoryDomain() {

}
