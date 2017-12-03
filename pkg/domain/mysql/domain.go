package mysql

import (
	"database/sql"
	"sync"
	"time"

	"github.com/satori/go.uuid"

	"openauth/api/exception"
	"openauth/pkg/domain"
)

var (
	createPrepare *sql.Stmt
	deletePrepare *sql.Stmt
)

// NewDomainManager use to create domain storage service
func NewDomainManager(db *sql.DB) domain.Manager {
	return &manager{db: db}
}

// DomainManager is use mongodb as storage
type manager struct {
	db *sql.DB
}

// CreateDomain use to create an domain
func (m *manager) CreateDomain(name, description, displayName string, enabled bool) (*domain.Domain, error) {
	var (
		once sync.Once
		err  error
	)

	once.Do(func() {
		createPrepare, err = m.db.Prepare("INSERT INTO `domain` (id, name, display_name, description, enabled, extra, create_at) VALUES (?,?,?,?,?,?,?)")
	})
	if err != nil {
		return nil, exception.NewInternalServerError("prepare insert domain stmt error, domain: %s, %s", name, err)
	}

	dom := domain.Domain{ID: uuid.NewV4().String(), Name: name, DisplayName: displayName, Description: description, CreateAt: time.Now().Unix(), Enabled: enabled}
	_, err = createPrepare.Exec(dom.ID, dom.Name, dom.DisplayName, dom.Description, dom.Enabled, "", dom.CreateAt)
	if err != nil {
		return nil, exception.NewInternalServerError("insert domain exec sql err, %s", err)
	}
	return &dom, nil
}

// GetDomain use to get domain detail
func (m *manager) GetDomain(domainID string) (*domain.Domain, error) {
	dom := domain.Domain{}
	err := m.db.QueryRow("SELECT id,name,display_name,description,enabled,create_at,update_at FROM domain WHERE id = ?", domainID).Scan(
		&dom.ID, &dom.Name, &dom.DisplayName, &dom.Description, &dom.Enabled, &dom.CreateAt, &dom.UpdateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("domain %s not find", domainID)
		}

		return nil, exception.NewInternalServerError("query single domain error, %s", err)
	}

	return &dom, nil
}

// ListDomain use to list all domains
func (m *manager) ListDomain() ([]*domain.Domain, error) {
	rows, err := m.db.Query("SELECT id,name,display_name,description,enabled,create_at,update_at FROM domain")
	if err != nil {
		return nil, exception.NewInternalServerError("query domain list error, %s", err)
	}
	defer rows.Close()

	domains := []*domain.Domain{}
	for rows.Next() {
		dom := domain.Domain{}
		if err := rows.Scan(&dom.ID, &dom.Name, &dom.DisplayName, &dom.Description, &dom.Enabled, &dom.CreateAt, &dom.UpdateAt); err != nil {
			return nil, exception.NewInternalServerError("scan domain record error, %s", err)
		}
		domains = append(domains, &dom)
	}

	return domains, nil
}

// UpdateDomain use to update an domain
func (m *manager) UpdateDomain(id, name, description string) (*domain.Domain, error) {
	return nil, nil
}

// DeleteDomain use to delete an domain from db
func (m *manager) DeleteDomain(id string) error {
	var (
		once sync.Once
		err  error
	)

	once.Do(func() {
		deletePrepare, err = m.db.Prepare("DELETE FROM domain WHERE id = ?")
	})
	if err != nil {
		return exception.NewInternalServerError("prepare delete domain stmt error, %s", err)
	}

	if _, err := deletePrepare.Exec(id); err != nil {
		return exception.NewInternalServerError("delete domain exec sql error, %s", err)
	}

	return nil
}

func (m *manager) CheckDomainIsExist(domainID string) error {
	var id string
	if err := m.db.QueryRow("SELECT id FROM domain WHERE id = ?", domainID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return exception.NewNotFound("domain %s not find", domainID)
		}

		return exception.NewInternalServerError("query single domain error, %s", err)
	}

	return nil
}
