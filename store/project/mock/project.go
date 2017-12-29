package mock

import "openauth/store/project"

// ProjectStore mock
type ProjectStore struct {
	CreateProjectFn                func(domainID, name, description string, enabled bool) (*project.Project, error)
	CreateProjectInvoked           bool
	GetProjectFn                   func(id string) (*project.Project, error)
	GetProjectInvoked              bool
	ListDomainProjectsFn           func(domainID string) ([]*project.Project, error)
	ListDomainProjectsInvoked      bool
	UpdateProjectFn                func(id, name, description string) (*project.Project, error)
	UpdateProjectInvoked           bool
	DeleteProjectFn                func(id string) error
	DeleteProjectInvoked           bool
	CheckProjectIsExistByIDFn      func(id string) (bool, error)
	CheckProjectIsExistByIDInvoked bool
	ListProjectUsersFn             func(projectID string) ([]string, error)
	ListProjectUsersInvoked        bool
	AddUsersToProjectFn            func(projectID string, userIDs ...string) error
	AddUsersToProjectInvoked       bool
	RemoveUsersFromProjectFn       func(projectID string, userIDs ...string) error
	RemoveUsersFromProjectInvoked  bool
}

// Close mock
func (s *ProjectStore) Close() error {
	return nil
}

// CreateProject mock
func (s *ProjectStore) CreateProject(domainID, name, description string, enabled bool) (*project.Project, error) {
	s.CreateProjectInvoked = true
	p, err := s.CreateProjectFn(domainID, name, description, enabled)
	return p, err
}

// GetProject mock
func (s *ProjectStore) GetProject(id string) (*project.Project, error) {
	s.GetProjectInvoked = true
	p, err := s.GetProjectFn(id)
	return p, err
}

// ListDomainProjects mock
func (s *ProjectStore) ListDomainProjects(domainID string) ([]*project.Project, error) {
	s.ListDomainProjectsInvoked = true
	ds, err := s.ListDomainProjectsFn(domainID)
	return ds, err
}

// UpdateProject mock
func (s *ProjectStore) UpdateProject(id, name, description string) (*project.Project, error) {
	s.UpdateProjectInvoked = true
	p, err := s.UpdateProjectFn(id, name, description)
	return p, err
}

// DeleteProject mock
func (s *ProjectStore) DeleteProject(id string) error {
	s.DeleteProjectInvoked = true
	err := s.DeleteProjectFn(id)
	return err
}

// CheckProjectIsExistByID mock
func (s *ProjectStore) CheckProjectIsExistByID(id string) (bool, error) {
	s.CheckProjectIsExistByIDInvoked = true
	ok, err := s.CheckProjectIsExistByIDFn(id)
	return ok, err
}

// ListProjectUsers mock
func (s *ProjectStore) ListProjectUsers(projectID string) ([]string, error) {
	s.ListProjectUsersInvoked = true
	uids, err := s.ListProjectUsersFn(projectID)
	return uids, err
}

// AddUsersToProject mock
func (s *ProjectStore) AddUsersToProject(projectID string, userIDs ...string) error {
	s.AddUsersToProjectInvoked = true
	err := s.AddUsersToProjectFn(projectID, userIDs...)
	return err
}

// RemoveUsersFromProject mock
func (s *ProjectStore) RemoveUsersFromProject(projectID string, userIDs ...string) error {
	s.RemoveUsersFromProjectInvoked = true
	err := s.RemoveUsersFromProjectFn(projectID, userIDs...)
	return err
}
