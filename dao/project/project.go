package project

import "github.com/defineiot/keyauth/dao/models"

// Store is project service
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader read project information from store
type Reader interface {
	GetProjectByID(id string) (*models.Project, error)
	ListDomainProjects(domainID string) ([]*models.Project, error)
	ListDepartmentProjects(departmentID string) ([]*models.Project, error)
	ListUserProjects(domainID, userID string) ([]*models.Project, error)
	CheckProjectIsExistByID(id string) (bool, error)
}

// Writer write project information to store
type Writer interface {
	CreateProject(p *models.Project) error
	DeleteProjectByID(id string) error
	DeleteProjectByName(domainName, prjectName string) error
	AddUsersToProject(projectID string, userIDs ...string) error
	RemoveUsersFromProject(projectID string, userIDs ...string) error
}
