package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	CreateDomain       = "create-domain"
	DeleteDomainByID   = "delete-domain-by-id"
	DeleteDomainByName = "delete-domain-by-name"
	UpdateDomainByID   = "update-domain-by-id"

	FindDomains         = "find-domains"
	FindDomainsWithPage = "find-domains-with-page"
	FindDomainByID      = "find-domain-by-id"
	FindDomainByName    = "find-domain-by-name"
	DomainCount         = "domain-count"
	FindDomainID        = "find-domain-id"
	FindDomainName      = "find-domain-name"
)

// NewDomainStore use to create domain storage service
func NewDomainStore(opt *dao.Options) (domain.Store, error) {
	if opt.DB == nil {
		return nil, exception.NewInternalServerError("miss db connection")
	}

	unprepared := map[string]string{
		CreateDomain: `
			INSERT INTO domains (id, name, display_name, logo_path, description, enabled, type, create_at, size, location, industry, address, fax, phone, contacts_name, contacts_title, contacts_mobile, contacts_email, owner_id)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
		`,
		FindDomainByID: `
			SELECT id, name, display_name, logo_path, description, enabled, type, create_at, update_at, size, location, industry, address, fax, phone, contacts_name, contacts_title, contacts_mobile, contacts_email, owner_id
			FROM domains 
			WHERE id = ?;
		`,
		FindDomains: `
			SELECT id, name, display_name, logo_path, description, enabled, type, create_at, update_at, size, location, industry, address, fax, phone, contacts_name, contacts_title, contacts_mobile, contacts_email, owner_id 
			FROM domains
			ORDER BY create_at 
			DESC;
		`,
		FindDomainsWithPage: `
			SELECT id, name, display_name, logo_path, description, enabled, type, create_at, update_at, size, location, industry, address, fax, phone, contacts_name, contacts_title, contacts_mobile, contacts_email, owner_id 
			FROM domains
			ORDER BY create_at 
			DESC LIMIT ?,?;
		`,
		FindDomainByName: `
			SELECT id, name, display_name, logo_path, description, enabled, type, create_at, update_at, size, location, industry, address, fax, phone, contacts_name, contacts_title, contacts_mobile, contacts_email, owner_id
			FROM domains 
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
	stmts, err := tools.PrepareStmts(opt.DB, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare domain query statment error, %s", err)
	}

	s := store{
		db:    opt.DB,
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

func init() {
	dao.Registe(NewDomainStore)
}
