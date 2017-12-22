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
func InitController(logger logger.OpenAuthLogger, ds domain.Store, ps project.Storage, us user.Storage) {
	once.Do(func() {
		controller = &Controller{log: logger, ds: ds, ps: ps, us: us}
		controller.log.Debug("initial project controller successful")
	})
	controller.log.Info("project controller already initialed")
}

// Controller is domain pkg
type Controller struct {
	log logger.OpenAuthLogger
	ps  project.Storage
	ds  domain.Store
	us  user.Storage
}

// CreateProject use to create domain
func (c *Controller) CreateProject(domainID, name, description string) (*project.Project, error) {
	ok, err := c.ds.CheckDomainIsExistByID(domainID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("domain %s not exist", domainID)
	}

	proj, err := c.ps.CreateProject(domainID, name, description, true)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

// ListProject list an domain's project
func (c *Controller) ListProject(domainID string) ([]*project.Project, error) {
	ok, err := c.ds.CheckDomainIsExistByID(domainID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("domain %s not exist", domainID)
	}

	projects, err := c.ps.ListDomainProjects(domainID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// GetProject use to get one project
func (c *Controller) GetProject(id string) (*project.Project, error) {
	proj, err := c.ps.GetProject(id)
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
func (c *Controller) DestroyProject(id string) error {
	// TODO: check the project is for this user

	if err := c.ps.DeleteProject(id); err != nil {
		return err
	}

	return nil
}

// ListProjectUsers use to list project's users
func (c *Controller) ListProjectUsers(projectID string) ([]*user.User, error) {
	userIDs, err := c.ps.ListProjectUsers(projectID)
	if err != nil {
		return nil, err
	}

	users := []*user.User{}
	for _, uid := range userIDs {
		u, err := c.us.GetUserByID(uid)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// AddUsersToProject add users
func (c *Controller) AddUsersToProject(projectID string, userIDs ...string) error {
	// user and project must be in one domain
	p, err := c.ps.GetProject(projectID)
	if err != nil {
		return err
	}

	for _, uid := range userIDs {
		u, err := c.us.GetUserByID(uid)
		if err != nil {
			return err
		}

		if p.DomainID != u.DomainID {
			return exception.NewBadRequest("user %s and project %s not in one domain", uid, projectID)
		}
	}

	// insert
	if err := c.ps.AddUsersToProject(projectID, userIDs...); err != nil {
		return err
	}
	return nil
}

// RemoveUsersFromProject remove users
func (c *Controller) RemoveUsersFromProject(projectID string, userIDs ...string) error {
	// user and project must be in one domain
	p, err := c.ps.GetProject(projectID)
	if err != nil {
		return err
	}

	for _, uid := range userIDs {
		u, err := c.us.GetUserByID(uid)
		if err != nil {
			return err
		}

		if p.DomainID != u.DomainID {
			return exception.NewBadRequest("user %s and project %s not in one domain", uid, projectID)
		}
	}

	// insert
	if err := c.ps.RemoveUsersFromProject(projectID, userIDs...); err != nil {
		return err
	}
	return nil
}
