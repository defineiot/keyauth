package mysql

import (
	"database/sql"
	"fmt"

	"openauth/api/exception"
	"openauth/store/domain"
)

const (
	CreateDomain        = "create-domain"
	DeleteDomain        = "delete-domain"
	UpdateDomainByID    = "update-domain-by-id"
	FindDomains         = "find-domains"
	FindDomainsWithPage = "find-domains-with-page"
	FindDomainByID      = "find-domain-by-id"
	FindDomainByName    = "find-domain-by-name"

	DomainCount    = "domain-count"
	FindDomainID   = "find-domain-id"
	FindDomainName = "find-domain-name"
)

// NewDomainStore use to create domain storage service
func NewDomainStore(db *sql.DB) (domain.Store, error) {
	unprepared := map[string]string{
		CreateDomain: `
			INSERT INTO domain (id, name, display_name, description, enabled, extra, create_at)
			VALUES (?,?,?,?,?,?,?);
		`,
		FindDomains: `
			SELECT d.id, d.name, d.display_name, d.description, d.enabled, d.create_at, d.update_at
			FROM domain d
			ORDER BY create_at 
			DESC;
		`,
		FindDomainsWithPage: `
			SELECT d.id, d.name, d.display_name, d.description, d.enabled, d.create_at, d.update_at
			FROM domain d
			ORDER BY create_at 
			DESC LIMIT ?,?;
		`,
		FindDomainByID: `
			SELECT d.id, d.name, d.display_name, d.description, d.enabled, d.create_at, d.update_at
			FROM domain d
			WHERE id = ?;
		`,
		FindDomainByName: `
			SELECT d.id, d.name, d.display_name, d.description, d.enabled, d.create_at, d.update_at
			FROM domain d
			WHERE name = ?;
		`,
		DeleteDomain: `
			DELETE FROM domain 
			WHERE id = ?;
		`,
		DomainCount: `
			SELECT COUNT(*) 
			FROM domain;
		`,
		FindDomainID: `
			SELECT id 
			FROM domain 
			WHERE id = ?;
		`,
		FindDomainName: `
			SELECT id 
			FROM domain 
			WHERE name = ?;
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
			return nil, fmt.Errorf("prepare statment: %s, %s", k, err)
		}
		prepared[k] = stmt
	}

	return prepared, nil
}
