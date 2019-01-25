package mysql

import (
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/internal/exception"
)

// CreateDomain use to create an domain
func (s *store) CreateDomain(d *domain.Domain) error {
	if err := d.Validate(); err != nil {
		return err
	}

	ok, err := s.CheckDomainIsExistByName(d.Name)
	if err != nil {
		return err
	}

	if ok {
		return exception.NewBadRequest("domain %s exist", d.Name)
	}

	d.ID = uuid.NewV4().String()
	d.CreateAt = time.Now().Unix()
	_, err = s.stmts[CreateDomain].Exec(
		d.ID, d.Name, d.DisplayName, d.LogoPath, d.Description, d.Enabled,
		int(d.Type), d.CreateAt, d.Size, d.Location, d.Industry, d.Address,
		d.Fax, d.Phone, d.ContactsName, d.ContactsTitle, d.ContactsMobile,
		d.ContactsEmail, d.Owner)
	if err != nil {
		return exception.NewInternalServerError("insert domain exec sql err, %s", err)
	}
	return nil
}

// GetDomain use to get domain detail
func (s *store) GetDomainByID(domainID string) (*domain.Domain, error) {
	d := new(domain.Domain)
	err := s.stmts[FindDomainByID].QueryRow(domainID).Scan(
		&d.ID, &d.Name, &d.DisplayName, &d.LogoPath, &d.Description, &d.Enabled,
		&d.Type, &d.CreateAt, &d.UpdateAt, &d.Size, &d.Location, &d.Industry,
		&d.Address, &d.Fax, &d.Phone, &d.ContactsName, &d.ContactsTitle,
		&d.ContactsMobile, &d.ContactsEmail, &d.Owner)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("domain %s not find", domainID)
		}

		return nil, exception.NewInternalServerError("query single domain error, %s", err)
	}

	return d, nil
}

func (s *store) ListUserThirdDomains(userID string) ([]*domain.Domain, error) {
	rows, err := s.stmts[FindUserThirdDomains].Query(userID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user third domains error, %s", err)
	}

	domains := []*domain.Domain{}
	for rows.Next() {
		d := new(domain.Domain)
		err := rows.Scan(
			&d.ID, &d.Name, &d.DisplayName, &d.LogoPath, &d.Description, &d.Enabled,
			&d.Type, &d.CreateAt, &d.UpdateAt, &d.Size, &d.Location, &d.Industry,
			&d.Address, &d.Fax, &d.Phone, &d.ContactsName, &d.ContactsTitle,
			&d.ContactsMobile, &d.ContactsEmail, &d.Owner)
		if err != nil {
			return nil, exception.NewInternalServerError("scan domain record error, %s", err)
		}
		domains = append(domains, d)
	}

	return domains, nil
}

func (s *store) GetDomainByName(name string) (*domain.Domain, error) {
	d := new(domain.Domain)
	err := s.stmts[FindDomainByName].QueryRow(name).Scan(
		&d.ID, &d.Name, &d.DisplayName, &d.LogoPath, &d.Description, &d.Enabled,
		&d.Type, &d.CreateAt, &d.UpdateAt, &d.Size, &d.Location, &d.Industry,
		&d.Address, &d.Fax, &d.Phone, &d.ContactsName, &d.ContactsTitle,
		&d.ContactsMobile, &d.ContactsEmail, &d.Owner)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("domain %s not find", name)
		}

		return nil, exception.NewInternalServerError("query single domain error, %s", err)
	}

	return d, nil
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
		d := new(domain.Domain)
		err := rows.Scan(
			&d.ID, &d.Name, &d.DisplayName, &d.LogoPath, &d.Description, &d.Enabled,
			&d.Type, &d.CreateAt, &d.UpdateAt, &d.Size, &d.Location, &d.Industry,
			&d.Address, &d.Fax, &d.Phone, &d.ContactsName, &d.ContactsTitle,
			&d.ContactsMobile, &d.ContactsEmail, &d.Owner)
		if err != nil {
			return nil, 0, exception.NewInternalServerError("scan domain record error, %s", err)
		}
		domains = append(domains, d)
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

// DeleteDomainByID use to delete an domain from db
func (s *store) DeleteDomainByID(id string) error {
	ret, err := s.stmts[DeleteDomainByID].Exec(id)
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

// DeleteDomainByName use to delete an domain from db
func (s *store) DeleteDomainByName(name string) error {
	ret, err := s.stmts[DeleteDomainByName].Exec(name)
	if err != nil {
		return exception.NewInternalServerError("delete domain exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("domian %s not exist", name)
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
