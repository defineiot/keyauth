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
	DeleteClient  = "delete-client"
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
		FindAll: `
			SELECT s.id, s.name, s.description, s.enabled, s.status, s.status_update_at, s.version, s.create_at, c.id, c.secret
			FROM services s
			LEFT JOIN clients c
			ON s.id = c.service_id
			WHERE c.type = "confidential"
			ORDER BY s.create_at
			DESC;
		`,
		FindOneByID: `
			SELECT s.id, s.name, s.description, s.enabled, s.status, s.status_update_at, s.version, s.create_at, c.id, c.secret
			FROM services s
			LEFT JOIN clients c
			ON s.id = c.service_id
			WHERE c.type = "confidential"
			AND s.id = ?;
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
