package mysql

import (
	"database/sql"
	"time"

	"github.com/satori/go.uuid"

	"openauth/api/exception"
	"openauth/storage/domain"
)

// CreateDomain use to create an domain
func (s *store) CreateDomain(name, description, displayName string, enabled bool) (*domain.Domain, error) {
	ok, err := s.CheckDomainIsExistByName(name)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, exception.NewBadRequest("domain %s exist", name)
	}

	dom := domain.Domain{ID: uuid.NewV4().String(), Name: name, DisplayName: displayName, Description: description, CreateAt: time.Now().Unix(), Enabled: enabled}
	_, err = s.stmts[CreateDomain].Exec(dom.ID, dom.Name, dom.DisplayName, dom.Description, dom.Enabled, "", dom.CreateAt)
	if err != nil {
		return nil, exception.NewInternalServerError("insert domain exec sql err, %s", err)
	}
	return &dom, nil
}

// GetDomain use to get domain detail
func (s *store) GetDomain(domainID string) (*domain.Domain, error) {
	dom := domain.Domain{}
	err := s.stmts[FindDomainByID].QueryRow(domainID).Scan(
		&dom.ID, &dom.Name, &dom.DisplayName, &dom.Description, &dom.Enabled, &dom.CreateAt, &dom.UpdateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("domain %s not find", domainID)
		}

		return nil, exception.NewInternalServerError("query single domain error, %s", err)
	}

	return &dom, nil
}

func (s *store) GetDomainByName(name string) (*domain.Domain, error) {
	dom := domain.Domain{}
	err := s.stmts[FindDomainByName].QueryRow(name).Scan(
		&dom.ID, &dom.Name, &dom.DisplayName, &dom.Description, &dom.Enabled, &dom.CreateAt, &dom.UpdateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("domain %s not find", name)
		}

		return nil, exception.NewInternalServerError("query single domain error, %s", err)
	}

	return &dom, nil
}

// ListDomain use to list all domains
func (s *store) ListDomain(pageNumber, pageSize int64) ([]*domain.Domain, int64, error) {
	var (
		rows   *sql.Rows
		err    error
		totalP int64
	)

	offset := (pageNumber - 1) * pageSize
	limit := pageSize

	if pageSize != 0 {
		rows, err = s.stmts[FindDomainsWithPage].Query(offset, limit)
	} else {
		rows, err = s.stmts[FindDomains].Query()
	}
	if err != nil {
		return nil, 0, exception.NewInternalServerError("query domain list error, %s", err)
	}
	defer rows.Close()

	domains := []*domain.Domain{}
	for rows.Next() {
		dom := domain.Domain{}
		if err := rows.Scan(&dom.ID, &dom.Name, &dom.DisplayName, &dom.Description, &dom.Enabled, &dom.CreateAt, &dom.UpdateAt); err != nil {
			return nil, 0, exception.NewInternalServerError("scan domain record error, %s", err)
		}
		domains = append(domains, &dom)
	}

	total, err := s.domainCount()
	if err != nil {
		return nil, 0, err
	}

	if pageSize != 0 {
		if total < pageSize {
			totalP = 1
		} else {
			ok := total % pageSize
			totalP = total / pageSize
			if ok != 0 {
				totalP++
			}
		}
	} else {
		totalP = 1
	}

	return domains, totalP, nil
}

func (s *store) domainCount() (int64, error) {
	count := int64(0)
	if err := s.stmts[DomainCount].QueryRow().Scan(&count); err != nil {
		return 0, exception.NewInternalServerError("count domain record error, %s", err)
	}

	return count, nil
}

// UpdateDomain use to update an domain
func (s *store) UpdateDomain(id, name, description string) (*domain.Domain, error) {
	return nil, nil
}

// DeleteDomain use to delete an domain from db
func (s *store) DeleteDomain(id string) error {
	ret, err := s.stmts[DeleteDomain].Exec(id)
	if err != nil {
		return exception.NewInternalServerError("delete domain exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("domian %s not exist", id)
	}

	return nil
}

func (s *store) CheckDomainIsExistByID(domainID string) (bool, error) {
	var id string
	if err := s.stmts[FindDomainID].QueryRow(domainID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("query single domain error, %s", err)
	}

	return true, nil
}

func (s *store) CheckDomainIsExistByName(domainName string) (bool, error) {
	var id string
	if err := s.stmts[FindDomainName].QueryRow(domainName).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("query single domain error, %s", err)
	}

	return true, nil
}
