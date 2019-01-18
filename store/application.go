package store

import (
	"github.com/defineiot/keyauth/dao/application"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateApplication todo
func (s *Store) CreateApplication(userID, name, redirectURI, description, website string) (*application.Application, error) {
	ok, err := s.user.CheckUserIsExistByID(userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("user %s not found", userID)
	}

	cli, err := s.client.CreateClient(client.Public, redirectURI)
	if err != nil {
		return nil, err
	}

	app, err := s.app.Registration(userID, name, description, website, cli.ID)
	if err != nil {
		if err := s.client.DeleteClient(cli.ID); err != nil {
			s.log.Error(err)
		}
		return nil, err
	}

	app.Client = cli

	return app, nil
}

// DeleteApplication todo
func (s *Store) DeleteApplication(id string) error {
	var err error

	cacheKey := "app_" + id

	app, err := s.app.GetApplication(id)
	if err != nil {
		return err
	}

	err = s.app.Unregistration(app.ID)
	if err != nil {
		return err
	}

	if app.ClientID != "" {
		err := s.client.DeleteClient(app.ClientID)
		if err != nil {
			return err
		}
	}

	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debugf("delete app from cache failed, key: %s", cacheKey)
		}
		s.log.Debugf("delete app from cache success, key: %s", cacheKey)
	}

	return nil
}

// ListUserApps todo
func (s *Store) ListUserApps(userID string) ([]*application.Application, error) {
	ok, err := s.user.CheckUserIsExistByID(userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("user %s not found", userID)
	}

	apps, err := s.app.ListApplications(userID)
	if err != nil {
		return nil, err
	}
	for _, app := range apps {
		if app.ClientID != "" {
			cli, err := s.client.GetClient(app.ClientID)
			if err != nil {
				return nil, err
			}
			app.Client = cli
		}
	}

	return apps, nil
}

// GetUserApp todo
func (s *Store) GetUserApp(appid string) (*application.Application, error) {
	var err error

	app := new(application.Application)
	cacheKey := "app_" + appid

	if s.isCache {
		if s.cache.Get(cacheKey, app) {
			s.log.Debugf("get app from cache key: %s", cacheKey)
			return app, nil
		}
		s.log.Debugf("get app from cache failed, key: %s", cacheKey)
	}

	app, err = s.app.GetApplication(appid)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, exception.NewBadRequest("app %s not found", appid)
	}

	if app.ClientID != "" {
		cli, err := s.client.GetClient(app.ClientID)
		if err != nil {
			return nil, err
		}
		app.Client = cli
	}

	if s.isCache {
		if !s.cache.Set(cacheKey, app, s.ttl) {
			s.log.Debugf("set app cache failed, key: %s", cacheKey)
		}
		s.log.Debugf("set app cache ok, key: %s", cacheKey)
	}

	return app, nil
}
