package store

import (
	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateMemberUser use to create user
func (s *Store) CreateMemberUser(u *user.User) error {
	// 判断用户名是否存在
	if _, err := s.dao.User.GetUser(user.Account, u.Account); err != nil {
		if _, ok := err.(*exception.NotFound); !ok {
			return err
		}
	} else {
		return exception.NewBadRequest("account: %s is exist", u.Account)
	}

	// 判断用户的秘密是否符合复杂度要求

	// 如果用户未选择部门, 则使用默认部门
	if u.Department.ID == "" {
		defaultDep, err := s.dao.Department.GetDepartmentByName(u.Domain.ID, defaultDepartmentName)
		if err != nil {
			return err
		}

		u.Department.ID = defaultDep.ID
	}

	// 创建用户
	if err := s.dao.User.CreateUser(u); err != nil {
		return err
	}

	// 绑定成员角色
	role, err := s.dao.Role.GetRoleByName(memberUserRoleName)
	if err != nil {
		return err
	}

	if err := s.dao.User.BindRole(u.Domain.ID, u.ID, role.ID); err != nil {
		return err
	}

	// 查询出域的具体详情
	dom, err := s.dao.Domain.GetDomainByID(u.Domain.ID)
	if err != nil {
		return err
	}

	// 查询用户部门的详情
	dep, err := s.dao.Department.GetDepartment(u.Department.ID)
	if err != nil {
		return err
	}

	roles, err := s.dao.Role.ListUserRole(u.Domain.ID, u.ID)
	if err != nil {
		return err
	}

	u.Domain = dom
	u.Department = dep
	u.Roles = roles

	return nil
}

// ListMemberUsers list all user
func (s *Store) ListMemberUsers(domainID string) ([]*user.User, error) {
	users, err := s.dao.User.ListDomainUsers(domainID)
	if err != nil {
		return nil, err
	}

	for i := range users {
		u := users[i]
		// 查询出域的具体详情
		dom, err := s.dao.Domain.GetDomainByID(u.Domain.ID)
		if err != nil {
			return nil, err
		}

		// 查询用户部门的详情
		dep, err := s.dao.Department.GetDepartment(u.Department.ID)
		if err != nil {
			return nil, err
		}

		roles, err := s.dao.Role.ListUserRole(u.Domain.ID, u.ID)
		if err != nil {
			return nil, err
		}

		u.Domain = dom
		u.Department = dep
		u.Roles = roles
	}

	return users, nil
}

// GetUser get an user
func (s *Store) GetUser(domainID, userID string) (*user.User, error) {
	var err error

	u := new(user.User)
	cacheKey := s.cachePrefix.user + userID

	if s.isCache {
		if s.cache.Get(cacheKey, u) {
			s.log.Debug("get project from cache key: %s", cacheKey)
			return u, nil
		}
		s.log.Debug("get project from cache failed, key: %s", cacheKey)
	}

	u, err = s.dao.User.GetUser(user.UserID, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, exception.NewBadRequest("user %s not found", userID)
	}

	// 查询出域的具体详情
	dom, err := s.dao.Domain.GetDomainByID(u.Domain.ID)
	if err != nil {
		return nil, err
	}

	// 查询用户部门的详情
	dep, err := s.dao.Department.GetDepartment(u.Department.ID)
	if err != nil {
		return nil, err
	}

	roles, err := s.dao.Role.ListUserRole(u.Domain.ID, u.ID)
	if err != nil {
		return nil, err
	}

	u.Domain = dom
	u.Department = dep
	u.Roles = roles

	if u.DefaultProject.ID != "" {
		pro, err := s.dao.Project.GetProjectByID(u.DefaultProject.ID)
		if err != nil {
			return nil, err
		}
		u.DefaultProject = pro
	}

	if s.isCache {
		if !s.cache.Set(cacheKey, u, s.ttl) {
			s.log.Debug("set user cache failed, key: %s", cacheKey)
		}
		s.log.Debug("set user cache ok, key: %s", cacheKey)
	}

	return u, nil
}

// DeleteUser delete an user by id
func (s *Store) DeleteUser(domainID, userID string) error {
	var err error

	cacheKey := s.cachePrefix.user + userID

	err = s.dao.User.DeleteUser(domainID, userID)
	if err != nil {
		return err
	}

	if s.isCache {
		if !s.cache.Delete(cacheKey) {
			s.log.Debug("delete user from cache failed, key: %s", cacheKey)
		}
		s.log.Debug("delete user from cache success, key: %s", cacheKey)
	}

	return nil
}
