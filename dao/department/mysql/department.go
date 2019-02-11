package mysql

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/internal/exception"
)

const (
	// MaxDepartmentDeep 部门层级深度限制
	MaxDepartmentDeep = 10
)

func (s *store) CreateDepartment(d *models.Department) error {
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

	tx, err := s.db.Begin()
	if err != nil {
		return exception.NewInternalServerError("start create department transaction error, %s", err)
	}
	defer tx.Commit()

	// 创建部门
	sdPre, err := tx.Prepare(s.unprepared[SaveDepartment])
	if err != nil {
		return exception.NewInternalServerError("prepare save department sql error, %s", err)
	}
	defer sdPre.Close()

	if _, err := sdPre.Exec(d.ID, d.Number, d.Name, d.ParentID, d.Grade,
		d.Path, d.ManagerID, d.DomainID, d.CreateAt); err != nil {
		tx.Rollback()
		return exception.NewInternalServerError("insert save department exec sql err, %s", err)
	}

	// 部门项目
	if len(d.ProjectIDs) > 0 {
		sdpPre, err := tx.Prepare(s.unprepared[SaveDepartmentProject])
		if err != nil {
			return exception.NewInternalServerError("prepare save department project sql error, %s", err)
		}
		defer sdpPre.Close()

		for _, pid := range d.ProjectIDs {
			if _, err := sdpPre.Exec(d.ID, pid); err != nil {
				tx.Rollback()
				return exception.NewInternalServerError("save department project error, %s", err)
			}
		}
	}

	// 部门角色
	if len(d.RoleIDs) > 0 {
		sdrPre, err := tx.Prepare(s.unprepared[SaveDepartmentRole])
		if err != nil {
			return exception.NewInternalServerError("prepare save department role sql error, %s", err)
		}
		defer sdrPre.Close()

		for _, rid := range d.RoleIDs {
			if _, err := sdrPre.Exec(d.ID, rid); err != nil {
				tx.Rollback()
				return exception.NewInternalServerError("save department role error, %s", err)
			}
		}
	}

	return nil
}

func (s *store) GetDepartment(depID string) (*models.Department, error) {
	d := new(models.Department)

	err := s.stmts[FindDepartment].QueryRow(depID).Scan(&d.ID, &d.Number, &d.Name,
		&d.ParentID, &d.Grade, &d.Path, &d.ManagerID, &d.DomainID, &d.CreateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("department %s not find", depID)
		}

		return nil, exception.NewInternalServerError("query single verify code error, %s", err)
	}

	return d, nil
}

func (s *store) ListSubDepartments(domainID, parentDepID string) ([]*models.Department, error) {
	rows, err := s.stmts[FindSubDepartments].Query(parentDepID, domainID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user's invitation records error, %s", err)
	}
	defer rows.Close()

	deps := []*models.Department{}
	for rows.Next() {
		d := new(models.Department)
		if err := rows.Scan(&d.ID, &d.Number, &d.Name, &d.ParentID, &d.Grade, &d.Path,
			&d.ManagerID, &d.DomainID, &d.CreateAt); err != nil {
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

func (s *store) GetDepartmentByName(domainID, departmentName string) (*models.Department, error) {
	d := new(models.Department)

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
