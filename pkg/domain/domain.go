package domain

// Domain is tenant container.
type Domain struct {
	ID string `json:"id"`
	// domain name, allow repeat
	Name string `json:"name"`
	// DisplayName is to show
	DisplayName string `json:"display_name"`
	// domain description
	Description string `json:"description"`
	// Whether to enable
	Enabled bool `json:"enabled"`
	// Extend fields to facilitate the expansion of database tables
	Extra string `json:"-"`
	// CreateAt create domain at
	CreateAt int64 `json:"create_at"`
	// UpdateAt update domain time
	UpdateAt int64 `json:"update_at,omitempty"`
}

// Manager is an domain service
type Manager interface {
	// Create Domain, Only super admin are allowed
	// to operate, Named globally only,
	// renaming is not allowed
	CreateDomain(name, description, displayName string, enabled bool) (*Domain, error)
	// GetDomain get a domain by domain id or domain name,
	// super admin & domain admin are allowed to operate
	GetDomain(domainID string) (*Domain, error)
	// List all Domain, Only super admin are allowed to operate
	ListDomain() ([]*Domain, error)
	// Update a Domain, super admin & domain admin are allowed to operate
	UpdateDomain(id, name, description string) (*Domain, error)
	// Soft Delete a Domain, Domain still in persistence storage, Only super admin are allowed to operate
	DeleteDomain(id string) error
	// CheckDomainIsExist use to check the domain is exist by domain id
	CheckDomainIsExistByID(domainID string) (bool, error)
}
