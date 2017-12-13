package project

// Project tenant resource container
type Project struct {

	// Project id, UUID as a unique logo
	ID string `json:"id"`
	// Project name, allow repeat
	Name string `json:"name"`
	// Project description
	Description string `json:"description"`
	// Whether to enable
	Enabled bool `json:"enabled"`
	// CrateAt create time
	CreateAt int64 `json:"create_at"`
	// the project's domain id
	DomainID string `json:"domain_id"`
	// Extend fields to facilitate the expansion of database tables
	Extra string `json:"-"`
}

// Service is project service
type Service interface {
	// Create a Project, super admin & domain admin are
	// allowed to operate, Named in Domain, does not allow renaming
	CreateProject(domainID, name, description string, enabled bool) (*Project, error)

	// Get a project with project id
	GetProject(id string) (*Project, error)
	// List all Project in domain_id, else all project
	ListDomainProjects(domainID string) ([]*Project, error)
	// Update a Project, super admin & domain admin are allowed to operate
	UpdateProject(id, name, description string) (*Project, error)
	// Soft Delete a Project,project still in persistence storage, super admin & domain admin are allowed to operate
	DeleteProject(id string) error
	// CheckProjectIsExist use to check the project is exist by project id
	CheckProjectIsExistByID(id string) (bool, error)
	// ListProjectUsers use to list all users
	ListProjectUsers(projectID string) ([]string, error)
	// add users to prject
	AddUsersToProject(projectID string, userIDs ...string) error
	// remove users from project
	RemoveUsersFromProject(projectID string, userIDs ...string) error
}
