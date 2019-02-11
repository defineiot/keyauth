package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/department"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/logger"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	SaveDepartment        = "save-department"
	SaveDepartmentProject = "save-department-project"
	SaveDepartmentRole    = "save-department-role"
	FindDepartment        = "find-department"
	FindDepartmentByName  = "find-department-by-name"
	FindSubDepartments    = "find-sub-departments"
	CountSubDepartments   = "count-sub-departments"
	DeleteDepartment      = "delete-departments"
)

// NewDepartmentStore use to create domain storage service
func NewDepartmentStore(opt *dao.Options) (department.Store, error) {
	if opt.DB == nil {
		return nil, exception.NewInternalServerError("the db connection required")
	}
	if opt.LOG == nil {
		return nil, exception.NewInternalServerError("the logger not config")
	}

	unprepared := map[string]string{
		SaveDepartment: `
			INSERT INTO departments (id, number, name, parent, grade, path, manager, domain_id, create_at) 
			VALUES (?,?,?,?,?,?,?,?,?);
		`,
		SaveDepartmentProject: `
			INSERT INTO department_project_mappings (department_id, project_id) 
			VALUES (?,?);
		`,
		SaveDepartmentRole: `
			INSERT INTO department_role_mappings (department_id, role_id) 
			VALUES (?,?);
		`,
		FindDepartment: `
			SELECT id, number, name, parent, grade, path, manager, domain_id, create_at 
			FROM departments 
			WHERE id = ?;
		`,
		FindDepartmentByName: `
			SELECT id, number, name, parent, grade, path, manager, domain_id, create_at 
			FROM departments 
			WHERE name = ? 
			AND domain_id = ?;
		`,
		FindSubDepartments: `
			SELECT id, number, name, parent, grade, path, manager, domain_id, create_at 
			FROM departments 
			WHERE parent = ? 
			AND domain_id = ? 
			ORDER BY create_at 
			DESC;
		`,
		CountSubDepartments: `
			SELECT count(*) 
			FROM departments 
			WHERE parent = ?;
		`,
		DeleteDepartment: `
			DELETE FROM departments 
			WHERE id = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(opt.DB, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare department store query statment error, %s", err)
	}

	s := store{
		db:         opt.DB,
		stmts:      stmts,
		unprepared: unprepared,
		log:        opt.LOG,
	}

	return &s, nil
}

// DomainManager is use mongodb as storage
type store struct {
	db         *sql.DB
	stmts      map[string]*sql.Stmt
	unprepared map[string]string
	key        string
	log        logger.Logger
}

// Close closes the database, releasing any open resources.
func (s *store) Close() error {
	return s.db.Close()
}

func init() {
	dao.Registe(NewDepartmentStore)
}
