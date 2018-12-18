package mysql

import (
	"database/sql"
	"time"

	"github.com/satori/go.uuid"

	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/internal/exception"
)

func (s *store) CreateRole(r *role.Role) error {
	ok, err := s.CheckRoleExist(r.Name)
	if err != nil {
		return err
	}
	if ok {
		return exception.NewBadRequest("role %s exist", r.Name)
	}

	r.CreateAt = time.Now().Unix()
	r.ID = uuid.NewV4().String()

	_, err = s.stmts[SaveRole].Exec(r.ID, r.Name, r.Description, r.CreateAt)
	if err != nil {
		return exception.NewInternalServerError("insert role exec sql err, %s", err)
	}
	return nil
}

func (s *store) CheckRoleExist(name string) (bool, error) {
	var n string
	if err := s.stmts[CheckRole].QueryRow(name).Scan(&n); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("query role exist error, %s", err)
	}

	return true, nil
}

func (s *store) ListRole() ([]*role.Role, error) {
	rows, err := s.stmts[FindAllRole].Query()
	if err != nil {
		return nil, exception.NewInternalServerError("query role list error, %s", err)
	}
	defer rows.Close()

	roles := []*role.Role{}
	for rows.Next() {
		r := new(role.Role)
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.CreateAt, &r.UpdateAt); err != nil {
			return nil, exception.NewInternalServerError("scan project record error, %s", err)
		}

		//  查询该role的功能列表
		features, err := s.getRoleFeatures(r.ID)
		if err != nil {
			return nil, err
		}
		r.Featrues = features

		roles = append(roles, r)
	}

	return roles, nil
}

func (s *store) GetRole(id string) (*role.Role, error) {
	r := new(role.Role)
	err := s.stmts[FindOneRole].QueryRow(id).Scan(&r.ID, &r.Name, &r.Description, &r.CreateAt, &r.UpdateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("role %s not find", id)
		}

		return nil, exception.NewInternalServerError("query single role error, %s", err)
	}

	//  查询该role的功能列表
	features, err := s.getRoleFeatures(r.ID)
	if err != nil {
		return nil, err
	}
	r.Featrues = features

	return r, nil
}

func (s *store) UpdateRole(name, description string) (*role.Role, error) {
	return nil, nil
}

func (s *store) DeleteRole(id string) error {
	ret, err := s.stmts[DeleteRole].Exec(id)
	if err != nil {
		return exception.NewInternalServerError("delete role exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete role affected rows error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("role %s not exist", id)
	}

	return nil
}

func (s *store) getRoleFeatures(id string) ([]*service.Feature, error) {
	rows, err := s.stmts[GetRoleFeatures].Query(id)
	if err != nil {
		return nil, exception.NewInternalServerError("query role features error, %s", err)
	}
	defer rows.Close()

	features := []*service.Feature{}
	for rows.Next() {
		f := new(service.Feature)
		if err := rows.Scan(&f.ID, &f.Name, &f.Method, &f.Endpoint, &f.Description, &f.IsDeleted, &f.WhenDeletedVersion,
			&f.IsAdded, &f.WhenAddedVersion, &f.ServiceID); err != nil {
			return nil, exception.NewInternalServerError("scan role feature mapping record error, %s", err)
		}
		features = append(features, f)
	}

	return features, nil
}

// func (s *store) AssociateFeaturesToRole(name string, features ...int64) error {
// 	// start transaction
// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		return fmt.Errorf("start associate features to role transaction error, %s", err)
// 	}

// 	// prepare insert feature
// 	mappingPre, err := tx.Prepare("INSERT INTO roles_features_mapping (feature_id, role_name) VALUES (?,?);")
// 	if err != nil {
// 		return exception.NewInternalServerError("prepare insert feature role mapping stmt error, name: %s, %s", name, err)
// 	}
// 	defer mappingPre.Close()

// 	for _, f := range features {
// 		// exec sql
// 		_, err := mappingPre.Exec(f, name)
// 		if err != nil {
// 			if err := tx.Rollback(); err != nil {
// 				s.log.Errorf("insert feature role mapping transaction rollback error, %s", err)
// 			}
// 			return exception.NewInternalServerError("insert feature role mapping exec sql err, %s", err)
// 		}

// 	}

// 	// commit transaction
// 	if err := tx.Commit(); err != nil {
// 		s.log.Errorf("insert feature transaction rollback error, %s", err)
// 		return exception.NewInternalServerError("insert feature transaction commit error, but rollback success, %s", err)
// 	}

// 	return nil
// }

// func (s *store) UnlinkFeatureFromRole(name string, features ...int64) (bool, error) {
// 	// start transaction
// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		return false, fmt.Errorf("start unlink features from role transaction error, %s", err)
// 	}

// 	// prepare insert feature
// 	mappingPre, err := tx.Prepare("DELETE FROM roles_features_mapping WHERE feature_id = ? AND role_name = ?;")
// 	if err != nil {
// 		return false, exception.NewInternalServerError("prepare insert feature role mapping stmt error, name: %s, %s", name, err)
// 	}
// 	defer mappingPre.Close()

// 	for _, f := range features {
// 		// exec sql
// 		_, err := mappingPre.Exec(f, name)
// 		if err != nil {
// 			if err := tx.Rollback(); err != nil {
// 				s.log.Errorf("unlik feature role mapping transaction rollback error, %s", err)
// 			}
// 			return false, exception.NewInternalServerError("unlik feature role mapping exec sql err, %s", err)
// 		}

// 	}

// 	// commit transaction
// 	if err := tx.Commit(); err != nil {
// 		s.log.Errorf("unlik feature transaction rollback error, %s", err)
// 		return false, exception.NewInternalServerError("unlik feature transaction commit error, but rollback success, %s", err)
// 	}

// 	return true, nil
// }
