package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	CreateDomain        = "create-domain"
	DeleteDomainByID    = "delete-domain-by-id"
	DeleteDomainByName  = "delete-domain-by-name"
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
		DeleteDomainByID: `
			DELETE FROM domains 
			WHERE id = ?;
		`,
		DeleteDomainByName: `
			DELETE FROM domains 
			WHERE name = ?;
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
