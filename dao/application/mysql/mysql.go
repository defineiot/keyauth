package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao/application"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	CreateAPP    = "create-application"
	CreateClient = "create-client"
	DeleteAPP    = "delete-app"
	DeleteClient = "delete-client"
	ListUserAPPS = "list-user-apps"
	GetUserAPP   = "get-user-app"
	GetClient    = "get-client"

	CheckExistByID   = "check-exist-by-id"
	CheckExistByName = "check-exist-by-name"
)

// NewAppStore use to create domain storage service
func NewAppStore(db *sql.DB) (application.Store, error) {
	unprepared := map[string]string{
		CreateAPP: `
			INSERT INTO applications (id, name, user_id, website, logo_image, description, create_at, client_id) 
			VALUES (?,?,?,?,?,?,?,?);
		`,
		ListUserAPPS: `
			SELECT a.id, a.name, a.user_id, a.website, a.logo_image, a.description, a.create_at, a.client_id 
			FROM applications a
			WHERE user_id = ? 
			ORDER BY a.create_at 
			DESC;
		`,
		GetUserAPP: `
		    SELECT a.id, a.name, a.user_id, a.website, a.logo_image, a.description, a.create_at, a.client_id
		    FROM applications a
		    WHERE a.id = ? 
		    ORDER BY a.create_at 
		    DESC;
	    `,
		DeleteAPP: `
			DELETE FROM applications 
			WHERE id = ?;
		`,
		CheckExistByID: `
			SELECT id FROM applications 
			WHERE id = ?;
		`,
		CheckExistByName: `
			SELECT id FROM applications 
			WHERE name = ? 
			AND user_id = ?;
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
