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
			INSERT INTO roles (name, description, create_at) 
			VALUES (?,?,?);
		`,
		FindAllRole: `
			SELECT id, name, description, create_at 
			FROM roles;
		`,
		FindOneRole: `
		    SELECT id, name, description, create_at 
			FROM roles
			WHERE name = ?;
	    `,
		DeleteRole: `
			DELETE FROM roles 
			WHERE name = ?;
		`,
		CheckRole: `
		    SELECT name 
		    FROM roles
		    WHERE name = ?;
		`,
		GetRoleFeatures: `
			SELECT feature_id
			FROM roles_features_mapping
			WHERE role_name = ?; 
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
