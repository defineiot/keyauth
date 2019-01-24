package store

import (
	"github.com/defineiot/keyauth/dao/role"
)

// CreateRole todo
func (s *Store) CreateRole(r *role.Role) error {
	return s.dao.Role.CreateRole(r)
}

// ListRoles todo
func (s *Store) ListRoles() ([]*role.Role, error) {
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
func (s *Store) GetRole(name string) (*role.Role, error) {
	r, err := s.dao.Role.GetRole(name)
	if err != nil {
		return nil, err
	}

	features, err := s.dao.Service.ListRoleFeatures(r.Name)
	if err != nil {
		return nil, err
	}

	r.Features = features

	return r, nil
}

// DeleteRole todo
func (s *Store) DeleteRole(name string) error {
	return s.dao.Role.DeleteRole(name)
}

// CheckRoleExist todo
func (s *Store) CheckRoleExist(roleName string) (bool, error) {
	return s.dao.Role.CheckRoleExist(roleName)
}

// AddFeaturesToRole todo
func (s *Store) AddFeaturesToRole(id string, features ...string) error {
	// errMsg := []string{}
	// notExist := []string{}

	// // 检查feature是否存在
	// for _, fid := range features {
	// 	ok, err := s.dao.Service.CheckFeatureIsExist(fid)
	// 	if err != nil {
	// 		errMsg = append(errMsg, err.Error())
	// 	}
	// 	if !ok {
	// 		notExist = append(notExist, fid)
	// 	}
	// }

	// if len(errMsg) != 0 {
	// 	return exception.NewInternalServerError(strings.Join(errMsg, ","))
	// }
	// if len(notExist) != 0 {
	// 	return exception.NewBadRequest("feature %v not exist", notExist)
	// }

	// // 对比, 哪些已经添加, 哪些未添加
	// exist, err := s.dao.Service.ListRoleFeatures(id)
	// if err != nil {
	// 	return err
	// }

	// needAdded := []*service.Feature{}
	// isAdded := []*service.Feature{}-
	// for _, infid := range features {
	// 	var inExist bool
	// 	for _, efif := range exist {
	// 		if efif.ID == infid {
	// 			inExist = true
	// 			isAdded = append(isAdded, efif)
	// 		}
	// 	}
	// 	if !inExist {
	// 		needAdded = append(needAdded, infid)
	// 	}
	// }

	// if len(isAdded) != 0 {
	// 	return exception.NewBadRequest("the feature %v has added to role: %s", isAdded, id)
	// }

	// return s.dao.Service.AssociateFeaturesToRole(id, needAdded...)
	return nil
}

// RemoveFeaturesFromRole todo
func (s *Store) RemoveFeaturesFromRole(name string, features ...int64) (bool, error) {
	return true, nil
}
