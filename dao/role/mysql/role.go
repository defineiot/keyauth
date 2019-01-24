package mysql

import (
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/internal/exception"
)

func (s *store) CreateRole(r *role.Role) error {
	if err := r.Validate(); err != nil {
		return err
	}

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

		roles = append(roles, r)
	}

	return roles, nil
}

func (s *store) GetRole(id string) (*role.Role, error) {
	r := new(role.Role)
	err := s.stmts[FindRoleByID].QueryRow(id).Scan(&r.ID, &r.Name, &r.Description, &r.CreateAt, &r.UpdateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("role %s not find", id)
		}

		return nil, exception.NewInternalServerError("query single role error, %s", err)
	}

	return r, nil
}

func (s *store) GetRoleByName(name string) (*role.Role, error) {
	r := new(role.Role)
	err := s.stmts[FindRoleByName].QueryRow(name).Scan(&r.ID, &r.Name, &r.Description, &r.CreateAt, &r.UpdateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("role %s not find", name)
		}

		return nil, exception.NewInternalServerError("query single role error, %s", err)
	}

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

func (s *store) ListUserRole(domainID, userID string) ([]*role.Role, error) {
	rows, err := s.stmts[FindRoleByUser].Query()
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

		roles = append(roles, r)
	}

	return roles, nil
}
