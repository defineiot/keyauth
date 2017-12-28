package mock

import "openauth/store/user"

// UserStore for mock user store
type UserStore struct {
	CreateUserFn                  func(domainID, name, password string, enabled bool, userExpires, passExpires int) (*user.User, error)
	CreateUserInvoked             bool
	ListUserProjectsFn            func(userID string) ([]string, error)
	ListUserProjectsInvoked       bool
	SetDefaultProjectFn           func(userID, projectID string) error
	SetDefaultProjectInvoked      bool
	AddProjectsToUserFn           func(userID string, projectIDs ...string) error
	AddProjectsToUserInvoked      bool
	RemoveProjectsFromUserFn      func(userID string, projectIDs ...string) error
	RemoveProjectsFromUserInvoked bool
	GetUserByIDFn                 func(userID string) (*user.User, error)
	GetUserByIDInvoked            bool
	QueryPhoneFn                  func(userID string) ([]*user.Phone, error)
	QueryPhoneInvoked             bool
	QueryEmailFn                  func(userID string) ([]*user.Email, error)
	QueryEmailInvoked             bool
	QueryPasswordFn               func(userID string) (*user.Password, error)
	QueryPasswordInvoked          bool
	ListUserFn                    func(domainID string) ([]*user.User, error)
	ListUserInvoked               bool
	DeleteUserFn                  func(userID string) error
	DeleteUserInvoked             bool
	GetUserByNameFn               func(domainID, userName string) (*user.User, error)
	GetUserByNameInvoked          bool
	ValidateUserFn                func(domainID, userName, password string) (string, error)
	ValidateUserInvoked           bool
	CheckUserNameIsExistFn        func(domainID, userName string) (bool, error)
	CheckUserNameIsExistInvoked   bool
	CheckUserIsExistByIDFn        func(userID string) (bool, error)
	CheckUserIsExistByIDInvoked   bool
}

// Close mock
func (s *UserStore) Close() error {
	return nil
}

// CreateUser mock
func (s *UserStore) CreateUser(domainID, name, password string, enabled bool, userExpires, passExpires int) (*user.User, error) {
	s.CreateUserInvoked = true
	u, err := s.CreateUserFn(domainID, name, password, enabled, userExpires, passExpires)
	return u, err
}

// ListUserProjects mock
func (s *UserStore) ListUserProjects(userID string) ([]string, error) {
	s.ListUserProjectsInvoked = true
	pids, err := s.ListUserProjectsFn(userID)
	return pids, err
}

// SetDefaultProject mock
func (s *UserStore) SetDefaultProject(userID, projectID string) error {
	s.SetDefaultProjectInvoked = true
	err := s.SetDefaultProjectFn(userID, projectID)
	return err
}

// AddProjectsToUser mock
func (s *UserStore) AddProjectsToUser(userID string, projectIDs ...string) error {
	s.AddProjectsToUserInvoked = true
	err := s.AddProjectsToUserFn(userID, projectIDs...)
	return err
}

// RemoveProjectsFromUser mock
func (s *UserStore) RemoveProjectsFromUser(userID string, projectIDs ...string) error {
	s.RemoveProjectsFromUserInvoked = true
	err := s.RemoveProjectsFromUserFn(userID, projectIDs...)
	return err
}

// GetUserByID mock
func (s *UserStore) GetUserByID(userID string) (*user.User, error) {
	s.GetUserByIDInvoked = true
	u, err := s.GetUserByIDFn(userID)
	return u, err
}

// QueryPhone mock
func (s *UserStore) QueryPhone(userID string) ([]*user.Phone, error) {
	s.QueryPhoneInvoked = true
	phone, err := s.QueryPhoneFn(userID)
	return phone, err
}

// QueryEmail mock
func (s *UserStore) QueryEmail(userID string) ([]*user.Email, error) {
	s.QueryEmailInvoked = true
	email, err := s.QueryEmailFn(userID)
	return email, err
}

// QueryPassword mock
func (s *UserStore) QueryPassword(userID string) (*user.Password, error) {
	s.QueryPasswordInvoked = true
	pass, err := s.QueryPasswordFn(userID)
	return pass, err
}

// ListUser mock
func (s *UserStore) ListUser(domainID string) ([]*user.User, error) {
	s.ListUserInvoked = true
	users, err := s.ListUserFn(domainID)
	return users, err
}

// DeleteUser mock
func (s *UserStore) DeleteUser(userID string) error {
	s.DeleteUserInvoked = true
	err := s.DeleteUserFn(userID)
	return err
}

// GetUserByName mock
func (s *UserStore) GetUserByName(domainID, userName string) (*user.User, error) {
	s.GetUserByIDInvoked = true
	u, err := s.GetUserByNameFn(domainID, userName)
	return u, err
}

// ValidateUser mock
func (s *UserStore) ValidateUser(domainID, userName, password string) (string, error) {
	s.ValidateUserInvoked = true
	id, err := s.ValidateUserFn(domainID, userName, password)
	return id, err
}

// CheckUserNameIsExist mock
func (s *UserStore) CheckUserNameIsExist(domainID, userName string) (bool, error) {
	s.CheckUserNameIsExistInvoked = true
	ok, err := s.CheckUserNameIsExistFn(domainID, userName)
	return ok, err
}

// CheckUserIsExistByID mock
func (s *UserStore) CheckUserIsExistByID(userID string) (bool, error) {
	s.CheckUserIsExistByIDInvoked = true
	ok, err := s.CheckUserIsExistByIDFn(userID)
	return ok, err
}
