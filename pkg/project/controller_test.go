package project_test

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"openauth/pkg/project"
	domainmock "openauth/store/domain/mock"
	projectstr "openauth/store/project"
	projectmock "openauth/store/project/mock"
	userstr "openauth/store/user"
	usermock "openauth/store/user/mock"
)

func NewProjectController() *project.Controller {
	log := logrus.New()
	ds := new(domainmock.DomainStore)
	us := new(usermock.UserStore)
	ps := new(projectmock.ProjectStore)

	ds.CheckDomainIsExistByIDFn = func(domainID string) (bool, error) {
		if domainID == "validated-domain-id" {
			return true, nil
		}
		return false, nil
	}

	p := projectstr.Project{
		ID:       "validated-id",
		Name:     "validated-name",
		Enabled:  true,
		DomainID: "validated-domain-id",
	}
	ps.CreateProjectFn = func(domainID, name, description string, enabled bool) (*projectstr.Project, error) {
		return &p, nil
	}
	ps.ListDomainProjectsFn = func(domainID string) ([]*projectstr.Project, error) {
		return []*projectstr.Project{&p}, nil
	}
	ps.GetProjectFn = func(id string) (*projectstr.Project, error) {
		if id == "validated-project-id" {
			return &p, nil
		}
		return nil, fmt.Errorf("project %s not found", id)
	}
	ps.DeleteProjectFn = func(id string) error {
		if id == "validated-id" {
			return nil
		}
		return fmt.Errorf("project %s not found", id)
	}
	ps.AddUsersToProjectFn = func(projectID string, userIDs ...string) error {
		return nil
	}
	ps.RemoveUsersFromProjectFn = func(projectID string, userIDs ...string) error {
		return nil
	}
	ps.ListProjectUsersFn = func(projectID string) ([]string, error) {
		if projectID == "validated-project-id" {
			return []string{"validated-user-id"}, nil
		}
		return nil, fmt.Errorf("project %s not found", projectID)
	}

	u := userstr.User{
		ID:       "validated-user-id",
		Name:     "validated-name",
		Enabled:  true,
		DomainID: "validated-domain-id",
		Password: &userstr.Password{
			ID:       1,
			Password: "validated-password",
			UserID:   "validated-user-id",
		},
	}
	us.GetUserByIDFn = func(userID string) (*userstr.User, error) {
		if userID == "validated-user-id" {
			return &u, nil
		}
		return nil, fmt.Errorf("user %s not found", userID)
	}

	usersvr := project.NewController(log, ds, ps, us)
	return usersvr
}
