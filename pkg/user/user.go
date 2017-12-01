package user

// User info
type User struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	LastActiveTime   int64    `json:"last_active_time"`
	Enabled          bool     `json:"enabled"`
	CreateAt         int64    `json:"create_at"`
	Password         Password `json:"password"`
	Phones           []*Phone `json:"phones,omitempty"`
	Emails           []*Email `json:"emails,omitempty"`
	DomainID         string   `json:"domain_id"`
	DefaultProjectID string   `json:"default_project_id"`
	ExpireActiveDays int64    `json:"expires_active_days"`

	// Extend fields to facilitate the expansion of database tables
	Extra string `json:"-"`
}

// Phone user's phone
type Phone struct {
	ID          int64  `json:"-"`
	UserID      string `json:"-"`
	Number      string `json:"number"`
	Primary     bool   `json:"primary"`
	Description string `json:"descrption"`
	// Extend fields to facilitate the expansion of database tables
	Extra string `json:"-"`
}

// Email use's email
type Email struct {
	ID          int64  `json:"-"`
	UserID      string `json:"-"`
	Address     string `json:"address"`
	Primary     bool   `json:"primary"`
	Description string `json:"description"`
	// Extend fields to facilitate the expansion of database tables
	Extra string `json:"-"`
}

// Password user's password
type Password struct {
	ID       int64  `json:"-"`
	Password string `json:"-"`
	ExpireAt int64  `json:"expire_at"`
	CreateAt int64  `json:"create_at"`
	UpdateAt int64  `json:"update_at,omitempty"`

	// Extend fields to facilitate the expansion of database tables
	Extra string `json:"-"`
}

// Credential is user's credential
type Credential struct {
	UserID      string
	Password    string
	AccessToken string
	DomainID    string
	UserName    string
}

// Manager is user service
type Manager interface {
	// Create user, in the same domain, user can not be renamed
	CreateUser(domainID, projectID, userName, password string, enabled bool, userExpires, passExpires int) (*User, error)
	// Get user with User name & user password
	GetUserByName(domainID, userName, userPassword string) (*User, error)
	// Get user with User id & user password
	GetUserByID(userID, userPassword string) (*User, error)
	// GetUser get an user
	GetUser(cert Credential) (*User, error)
	// Delete user from persistence storage
	DeleteUser(cert Credential) error
	// Add Phone to User
	AddPhone(cert Credential, number, phoneType, description string) error
	// Remove Phone from User
	RemovePhone(cert Credential, number string) error
	// Query user phones
	QueryPhone(cert Credential) (*[]Phone, error)
	// Add email to User
	AddEmail(cert Credential, address, description string, primary bool) error
	// Remove Email from User
	RemoveEmail(cert Credential, address string) error
	// Query user emails
	QueryEmail(cert Credential) (*[]Email, error)
	// Add role to user
	AddRoleToUser(cert Credential, roleID string) error
	// remove role from user
	RemoveRoleFromUser(cert Credential, roleID string) error
	// list user roles
	QueryRole(cert Credential) ([]string, error)
	// Verify user & feature
	HasFeatures(cert Credential, features ...string) (bool, error)
	// set user default project
	SetDefaultProject(cert Credential, projectID string) error
	// add user to project, a user can add multiple project
	AddUserToProject(cert Credential, projectID string) error
	// remove user from project, when the user does not belong to any project, the user disable
	RemoveUserFromProject(cert Credential, projectID string) error
	// list user projects
	ListUserProject(cert Credential) ([]string, error)
	// get user default project
	GetDefaultProject(cert Credential) (string, error)
	// determine if the user has a super administrator role
	IsSystemAdmin(cert Credential) (bool, error)
	// determine if the user has a cloud administrator role
	IsCloudAdmin(cert Credential) (bool, error)
}
