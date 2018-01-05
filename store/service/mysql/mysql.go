package mysql

import (
	"database/sql"

	"openauth/api/exception"
	"openauth/store/service"
	"openauth/tools"
)

const (
	SaveService   = "save-service"
	UpdateService = "update-service"
	DeleteService = "delete-service"
	FindAll       = "find-all"
	FindOneByID   = "find-one"
)

// NewServiceStore use to create domain storage service
func NewServiceStore(db *sql.DB) (service.Store, error) {
	unprepared := map[string]string{
		SaveService: `
			INSERT INTO services (id, name, description, enabled, status, status_update_at, version, create_at) 
			VALUES (?,?,?,?,?,?,?,?);
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(db, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare application store query statment error, %s", err)
	}

	s := store{
		db:    db,
		stmts: stmts,
	}

	return &s, nil
}

// DomainManager is use mongodb as storage
type store struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

// Close closes the database, releasing any open resources.
func (s *store) Close() error {
	return s.db.Close()
}
