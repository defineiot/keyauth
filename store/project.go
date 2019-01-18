package store

import (
	"errors"
	"strings"

	"github.com/defineiot/keyauth/dao/project"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateProject use to create an project
func (s *Store) CreateProject(p *project.Project) error {
	if p.DomainID == "" {
		return exception.NewBadRequest("domainID or domainName required one")
	}

	// check domain exist
	if _, err := s.dao.Domain.GetDomainByID(p.DomainID); err != nil {
		return err
	}

	return s.dao.Project.CreateProject(p)
}

// ListDomainProjects list domain projects
func (s *Store) ListDomainProjects(domainID string) ([]*project.Project, error) {
	// check domain exist
	if _, err := s.dao.Domain.GetDomainByID(domainID); err != nil {
		return nil, err
	}

	return s.dao.Project.ListDomainProjects(domainID)
}

// ListUserProjects todo
func (s *Store) ListUserProjects(domainID, userID string) ([]*project.Project, error) {
	return s.dao.Project.ListUserProjects(domainID, userID)
}

// GetProject get one project
func (s *Store) GetProject(id string) (*project.Project, error) {
	var err error

	pro := new(project.Project)
	cacheKey := "project_" + id

	if s.isCache {
		if s.cache.Get(cacheKey, pro) {
			s.log.Debug("get project from cache key: %s", cacheKey)
			return pro, nil
		}
		s.log.Debug("get project from cache failed, key: %s", cacheKey)
	}

	pro, err = s.dao.Project.GetProjectByID(id)
	if err != nil {
		return nil, err
	}
	if pro == nil {
		return nil, exception.NewBadRequest("project %s not found", id)
	}

	if s.isCache {
		if !s.cache.Set(cacheKey, pro, s.ttl) {
			s.log.Debug("set project cache failed, key: %s", cacheKey)
		}
		s.log.Debug("set project cache ok, key: %s", cacheKey)
	}

	return pro, nil
}

// DeleteProjectByID project
func (s *Store) DeleteProjectByID(id string) error {
	var err error

	cacheKey := "project_" + id

	err = s.dao.Project.DeleteProjectByID(id)
	if err != nil {
		return err
	}

	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debug("delete project from cache failed, key: %s", cacheKey)
		}
		s.log.Debug("delete project from cache success, key: %s", cacheKey)
	}

	return nil
}

// DeleteProjectByName project
func (s *Store) DeleteProjectByName(projectName, domainID string) error {
	var err error

	cacheKey := "project_" + domainID + projectName

	err = s.dao.Project.DeleteProjectByName(projectName, domainID)
	if err != nil {
		return err
	}

	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debug("delete project from cache failed, key: %s", cacheKey)
		}
		s.log.Debug("delete project from cache success, key: %s", cacheKey)
	}

	return nil
}

// AddUsersToProject add user
func (s *Store) AddUsersToProject(accessToken, projectID string, userIDs ...string) error {
	if err := s.checkUserExist(userIDs...); err != nil {
		return err
	}

	err := s.dao.Project.AddUsersToProject(projectID, userIDs...)

	cacheKey := "token_" + accessToken
	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debug("delete token from cache failed, key: %s", cacheKey)
		}
		s.log.Debug("delete token from cache success, key: %s", cacheKey)
	}

	return err
}

// RemoveUsersFromProject remove user
func (s *Store) RemoveUsersFromProject(accessToken, projectID string, userIDs ...string) error {
	if err := s.checkUserExist(userIDs...); err != nil {
		return err
	}

	err := s.dao.Project.RemoveUsersFromProject(projectID, userIDs...)

	cacheKey := "token_" + accessToken
	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debug("delete token from cache failed, key: %s", cacheKey)
		}
		s.log.Debug("delete token from cache success, key: %s", cacheKey)
	}

	return err
}

func (s *Store) checkUserExist(userIDs ...string) error {
	errs := make([]string, 0)
	noexist := make([]string, 0)

	for _, uid := range userIDs {

		ok, err := s.dao.User.CheckUserIsExistByID(uid)
		if err != nil {
			errs = append(errs, err.Error())
		}
		if !ok {
			noexist = append(noexist, uid)
		}
	}

	if len(errs) != 0 {
		err := strings.Join(errs, ",")
		return errors.New(err)
	}

	if len(noexist) != 0 {
		neu := strings.Join(noexist, ",")
		return exception.NewBadRequest("user %s not exist", neu)
	}

	return nil
}
