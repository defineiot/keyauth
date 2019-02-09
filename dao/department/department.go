package department

import (
	"github.com/defineiot/keyauth/dao/models"
)

// Store is user service
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader use to read department information form store
type Reader interface {
	GetDepartment(depID string) (*models.Department, error)
	GetDepartmentByName(domainID, departmentName string) (*models.Department, error)
	ListSubDepartments(domainID, parentDepID string) ([]*models.Department, error)
}

// Writer use to write department information to store
type Writer interface {
	CreateDepartment(d *models.Department) error
	DelDepartment(depID string) error
}
