package mysql

import (
	"database/sql"

	"openauth/api/exception"
	"openauth/api/logger"
	"openauth/store/user"
)

const (
	SaveToken = "save-token"
	FindToken = "find-token"
)

// NewUserStore use to create domain storage service
func NewUserStore(db *sql.DB, key string, logger logger.OpenAuthLogger) (user.Store, error) {
	unprepared := map[string]string{
		SaveToken: `
			INSERT INTO token (grant_type, access_token, refresh_token, type, create_at, expire_at, client_id, user_id, domain_id, project_id) 
			VALUES (?,?,?,?,?,?,?,?,?,?);
		`,
		FindToken: `
			SELECT t.grant_type, t.access_token, t.refresh_token, t.type, t.create_at, t.expire_at, t.client_id, t.user_id, t.domain_id, t.project_id 
			FROM token t
			WHERE access_token = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := prepareStmts(db, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare domain query statment error, %s", err)
	}

	s := store{
		db:    db,
		stmts: stmts,
		key:   key,
		log:   logger,
	}

	return &s, nil
}

// DomainManager is use mongodb as storage
type store struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
	key   string
	log   logger.OpenAuthLogger
}

// Close closes the database, releasing any open resources.
func (s *store) Close() error {
	return s.db.Close()
}

// prepareStmts will attempt to prepare each unprepared
// query on the database. If one fails, the function returns
// with an error.
func prepareStmts(db *sql.DB, unprepared map[string]string) (map[string]*sql.Stmt, error) {
	prepared := map[string]*sql.Stmt{}
	for k, v := range unprepared {
		stmt, err := db.Prepare(v)
		if err != nil {
			return nil, err
		}
		prepared[k] = stmt
	}

	return prepared, nil
}
