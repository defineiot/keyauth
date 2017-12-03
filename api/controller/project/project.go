package project

import (
	"database/sql"
	"sync"

	"openauth/api/exception"
	"openauth/api/logger"
	"openauth/pkg/project"
	"openauth/pkg/user"
)

var (
	controller *Controller
	once       sync.Once
)

// GetController use to use an domain controller
func GetController() (*Controller, error) {
	if controller == nil {
		return nil, exception.NewInternalServerError("project controller isn't initial")
	}

	return controller, nil
}

// InitController use to initial an domain controller instance
func InitController(db *sql.DB, logger logger.OpenAuthLogger, pm project.Manager) error {
	once.Do(func() {
		controller = &Controller{db: db, logger: logger, pm: pm}
	})

	return nil
}

// Controller is domain controller
type Controller struct {
	db     *sql.DB
	logger logger.OpenAuthLogger
	pm     project.Manager
}

// CreateProject use to create domain
func (c *Controller) CreateProject(domainID, name, description string, cred user.Credential) (*project.Project, error) {
	proj, err := c.pm.CreateProject(domainID, name, description, true)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

// ListProject list an domain's project
func (c *Controller) ListProject(domainID string) ([]*project.Project, error) {
	projects, err := c.pm.ListDomainProjects(domainID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// GetProject use to get one project
func (c *Controller) GetProject(id string, cred user.Credential) (*project.Project, error) {
	proj, err := c.pm.GetProject(id)
	if err != nil {
		return nil, err
	}

	// TODO: check the project is for this user

	return proj, nil
}

// UpdateProject use to update one project
func (c *Controller) UpdateProject(cred user.Credential) (*project.Project, error) {
	return nil, nil
}

// DestroyProject use to delete one project
func (c *Controller) DestroyProject(id string, cred user.Credential) error {
	// TODO: check the projcet is for this user

	if err := c.pm.DeleteProject(id); err != nil {
		return err
	}

	return nil
}
