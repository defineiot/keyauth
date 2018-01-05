package mysql

import (
	"database/sql"

	"openauth/api/exception"
	"openauth/store/token"
	"openauth/tools"
)

const (
	SaveToken   = "save-token"
	FindToken   = "find-token"
	DeleteToken = "delete-token"
)

// NewTokenStore use to create domain storage service
func NewTokenStore(db *sql.DB) (token.Store, error) {
	unprepared := map[string]string{
		SaveToken: `
			INSERT INTO tokens (grant_type, access_token, refresh_token, type, create_at, expire_at, client_id, user_id, domain_id, project_id) 
			VALUES (?,?,?,?,?,?,?,?,?,?);
		`,
		FindToken: `
			SELECT t.grant_type, t.access_token, t.refresh_token, t.type, t.create_at, t.expire_at, t.client_id, t.user_id, t.domain_id, t.project_id 
			FROM tokens t
			WHERE access_token = ?;
		`,
		DeleteToken: `
			DELETE FROM tokens 
			WHERE access_token = ?;
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
