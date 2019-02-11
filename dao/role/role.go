package role

import (
	"github.com/defineiot/keyauth/dao/models"
)

// Store is an role service
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader for read data from store
type Reader interface {
	GetRole(id string) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	CheckRoleExist(name string) (bool, error)
	ListRole() ([]*models.Role, error)
	ListDepartmentRoles(departmentID string) ([]*models.Role, error)
	ListUserRole(domainID, userID string) ([]*models.Role, error)
}

// Writer for write data to store
type Writer interface {
	CreateRole(role *models.Role) error
	UpdateRole(name, description string) (*models.Role, error)
	DeleteRole(name string) error
}
