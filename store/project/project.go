package project

// Project tenant resource container
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	CreateAt    int64  `json:"create_at"`
	DomainID    string `json:"domain_id"`
	Extra       string `json:"-"`
}

// Store is project service
type Store interface {
	StoreReader
	StoreWriter
	Close() error
}

// StoreReader read project information from store
type StoreReader interface {
	GetProject(id string) (*Project, error)
	ListDomainProjects(domainID string) ([]*Project, error)
	CheckProjectIsExistByID(id string) (bool, error)
	ListProjectUsers(projectID string) ([]string, error)
}

// StoreWriter write project information to store
type StoreWriter interface {
	CreateProject(domainID, name, description string, enabled bool) (*Project, error)
	UpdateProject(id, name, description string) (*Project, error)
	// Soft Delete a Project,project still in persistence storage, super admin & domain admin are allowed to operate
	DeleteProject(id string) error
	// add users to project
	AddUsersToProject(projectID string, userIDs ...string) error
	// remove users from project
	RemoveUsersFromProject(projectID string, userIDs ...string) error
}
