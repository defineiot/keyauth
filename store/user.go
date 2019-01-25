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

	return nil
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

	// query user's roles
	roles, err := s.dao.Role.ListUserRole(domainID, userID)
	if err != nil {
		return nil, err
	}
	for i := range roles {
		u.RoleNames = append(u.RoleNames, roles[i].Name)
	}

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
