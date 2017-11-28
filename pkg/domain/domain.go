package domain

// Domain is tenant container.
type Domain struct {
	DomainID string
	// domain name, allow repeat
	Name string
	// DisplayName is to show
	DisplayName string
	// domain description
	Description string
	// Whether to enable
	Enabled bool
	// Extend fields to facilitate the expansion of database tables
	Extra string
	// CreateAt create domain at
	CreateAt int64
	// UpdateAt update domain time
	UpdateAt int64
}

// Manager is an domain service
type Manager interface {
	// Create Domain, Only super admin are allowed
	// to operate, Named globally only,
	// renaming is not allowed
	CreateDomain(name, description, displayName string) (*Domain, error)
	// GetDomain get a domain by domain id or domain name,
	// super admin & domain admin are allowed to operate
	GetDomain(domainID string) (*Domain, error)
	// List all Domain, Only super admin are allowed to operate
	ListDomain() (*[]Domain, error)
	// Update a Domain, super admin & domain admin are allowed to operate
	UpdateDomain(id, name, description string) (*Domain, error)
	// Soft Delete a Domain, Domain still in persistence storage, Only super admin are allowed to operate
	DeleteDomain(id string) error
}
