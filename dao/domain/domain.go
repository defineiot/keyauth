package domain

// Domain is tenant container.
type Domain struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	CreateAt    int64  `json:"create_at"`
	UpdateAt    int64  `json:"update_at,omitempty"`
	Extra       string `json:"-"`
}

// Store is an domain service
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader for read data from store
type Reader interface {
	GetDomainByID(domainID string) (*Domain, error)
	GetDomainByName(name string) (*Domain, error)
	CheckDomainIsExistByID(domainID string) (bool, error)
	CheckDomainIsExistByName(domainName string) (bool, error)
	ListDomain(pageNumber, pageSize int64) (domains []*Domain, totalPage int64, err error)
}

// Writer for write data to store
type Writer interface {
	CreateDomain(name, description, displayName string, enabled bool) (*Domain, error)
	UpdateDomain(id, name, description string) (*Domain, error)
	DeleteDomainByID(id string) error
	DeleteDomainByName(name string) error
}
