package user_test

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"openauth/pkg/user"
	domainmock "openauth/store/domain/mock"
	"openauth/store/project"
	projectmock "openauth/store/project/mock"
	userstr "openauth/store/user"
	usermock "openauth/store/user/mock"
)

func NewUserController() *user.Controller {
	log := logrus.New()

	ds := new(domainmock.DomainStore)
	us := new(usermock.UserStore)
	ps := new(projectmock.ProjectStore)

	mockp := project.Project{
		ID:       "validated-id",
		Name:     "validated-name",
		Enabled:  true,
		DomainID: "validated-domain-id",
	}
	ps.GetProjectFn = func(id string) (*project.Project, error) {
		if id == "validated-project-id" {
			return &mockp, nil
		}
		return nil, fmt.Errorf("project %s not found", id)
	}
	ps.CheckProjectIsExistByIDFn = func(id string) (bool, error) {
		if id == "validated-project-id" {
			return true, nil
		}
		return false, nil
	}

	ds.CheckDomainIsExistByIDFn = func(domainID string) (bool, error) {
		if domainID == "validated-domain-id" {
			return true, nil
		}
		return false, nil
	}

	mocku := userstr.User{
		ID:               "validated-user-id",
		Name:             "validated-name",
		Enabled:          true,
		ExpireActiveDays: 4096,
		DomainID:         "validated-domain-id",
		DefaultProjectID: "validated-project-id",
		Password: &userstr.Password{
			ID:       1,
			Password: "validated-password",
			ExpireAt: 4096,
		},
	}
	us.CreateUserFn = func(domainID, name, password string, enabled bool, userExpires, passExpires int) (*userstr.User, error) {
		u := userstr.User{
			ID:               "validated-user-id",
			Name:             name,
			Enabled:          enabled,
			ExpireActiveDays: int64(userExpires),
			DomainID:         domainID,
			Password: &userstr.Password{
				ID:       1,
				Password: password,
				ExpireAt: int64(passExpires),
			},
		}
		return &u, nil
	}
	us.GetUserByIDFn = func(userID string) (*userstr.User, error) {
		if userID == "validated-user-id" {
			return &mocku, nil
		}
		return nil, fmt.Errorf("user %s not found", userID)
	}
	us.ListUserFn = func(domainID string) ([]*userstr.User, error) {
		if domainID == "validated-domain-id" {
			return []*userstr.User{&mocku}, nil
		}
		return []*userstr.User{}, nil
	}
	us.SetDefaultProjectFn = func(userID, projectID string) error {
		if userID == "validated-user-id" {
			return nil
		}
		return fmt.Errorf("user %s not found", userID)
	}
	us.ListUserProjectsFn = func(userID string) ([]string, error) {
		if userID == "validated-user-id" {
			return []string{"validated-project-id"}, nil
		}

		return nil, fmt.Errorf("user %s not found", userID)
	}
	us.AddProjectsToUserFn = func(userID string, projectIDs ...string) error {
		if userID == "validated-user-id" {
			return nil
		}

		return fmt.Errorf("user %s not found", userID)
	}

	us.RemoveProjectsFromUserFn = func(userID string, projectIDs ...string) error {
		if userID == "validated-user-id" {
			return nil
		}

		return fmt.Errorf("user %s not found", userID)
	}

	return user.NewController(log, us, ds, ps)
}
