package mysql

import (
	"database/sql"

	"openauth/api/exception"
	"openauth/store/application"
)

const (
	CreateAPP   = "create-application"
	DeleteAPP   = "delete-app"
	GetUserAPPS = "get-user-apps"
	GetClient   = "get-client"

	CheckExistByID   = "check-exist-by-id"
	CheckExistByName = "check-exist-by-name"
)

// NewAppStore use to create domain storage service
func NewAppStore(db *sql.DB) (application.Store, error) {
	unprepared := map[string]string{
		CreateAPP: `
			INSERT INTO application (id, name, user_id, client_id, client_secret, client_type, website, logo_image, description, redirect_uri, create_at) 
			VALUES (?,?,?,?,?,?,?,?,?,?,?);
		`,
		GetUserAPPS: `
			SELECT a.id, a.name, a.user_id, a.client_id, a.client_secret, a.client_type, a.website, a.logo_image, a.description, a.redirect_uri, a.create_at 
			FROM application a 
			WHERE user_id = ? 
			ORDER BY create_at 
			DESC;
		`,
		GetClient: `
			SELECT a.client_id, a.client_secret, a.client_type, a.redirect_uri 
			FROM application a 
			WHERE client_id = ?;
		`,
		DeleteAPP: `
			DELETE FROM application 
			WHERE id = ?;
		`,
		CheckExistByID: `
			SELECT id FROM application 
			WHERE id = ?;
		`,
		CheckExistByName: `
			SELECT id FROM application 
			WHERE name = ? 
			AND user_id = ?;
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
