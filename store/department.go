package store

import (
	"github.com/defineiot/keyauth/dao/application"
	"github.com/defineiot/keyauth/dao/department"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateDepartment todo
func (s *Store) CreateDepartment(dep *department.Department) error {
	has, err := s.dao.Department.GetDepartmentByName(dep.DomainID, dep.Name)
	if _, ok := err.(*exception.NotFound); !ok {
		return err
	}
	if has != nil && has.ID != "" {
		return exception.NewBadRequest("the department %s has exist", dep.Name)
	}

	if err := s.dao.Department.CreateDepartment(dep); err != nil {
		return err
	}

	return nil
}

// DeleteDepartment todo
func (s *Store) DeleteDepartment(id string) error {
	var err error

	cacheKey := s.cachePrefix.dep + id

	dep, err := s.dao.Department.GetDepartment(id)
	if err != nil {
		return err
	}

	err = s.dao.Department.DelDepartment(dep.ID)
	if err != nil {
		return err
	}
	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debug("delete department from cache failed, key: %s", cacheKey)
		}
		s.log.Debug("delete department from cache success, key: %s", cacheKey)
	}

	return nil
}

// ListSubDepartments todo
func (s *Store) ListSubDepartments(domainID, parentID string) ([]*department.Department, error) {
	depts, err := s.dao.Department.ListSubDepartments(domainID, parentID)
	if err != nil {
		return nil, err
	}

	return depts, nil
}

// GetDepartment todo
func (s *Store) GetDepartment(appid string) (*application.Application, error) {
	var err error

	app := new(application.Application)
	cacheKey := s.cachePrefix.dep + appid

	if s.isCache {
		if s.cache.Get(cacheKey, app) {
			s.log.Debug("get app from cache key: %s", cacheKey)
			return app, nil
		}
		s.log.Debug("get app from cache failed, key: %s", cacheKey)
	}

	app, err = s.dao.Application.GetApplication(appid)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, exception.NewBadRequest("app %s not found", appid)
	}

	if s.isCache {
		if !s.cache.Set(cacheKey, app, s.ttl) {
			s.log.Debug("set app cache failed, key: %s", cacheKey)
		}
		s.log.Debug("set app cache ok, key: %s", cacheKey)
	}

	return app, nil
}
