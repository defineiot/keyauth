package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao/client"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	CreateClient = "create-client"
	UpdateClient = "update-client"
	DeleteClient = "delete-client"
	FindAll      = "find-client-all"
	FindOneByID  = "find-client-one"
)

// NewClientStore use to create domain storage service
func NewClientStore(db *sql.DB) (client.Store, error) {
	unprepared := map[string]string{
		CreateClient: `
		    INSERT INTO clients (id, secret, type, redirect_uri)
		    VALUES (?,?,?,?)
		`,
		FindOneByID: `
		    SELECT c.id, c.secret, c.type, c.redirect_uri, create_at 
		    FROM clients c
		    WHERE c.id = ?;
	    `,
		FindAll: `
		    SELECT c.id, c.secret, c.type, c.redirect_uri 
		    FROM clients c;
		`,
		DeleteClient: `
		    DELETE FROM clients
		    WHERE id = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(db, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare service store query statment error, %s", err)
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
