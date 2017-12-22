package role

// Role is rbac's role
type Role struct {
	// Role id, UUID as a unique logo
	RoleID string
	// role name, allow repeat
	Name string
	// Role description
	Description string
	// Whether to enable
	Enabled bool
	// Extend fields to facilitate the expansion of database tables
	Extra string
}

// Manager is role service
type Manager interface {
	// Create a role with features, super admin only
	CreateRole(name, description string, features ...string) (*Role, error)
	// Associate features to roles
	AssociateFeaturesToRole(roleID string, features ...string) (bool, error)
	// Unlink feature to role
	// Note that when a character's list of properties is empty,
	// it should be a space, and can not use "" instead
	UnlinkFeatureFromRole(roleID string, features ...string) (bool, error)
	// Get a role with id, domain admin & super admin Only
	// Get the role information, it will automatically refresh
	// the role of the list of properties, filter to the offline features
	GetRole(id string) (*Role, error)
	GetRoleFeature(id string) ([]string, error)
	// List role, super admin only
	ListRole() (*[]Role, error)
	// Update role with id, modify name or description
	UpdateRole(roleID, name, description string) (*Role, error)
	// Verify that the role has permission to operate a function
	VerifyRole(roleID string, feature string) (bool, error)
	// Soft delete role, only target it delete,not real delete
	DeleteRole(roleID string) error
}
