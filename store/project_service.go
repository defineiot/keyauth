package store

import (
	"errors"
	"strings"

	"github.com/defineiot/keyauth/dao/project"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateProject use to create an project
func (s *Store) CreateProject(domainID, name, description string, enabled bool) (*project.Project, error) {
	if domainID == "" {
		return nil, exception.NewBadRequest("domainID or domainName required one")
	}

	// check domain exist
	ok, err := s.domain.CheckDomainIsExistByID(domainID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("domain %s not exist", domainID)
	}

	return s.project.CreateProject(domainID, name, description, enabled)
}

// ListDomainProjects list domain projects
func (s *Store) ListDomainProjects(domainID string) ([]*project.Project, error) {
	// check domain exist
	ok, err := s.domain.CheckDomainIsExistByID(domainID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("domain %s not exist", domainID)
	}

	return s.project.ListDomainProjects(domainID)
}

// ListUserProjects todo
func (s *Store) ListUserProjects(domainID, userID string) ([]*project.Project, error) {
	pids, err := s.user.ListUserProjects(domainID, userID)
	if err != nil {
		return nil, err
	}

	projects := []*project.Project{}
	for _, pid := range pids {
		p, err := s.project.GetProject(pid)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

// GetProject get one project
func (s *Store) GetProject(id string) (*project.Project, error) {
	var err error

	pro := new(project.Project)
	cacheKey := "project_" + id

	if s.isCache {
		if s.cache.Get(cacheKey, pro) {
			s.log.Debugf("get project from cache key: %s", cacheKey)
			return pro, nil
		}
		s.log.Debugf("get project from cache failed, key: %s", cacheKey)
	}

	pro, err = s.project.GetProject(id)
	if err != nil {
		return nil, err
	}
	if pro == nil {
		return nil, exception.NewBadRequest("project %s not found", id)
	}

	if s.isCache {
		if !s.cache.Set(cacheKey, pro, s.ttl) {
			s.log.Debugf("set project cache failed, key: %s", cacheKey)
		}
		s.log.Debugf("set project cache ok, key: %s", cacheKey)
	}

	return pro, nil
}

// DeleteProjectByID project
func (s *Store) DeleteProjectByID(id string) error {
	var err error

	cacheKey := "project_" + id

	err = s.project.DeleteProjectByID(id)
	if err != nil {
		return err
	}

	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debugf("delete project from cache failed, key: %s", cacheKey)
		}
		s.log.Debugf("delete project from cache success, key: %s", cacheKey)
	}

	return nil
}

// DeleteProjectByName project
func (s *Store) DeleteProjectByName(projectName, domainID string) error {
	var err error

	cacheKey := "project_" + domainID + projectName

	err = s.project.DeleteProjectByName(projectName, domainID)
	if err != nil {
		return err
	}

	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debugf("delete project from cache failed, key: %s", cacheKey)
		}
		s.log.Debugf("delete project from cache success, key: %s", cacheKey)
	}

	return nil
}

// AddUsersToProject add user
func (s *Store) AddUsersToProject(accessToken, projectID string, userIDs ...string) error {
	if err := s.checkUserExist(userIDs...); err != nil {
		return err
	}

	err := s.project.AddUsersToProject(projectID, userIDs...)

	cacheKey := "token_" + accessToken
	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debugf("delete token from cache failed, key: %s", cacheKey)
		}
		s.log.Debugf("delete token from cache success, key: %s", cacheKey)
	}

	return err
}

// RemoveUsersFromProject remove user
func (s *Store) RemoveUsersFromProject(accessToken, projectID string, userIDs ...string) error {
	if err := s.checkUserExist(userIDs...); err != nil {
		return err
	}

	err := s.project.RemoveUsersFromProject(projectID, userIDs...)

	cacheKey := "token_" + accessToken
	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debugf("delete token from cache failed, key: %s", cacheKey)
		}
		s.log.Debugf("delete token from cache success, key: %s", cacheKey)
	}

	return err
}

func (s *Store) checkUserExist(userIDs ...string) error {
	errs := make([]string, 0)
	noexist := make([]string, 0)

	for _, uid := range userIDs {
		ok, err := s.user.CheckUserIsExistByID(uid)
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
