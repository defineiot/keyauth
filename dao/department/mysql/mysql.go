package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao/department"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/logger"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	SaveDepartment      = "save-department"
	FindDepartment      = "find-department"
	FindSubDepartments  = "find-sub-departments"
	CountSubDepartments = "count-sub-departments"
	DeleteDepartment    = "delete-departments"
)

// NewDepartmentStore use to create domain storage service
func NewDepartmentStore(db *sql.DB, log logger.Logger) (department.Store, error) {
	unprepared := map[string]string{
		SaveDepartment: `
			INSERT INTO departments (id, name, parent, grade, path, manager, domain_id, create_at) 
			VALUES (?,?,?,?,?,?,?,?);
		`,
		FindDepartment: `
			SELECT id, name, parent, grade, path, manager, domain_id, create_at 
			FROM departments 
			WHERE id = ?;
		`,
		FindSubDepartments: `
			SELECT id, name, parent, grade, path, manager, domain_id, create_at 
			FROM departments 
			WHERE parent = ?;
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
	stmts, err := tools.PrepareStmts(db, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare department store query statment error, %s", err)
	}

	s := store{
		db:         db,
		stmts:      stmts,
		unprepared: unprepared,
		log:        log,
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
