package user

import (
	"openauth/api/exception"
	"openauth/store/project"
	"openauth/store/user"
)

// CreateUser create user
func (c *Controller) CreateUser(domainID, userName, password, description string) (*user.User, error) {
	ok, err := c.ds.CheckDomainIsExistByID(domainID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("domain %s not exist", domainID)
	}

	u, err := c.us.CreateUser(domainID, userName, password, true, 4096, 4096)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetUser get on user
func (c *Controller) GetUser(userID string) (*user.User, error) {
	u, err := c.us.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// ListUser list user
func (c *Controller) ListUser(domainID string) ([]*user.User, error) {
	ok, err := c.ds.CheckDomainIsExistByID(domainID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("domain %s not exist", domainID)
	}

	users, err := c.us.ListUser(domainID)
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.DefaultProjectID != "" {
			dp, err := c.ps.GetProject(u.DefaultProjectID)
			if err != nil {
				return nil, exception.NewInternalServerError("get user %s project error, %s", u.Name, err)
			}

			u.DefaultProject = dp
		}
	}

	return users, nil
}

// DeleteUser delete user
func (c *Controller) DeleteUser(userID string) error {
	if err := c.us.DeleteUser(userID); err != nil {
		return err
	}
	return nil
}

// SetUserDefaultProject use to set default prject
func (c *Controller) SetUserDefaultProject(userID, projectID string) error {
	ok, err := c.ps.CheckProjectIsExistByID(projectID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("project %s not exist", projectID)
	}

	if err := c.us.SetDefaultProject(userID, projectID); err != nil {
		return err
	}

	return nil
}

// ListUserProjects list all projects
func (c *Controller) ListUserProjects(userID string) ([]*project.Project, error) {
	projectIDs, err := c.us.ListUserProjects(userID)
	if err != nil {
		return nil, err
	}

	projects := []*project.Project{}
	for _, pid := range projectIDs {
		pj, err := c.ps.GetProject(pid)
		if err != nil {
			return nil, err
		}
		projects = append(projects, pj)
	}

	return projects, nil
}

// AddProjectsToUser add projects
func (c *Controller) AddProjectsToUser(userID string, projectIDs ...string) error {
	u, err := c.us.GetUserByID(userID)
	if err != nil {
		return err
	}

	for _, pid := range projectIDs {
		p, err := c.ps.GetProject(pid)
		if err != nil {
			return err
		}

		if p.DomainID != u.DomainID {
			return exception.NewBadRequest("user %s and project %s not in one domain", userID, pid)
		}
	}

	if err := c.us.AddProjectsToUser(userID, projectIDs...); err != nil {
		return err
	}

	return nil
}

// RemoveProjectsFromUser remove projects
func (c *Controller) RemoveProjectsFromUser(userID string, projectIDs ...string) error {
	for _, pid := range projectIDs {
		ok, err := c.ps.CheckProjectIsExistByID(pid)
		if err != nil {
			return err
		}
		if !ok {
			return exception.NewBadRequest("project %s not exist", pid)
		}
	}

	if err := c.us.RemoveProjectsFromUser(userID, projectIDs...); err != nil {
		return err
	}
	return nil
}
