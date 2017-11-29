package user

import (
	"time"
)

// User info
type User struct {
	// user id, UUID as a unique logo
	ID string
	// user name, not allow repeat
	Name string
	// Default Project ID
	DefaultProjectID string
	// Whether to enable
	Enabled bool
	// Extend fields to facilitate the expansion of database tables
	Extra string
}

// Phone user's phone
type Phone struct {
	ID          uint
	UserID      string
	Number      string
	Type        string
	Description string
	// Extend fields to facilitate the expansion of database tables
	Extra string
}

// Email use's email
type Email struct {
	ID          uint
	UserID      string
	Address     string
	Primary     bool
	Description string
	// Extend fields to facilitate the expansion of database tables
	Extra string
}

// Password user's password
type Password struct {
	UserID    string
	Password  string
	ExpiredAt time.Time
	// Extend fields to facilitate the expansion of database tables
	Extra string
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
	CreateUser(projectID, userName, password string) (*User, error)
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
