package mysql

import (
	"database/sql"

	"openauth/api/exception"
	"openauth/store/domain"
	"openauth/tools"
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
			INSERT INTO domains (id, name, display_name, description, enabled, extra, create_at)
			VALUES (?,?,?,?,?,?,?);
		`,
		FindDomains: `
			SELECT d.id, d.name, d.display_name, d.description, d.enabled, d.create_at, d.update_at
			FROM domains d
			ORDER BY create_at 
			DESC;
		`,
		FindDomainsWithPage: `
			SELECT d.id, d.name, d.display_name, d.description, d.enabled, d.create_at, d.update_at
			FROM domains d
			ORDER BY create_at 
			DESC LIMIT ?,?;
		`,
		FindDomainByID: `
			SELECT d.id, d.name, d.display_name, d.description, d.enabled, d.create_at, d.update_at
			FROM domains d
			WHERE id = ?;
		`,
		FindDomainByName: `
			SELECT d.id, d.name, d.display_name, d.description, d.enabled, d.create_at, d.update_at
			FROM domains d
			WHERE name = ?;
		`,
		DeleteDomain: `
			DELETE FROM domains 
			WHERE id = ?;
		`,
		DomainCount: `
			SELECT COUNT(*) 
			FROM domains;
		`,
		FindDomainID: `
			SELECT id 
			FROM domains 
			WHERE id = ?;
		`,
		FindDomainName: `
			SELECT id 
			FROM domains 
			WHERE name = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(db, unprepared)
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
