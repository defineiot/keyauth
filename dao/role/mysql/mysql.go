package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/log"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	SaveRole        = "save-role"
	FindAllRole     = "find-role"
	FindOneRole     = "find-role-by-name"
	DeleteRole      = "delete-role"
	CheckRole       = "check-role-exist"
	GetRoleFeatures = "get-role-features"
)

// NewRoleStore use to create domain storage service
func NewRoleStore(db *sql.DB, log log.IOTAuthLogger) (role.Store, error) {
	unprepared := map[string]string{
		SaveRole: `
			INSERT INTO roles (id, name, description, create_at) 
			VALUES (?, ?,?,?);
		`,
		FindAllRole: `
			SELECT id, name, description, create_at, update_at  
			FROM roles;
		`,
		FindOneRole: `
		    SELECT id, name, description, create_at, update_at 
			FROM roles
			WHERE id = ?;
	    `,
		DeleteRole: `
			DELETE FROM roles 
			WHERE id = ?;
		`,
		CheckRole: `
		    SELECT name 
		    FROM roles
		    WHERE name = ?;
		`,
		GetRoleFeatures: `
			SELECT f.id, f.name, f.method, f.endpoint, f.description, f.is_deleted, f.when_deleted_version, f.is_added, f.when_added_version, f.service_id
			FROM features f 
			LEFT JOIN role_feature_mappings m
			ON m.feature_id = f.id 
			WHERE m.role_id = ?; 
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(db, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare token store query statment error, %s", err)
	}

	s := store{
		db:    db,
		stmts: stmts,
		log:   log,
	}

	return &s, nil
}

// DomainManager is use mongodb as storage
type store struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
	log   log.IOTAuthLogger
}

// Close closes the database, releasing any open resources.
func (s *store) Close() error {
	return s.db.Close()
}
