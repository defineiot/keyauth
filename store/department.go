package store

import (
	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateDepartment todo
func (s *Store) CreateDepartment(dep *models.Department) error {
	has, err := s.dao.Department.GetDepartmentByName(dep.DomainID, dep.Name)
	if err != nil {
		if _, ok := err.(*exception.NotFound); !ok {
			return err
		}
	}

	if has != nil && has.ID != "" {
		return exception.NewBadRequest("the department %s has exist", dep.Name)
	}

	// 检查项目是否合法
	for _, pid := range dep.ProjectIDs {
		p, err := s.dao.Project.GetProjectByID(pid)
		if err != nil {
			return exception.NewBadRequest(err.Error())
		}
		dep.Projects = append(dep.Projects, p)
	}

	// 检查角色是否存在
	for _, rid := range dep.RoleIDs {
		r, err := s.dao.Role.GetRole(rid)
		if err != nil {
			return exception.NewBadRequest(err.Error())
		}
		dep.Roles = append(dep.Roles, r)
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
func (s *Store) ListSubDepartments(domainID, parentID string) ([]*models.Department, error) {
	depts, err := s.dao.Department.ListSubDepartments(domainID, parentID)
	if err != nil {
		return nil, err
	}

	return depts, nil
}

// GetDepartment todo
func (s *Store) GetDepartment(depID string) (*models.Department, error) {
	var err error

	dep := new(models.Department)
	cacheKey := s.cachePrefix.dep + depID

	if s.isCache {
		if s.cache.Get(cacheKey, dep) {
			s.log.Debug("get department from cache key: %s", cacheKey)
			return dep, nil
		}
		s.log.Debug("get department from cache failed, key: %s", cacheKey)
	}

	dep, err = s.dao.Department.GetDepartment(depID)
	if err != nil {
		return nil, err
	}
	if dep == nil {
		return nil, exception.NewBadRequest("department %s not found", depID)
	}

	// 填充部门关联的用户相关数据
	users, err := s.dao.User.ListDepartmentUsers(depID)
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		u.Domain = nil
		u.DefaultProject = nil
		u.Department = nil

		roles, err := s.dao.Role.ListUserRole(dep.DomainID, u.ID)
		if err != nil {
			return nil, err
		}
		u.Roles = roles
	}
	dep.Users = users

	// 填充部门相关的项目数据
	projects, err := s.dao.Project.ListDepartmentProjects(depID)
	if err != nil {
		return nil, err
	}
	dep.Projects = projects

	// 添加部门相关的角色数据
	roles, err := s.dao.Role.ListDepartmentRoles(depID)
	if err != nil {
		return nil, err
	}
	dep.Roles = roles

	if s.isCache {
		if !s.cache.Set(cacheKey, dep, s.ttl) {
			s.log.Debug("set app cache failed, key: %s", cacheKey)
		}
		s.log.Debug("set app cache ok, key: %s", cacheKey)
	}

	return dep, nil
}
