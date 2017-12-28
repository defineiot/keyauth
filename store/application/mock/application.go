package mock

import "openauth/store/application"

// AppStore 代表一个application.Store的mock实现
type AppStore struct {
	RegistrationFn               func(userID, name, redirectURI, clientType, description, website string) (*application.Application, error)
	RegistrationInvoked          bool
	CheckAPPIsExistByIDFn        func(appID string) (bool, error)
	CheckAPPIsExistByIDInvoked   bool
	CheckAPPIsExistByNameFn      func(userID, name string) (bool, error)
	CheckAPPIsExistByNameInvoked bool
	UnregistrationFn             func(id string) error
	UnregistrationInvoked        bool
	GetUserAppsFn                func(userID string) ([]*application.Application, error)
	GetUserAppsInvoked           bool
	GetClientFn                  func(clientID string) (*application.Client, error)
	GetClientInvoked             bool
}

// Close mock
func (s *AppStore) Close() error {
	return nil
}

// Registration mock
func (s *AppStore) Registration(userID, name, redirectURI, clientType, description, website string) (*application.Application, error) {
	s.RegistrationInvoked = true
	app, err := s.RegistrationFn(userID, name, redirectURI, clientType, description, website)
	return app, err
}

// CheckAPPIsExistByID mock
func (s *AppStore) CheckAPPIsExistByID(appID string) (bool, error) {
	s.CheckAPPIsExistByIDInvoked = true
	ok, err := s.CheckAPPIsExistByIDFn(appID)
	return ok, err
}

// CheckAPPIsExistByName mock
func (s *AppStore) CheckAPPIsExistByName(userID, name string) (bool, error) {
	s.CheckAPPIsExistByNameInvoked = true
	ok, err := s.CheckAPPIsExistByNameFn(userID, name)
	return ok, err
}

// Unregistration mock
func (s *AppStore) Unregistration(id string) error {
	s.UnregistrationInvoked = true
	err := s.UnregistrationFn(id)
	return err
}

// GetUserApps mock
func (s *AppStore) GetUserApps(userID string) ([]*application.Application, error) {
	s.GetClientInvoked = true
	app, err := s.GetUserAppsFn(userID)
	return app, err
}

// GetClient mock
func (s *AppStore) GetClient(clientID string) (*application.Client, error) {
	s.GetClientInvoked = true
	app, err := s.GetClientFn(clientID)
	return app, err
}
