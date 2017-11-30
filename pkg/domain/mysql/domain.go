package mysql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/satori/go.uuid"

	"openauth/pkg/domain"
)

var (
	createPrepare *sql.Stmt
	once          sync.Once
)

// NewDomainManager use to create domain storage service
func NewDomainManager(db *sql.DB) (domain.Manager, error) {
	return &manager{db: db}, nil
}

// DomainManager is use mongodb as storage
type manager struct {
	db *sql.DB
}

// CreateDomain use to create an domain
func (m *manager) CreateDomain(name, description, displayName string) (*domain.Domain, error) {
	var err error
	once.Do(func() {
		createPrepare, err = m.db.Prepare("INSERT INTO `domain` (id, name, display_name, description, enable, extra, create_at, update_at) VALUES (?,?,?,?,?,?,?)")
	})
	if err != nil {
		return nil, fmt.Errorf("insert domain: %s error, %s", name, err)
	}

	_, err = createPrepare.Exec(uuid.NewV4().String(), name, displayName, description, 1, "", time.Now().Unix(), 0)
	if err != nil {
		return nil, fmt.Errorf("insert device exec sql err, %s", err.Error())
	}
	return nil, nil
}

// GetDomain use to get domain detail
func (m *manager) GetDomain(domainID string) (*domain.Domain, error) {
	return nil, nil
}

// ListDomain use to list all domains
func (m *manager) ListDomain() (*[]domain.Domain, error) {
	return nil, nil
}

// UpdateDomain use to update an domain
func (m *manager) UpdateDomain(id, name, description string) (*domain.Domain, error) {
	return nil, nil
}

// DeleteDomain use to delete an domain from db
func (m *manager) DeleteDomain(id string) error {
	return nil
}
