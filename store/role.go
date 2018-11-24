package store

import (
	"strings"

	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateRole todo
func (s *Store) CreateRole(name, description string) (*role.Role, error) {
	return s.role.CreateRole(name, description)
}

// ListRoles todo
func (s *Store) ListRoles() ([]*role.Role, error) {
	roles, err := s.role.ListRole()
	if err != nil {
		return nil, err
	}

	for _, r := range roles {
		features, err := s.service.ListRoleFeatures(r.Name)
		if err != nil {
			return nil, err
		}
		r.Featrues = features
	}

	return roles, nil
}

// GetRole todo
func (s *Store) GetRole(name string) (*role.Role, error) {
	r, err := s.role.GetRole(name)
	if err != nil {
		return nil, err
	}

	features, err := s.service.ListRoleFeatures(r.Name)
	if err != nil {
		return nil, err
	}

	r.Featrues = features

	return r, nil
}

// DeleteRole todo
func (s *Store) DeleteRole(name string) error {
	return s.role.DeleteRole(name)
}

// CheckRoleExist todo
func (s *Store) CheckRoleExist(roleName string) (bool, error) {
	return s.role.CheckRoleExist(roleName)
}

// AddFeaturesToRole todo
func (s *Store) AddFeaturesToRole(name string, features ...int64) error {
	errMsg := []string{}
	notExist := []int64{}
	for _, fid := range features {
		ok, err := s.service.CheckFeatureIsExist(fid)
		if err != nil {
			errMsg = append(errMsg, err.Error())
		}
		if !ok {
			notExist = append(notExist, fid)
		}
	}

	if len(errMsg) != 0 {
		return exception.NewInternalServerError(strings.Join(errMsg, ","))
	}
	if len(notExist) != 0 {
		return exception.NewBadRequest("feature %v not exist", notExist)
	}

	exist, err := s.role.GetRoleFeature(name)
	if err != nil {
		return err
	}

	needAdded := []int64{}
	isAdded := []int64{}
	for _, infid := range features {
		var inExist bool
		for _, efif := range exist {
			if efif == infid {
				inExist = true
				isAdded = append(isAdded, infid)
			}
		}
		if !inExist {
			needAdded = append(needAdded, infid)
		}
	}

	if len(isAdded) != 0 {
		return exception.NewBadRequest("the feature %v has added to role: %s", isAdded, name)
	}

	return s.role.AssociateFeaturesToRole(name, needAdded...)
}

// RemoveFeaturesFromRole todo
func (s *Store) RemoveFeaturesFromRole(name string, features ...int64) (bool, error) {
	return true, nil
}
