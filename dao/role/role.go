package role

import (
	"github.com/defineiot/keyauth/dao/service"
)

// Role is rbac's role
type Role struct {
	ID          int64              `json:"-"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	CreateAt    int64              `json:"create_at"`
	Featrues    []*service.Feature `json:"features"`
}

// Store is an role service
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader for read data from store
type Reader interface {
	// Get a role with id, domain admin & super admin Only
	// Get the role information, it will automatically refresh
	// the role of the list of properties, filter to the offline features
	GetRole(name string) (*Role, error)
	GetRoleFeature(name string) ([]int64, error)
	CheckRoleExist(name string) (bool, error)
	// List role, super admin only
	ListRole() ([]*Role, error)
	// Verify that the role has permission to operate a function
	VerifyRole(name string, feature string) (bool, error)
}

// Writer for write data to store
type Writer interface {
	// Create a role with features, super admin only
	CreateRole(name, description string) (*Role, error)
	// Associate features to roles
	AssociateFeaturesToRole(name string, features ...int64) error
	// Unlink feature to role
	// Note that when a character's list of properties is empty,
	// it should be a space, and can not use "" instead
	UnlinkFeatureFromRole(name string, features ...int64) (bool, error)
	// Update role with id, modify name or description
	UpdateRole(name, description string) (*Role, error)

	// Soft delete role, only target it delete,not real delete
	DeleteRole(name string) error
}
