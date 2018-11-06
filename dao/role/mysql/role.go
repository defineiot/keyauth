package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/internal/exception"
)

func (s *store) CreateRole(name, description string) (*role.Role, error) {
	ok, err := s.CheckRoleExist(name)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, exception.NewBadRequest("role %s exist", name)
	}

	r := role.Role{Name: name, Description: description, CreateAt: time.Now().Unix()}
	_, err = s.stmts[SaveRole].Exec(r.Name, r.Description, r.CreateAt)
	if err != nil {
		return nil, exception.NewInternalServerError("insert role exec sql err, %s", err)
	}
	return &r, nil

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
		r := role.Role{}
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.CreateAt); err != nil {
			return nil, exception.NewInternalServerError("scan project record error, %s", err)
		}
		roles = append(roles, &r)
	}

	return roles, nil
}

func (s *store) GetRole(name string) (*role.Role, error) {
	r := role.Role{}
	err := s.stmts[FindOneRole].QueryRow(name).Scan(&r.ID, &r.Name, &r.Description, &r.CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("role %s not find", name)
		}

		return nil, exception.NewInternalServerError("query single role error, %s", err)
	}

	return &r, nil
}

func (s *store) UpdateRole(name, description string) (*role.Role, error) {
	return nil, nil
}

func (s *store) DeleteRole(name string) error {
	ret, err := s.stmts[DeleteRole].Exec(name)
	if err != nil {
		return exception.NewInternalServerError("delete role exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete role affected rows error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("role %s not exist", name)
	}

	return nil
}

func (s *store) GetRoleFeature(name string) ([]int64, error) {
	ok, err := s.CheckRoleExist(name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("role %s not exist", name)
	}

	rows, err := s.stmts[GetRoleFeatures].Query(name)
	if err != nil {
		return nil, exception.NewInternalServerError("query role features error, %s", err)
	}
	defer rows.Close()

	fids := []int64{}
	for rows.Next() {
		var fid int64
		if err := rows.Scan(&fid); err != nil {
			return nil, exception.NewInternalServerError("scan role feature mapping record error, %s", err)
		}
		fids = append(fids, fid)
	}

	return fids, nil
}

func (s *store) VerifyRole(name string, feature string) (bool, error) {
	return true, nil
}

func (s *store) AssociateFeaturesToRole(name string, features ...int64) error {
	// start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("start associate features to role transaction error, %s", err)
	}

	// prepare insert feature
	mappingPre, err := tx.Prepare("INSERT INTO roles_features_mapping (feature_id, role_name) VALUES (?,?);")
	if err != nil {
		return exception.NewInternalServerError("prepare insert feature role mapping stmt error, name: %s, %s", name, err)
	}
	defer mappingPre.Close()

	for _, f := range features {
		// exec sql
		_, err := mappingPre.Exec(f, name)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				s.log.Errorf("insert feature role mapping transaction rollback error, %s", err)
			}
			return exception.NewInternalServerError("insert feature role mapping exec sql err, %s", err)
		}

	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		s.log.Errorf("insert feature transaction rollback error, %s", err)
		return exception.NewInternalServerError("insert feature transaction commit error, but rollback success, %s", err)
	}

	return nil
}

func (s *store) UnlinkFeatureFromRole(name string, features ...int64) (bool, error) {
	// start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return false, fmt.Errorf("start unlink features from role transaction error, %s", err)
	}

	// prepare insert feature
	mappingPre, err := tx.Prepare("DELETE FROM roles_features_mapping WHERE feature_id = ? AND role_name = ?;")
	if err != nil {
		return false, exception.NewInternalServerError("prepare insert feature role mapping stmt error, name: %s, %s", name, err)
	}
	defer mappingPre.Close()

	for _, f := range features {
		// exec sql
		_, err := mappingPre.Exec(f, name)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				s.log.Errorf("unlik feature role mapping transaction rollback error, %s", err)
			}
			return false, exception.NewInternalServerError("unlik feature role mapping exec sql err, %s", err)
		}

	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		s.log.Errorf("unlik feature transaction rollback error, %s", err)
		return false, exception.NewInternalServerError("unlik feature transaction commit error, but rollback success, %s", err)
	}

	return true, nil
}
