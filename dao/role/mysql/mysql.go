package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/logger"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	SaveRole       = "save-role"
	FindAllRole    = "find-role"
	FindRoleByName = "find-role-by-name"
	FindRoleByID   = "find-role-by-id"
	DeleteRole     = "delete-role"
	CheckRole      = "check-role-exist"
	FindRoleByUser = "find-role-by-user"
)

// NewRoleStore use to create domain storage service
func NewRoleStore(opt *dao.Options) (role.Store, error) {
	unprepared := map[string]string{
		SaveRole: `
			INSERT INTO roles (id, name, description, create_at) 
			VALUES (?, ?,?,?);
		`,
		FindAllRole: `
			SELECT id, name, description, create_at, update_at  
			FROM roles;
		`,
		FindRoleByUser: `
			SELECT id, name, description, create_at, update_at  
			FROM roles r 
			LEFT JOIN role_user_mappings m
			ON r.id = m.role_id
			WHERE domain_id = ?  
			AND user_id = ?;
		`,
		FindRoleByName: `
		    SELECT id, name, description, create_at, update_at 
			FROM roles
			WHERE name = ?;
		`,
		FindRoleByID: `
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
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(opt.DB, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare token store query statment error, %s", err)
	}

	s := store{
		db:    opt.DB,
		stmts: stmts,
		log:   opt.LOG,
	}

	return &s, nil
}

// DomainManager is use mongodb as storage
type store struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
	log   logger.Logger
}

// Close closes the database, releasing any open resources.
func (s *store) Close() error {
	return s.db.Close()
}

func init() {
	dao.Registe(NewRoleStore)
}
