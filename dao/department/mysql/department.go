package mysql

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/defineiot/keyauth/dao/department"
	"github.com/defineiot/keyauth/internal/exception"
)

const (
	// MaxDepartmentDeep 部门层级深度限制
	MaxDepartmentDeep = 10
)

func (s *store) CreateDepartment(d *department.Department) error {
	if err := d.Validate(); err != nil {
		return err
	}

	if len(strings.Split(d.ParentID, "/")) > MaxDepartmentDeep {
		return exception.NewBadRequest("max department deep is %d, but overflow", MaxDepartmentDeep)
	}

	// 默认为顶层部门(根部门)
	d.ID = uuid.NewV4().String()
	d.CreateAt = time.Now().Unix()
	d.Number = d.DomainID[:2] + strconv.FormatInt(d.CreateAt, 36) + d.ID[:4]

	if d.ParentID == "" {
		d.ParentID = "/"
		d.Path = d.Number
	} else {
		parentDep, err := s.GetDepartment(d.ParentID)
		if err != nil {
			return err
		}
		d.Path = parentDep.Path + "/" + d.Number
	}

	_, err := s.stmts[SaveDepartment].Exec(d.ID, d.Number, d.Name, d.ParentID, d.Grade, d.Path, d.ManagerID, d.DomainID, d.CreateAt)
	if err != nil {
		return exception.NewInternalServerError("insert save department exec sql err, %s", err)
	}

	return nil
}

func (s *store) GetDepartment(depID string) (*department.Department, error) {
	d := new(department.Department)

	err := s.stmts[FindDepartment].QueryRow(depID).Scan(
		&d.ID, &d.Number, &d.Name, &d.ParentID, &d.Grade, &d.Path, &d.ManagerID, &d.DomainID, &d.CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("department %s not find", depID)
		}

		return nil, exception.NewInternalServerError("query single verify code error, %s", err)
	}

	return d, nil
}

func (s *store) ListSubDepartments(parentDepID string) ([]*department.Department, error) {
	rows, err := s.stmts[FindSubDepartments].Query(parentDepID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user's invitation records error, %s", err)
	}
	defer rows.Close()

	deps := []*department.Department{}
	for rows.Next() {
		d := new(department.Department)
		if err := rows.Scan(&d.ID, &d.Number, &d.Name, &d.ParentID, &d.Grade, &d.Path, &d.ManagerID, &d.DomainID, &d.CreateAt); err != nil {
			return nil, exception.NewInternalServerError("scan user's project id error, %s", err)
		}
		deps = append(deps, d)
	}
	return deps, nil
}

func (s *store) DelDepartment(depID string) error {
	var count int
	if err := s.stmts[CountSubDepartments].QueryRow(depID).Scan(&count); err != nil {
		return exception.NewInternalServerError("delete depepartment error when count sub depeartments, %s", err)
	}

	if count != 0 {
		return exception.NewBadRequest("the department has %d sub departments, your should delete them first!", count)
	}

	ret, err := s.stmts[DeleteDepartment].Exec(depID)
	if err != nil {
		return exception.NewInternalServerError("delete department exec sql error, %s", err)
	}
	affect, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if affect == 0 {
		return exception.NewBadRequest("department %s not exist", depID)
	}

	return nil
}

func (s *store) GetDepartmentByName(domainID, departmentName string) (*department.Department, error) {
	d := new(department.Department)

	err := s.stmts[FindDepartmentByName].QueryRow(departmentName, domainID).Scan(
		&d.ID, &d.Number, &d.Name, &d.ParentID, &d.Grade, &d.Path, &d.ManagerID, &d.DomainID, &d.CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("department %s not find", departmentName)
		}

		return nil, exception.NewInternalServerError("query single verify code error, %s", err)
	}

	return d, nil
}
