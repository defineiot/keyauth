package store

import (
	"fmt"

	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateRole todo
func (s *Store) CreateRole(r *models.Role) error {
	return s.dao.Role.CreateRole(r)
}

// ListRoles todo
func (s *Store) ListRoles() ([]*models.Role, error) {
	roles, err := s.dao.Role.ListRole()
	if err != nil {
		return nil, err
	}

	for _, r := range roles {
		features, err := s.dao.Service.ListRoleFeatures(r.ID)
		if err != nil {
			return nil, err
		}
		r.Features = features
	}

	return roles, nil
}

// GetRole todo
func (s *Store) GetRole(id string) (*models.Role, error) {
	r, err := s.dao.Role.GetRole(id)
	if err != nil {
		return nil, err
	}

	features, err := s.dao.Service.ListRoleFeatures(r.ID)
	if err != nil {
		return nil, err
	}

	r.Features = features

	return r, nil
}

// DeleteRole todo
func (s *Store) DeleteRole(name string) error {
	if name == systemAdminRoleName || name == domainAdminRoleName ||
		name == memberUserRoleName {
		return exception.NewForbidden("system initial role can't be delete")
	}

	return s.dao.Role.DeleteRole(name)
}

// CheckRoleExist todo
func (s *Store) CheckRoleExist(roleName string) (bool, error) {
	return s.dao.Role.CheckRoleExist(roleName)
}

// AddFeaturesToRole todo
func (s *Store) AddFeaturesToRole(id string, features ...string) error {
	// 获取当前角色
	ro, err := s.GetRole(id)
	if err != nil {
		return err
	}
	emap := make(map[string]bool)
	for i := range ro.Features {
		emap[ro.Features[i].ID] = true
	}

	// 获取系统所有功能列表
	fmap := make(map[string]*models.Feature)
	all, err := s.dao.Service.ListAllFeatures()
	if err != nil {
		return err
	}
	for i := range all {
		fmap[all[i].ID] = all[i]
	}

	// 判断要添加的功能是否存在, 是否已经添加
	notExist := []string{}
	hasAdded := []string{}
	needAdded := []*models.Feature{}
	for i := range features {
		if v, ok := fmap[features[i]]; !ok {
			notExist = append(notExist, features[i])
		} else {
			needAdded = append(needAdded, v)
		}

		if _, ok := emap[features[i]]; ok {
			hasAdded = append(hasAdded, features[i])
		}
	}

	if len(notExist) > 0 {
		return exception.NewBadRequest("the features: %s not exist", notExist)
	}
	if len(hasAdded) > 0 {
		return exception.NewBadRequest("the features: %s has added", hasAdded)
	}
	if len(needAdded) == 0 {
		return exception.NewBadRequest("no features need add")
	}

	return s.dao.Service.AssociateFeaturesToRole(id, needAdded...)
}

// RemoveFeaturesFromRole todo
func (s *Store) RemoveFeaturesFromRole(id string, features ...string) error {
	// 获取当前角色
	ro, err := s.GetRole(id)
	if err != nil {
		return err
	}
	emap := make(map[string]bool)
	for i := range ro.Features {
		emap[ro.Features[i].ID] = true
	}

	fmt.Println(emap)

	// 获取系统所有功能列表
	fmap := make(map[string]*models.Feature)
	all, err := s.dao.Service.ListAllFeatures()
	if err != nil {
		return err
	}
	for i := range all {
		fmap[all[i].ID] = all[i]
	}

	fmt.Println(fmap)

	// 判断要移除的功能是否存在
	notExist := []string{}
	notAdded := []string{}
	needRemove := []*models.Feature{}
	for i := range features {
		if v, ok := fmap[features[i]]; !ok {
			notExist = append(notExist, features[i])
		} else {
			needRemove = append(needRemove, v)
		}

		if _, ok := emap[features[i]]; !ok {
			notAdded = append(notAdded, features[i])
		}
	}

	if len(notExist) > 0 {
		return exception.NewBadRequest("the features: %s not exist", notExist)
	}
	if len(notAdded) > 0 {
		return exception.NewBadRequest("the features: %s not added", notAdded)
	}
	if len(needRemove) == 0 {
		return exception.NewBadRequest("no features need remove")
	}

	return s.dao.Service.UnlinkFeatureFromRole(id, needRemove...)
}
