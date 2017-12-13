package user

import (
	"openauth/pkg/project"
)

// User info
type User struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	LastActiveTime   int64            `json:"last_active_time"`
	Enabled          bool             `json:"enabled"`
	CreateAt         int64            `json:"create_at"`
	Password         *Password        `json:"-"`
	Phones           []*Phone         `json:"phones,omitempty"`
	Emails           []*Email         `json:"emails,omitempty"`
	DomainID         string           `json:"domain_id"`
	DefaultProjectID string           `json:"-"`
	DefaultProject   *project.Project `json:"default_project,omitempty"`
	ExpireActiveDays int64            `json:"expires_active_days"`

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
	UserID   string `json:"user_id"`

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

// Service is user service
type Service interface {
	CreateUser(domainID, userName, password string, enabled bool, userExpires, passExpires int) (*User, error)
	ListUser(domainID string) ([]*User, error)
	GetUserByID(userID string) (*User, error)
	DeleteUser(userID string) error
	CheckUserNameIsExist(domainID, userName string) (bool, error)
	CheckUserIsExistByID(userID string) (bool, error)

	ListUserProjects(userID string) ([]string, error)
	SetDefaultProject(userID, projectID string) error
	AddProjectsToUser(userID string, projectIDs ...string) error
	RemoveProjectsFromUser(userID string, projectIDs ...string) error

	// Get user with User name & user password
	GetUserByName(domainID, userName, userPassword string) (*User, error)
	// GetUser get an user
	GetUser(cert Credential) (*User, error)
	// Delete user from persistence storage

	// Add Phone to User
	AddPhone(cert Credential, number, phoneType, description string) error
	// Remove Phone from User
	RemovePhone(cert Credential, number string) error
	// Add email to User
	AddEmail(cert Credential, address, description string, primary bool) error
	// Remove Email from User
	RemoveEmail(cert Credential, address string) error
	// Add role to user
	AddRoleToUser(cert Credential, roleID string) error
	// remove role from user
	RemoveRoleFromUser(cert Credential, roleID string) error
	// list user roles
	QueryRole(cert Credential) ([]string, error)
	// Verify user & feature
	HasFeatures(cert Credential, features ...string) (bool, error)
	// add user to project, a user can add multiple project
	AddUserToProject(cert Credential, projectID string) error
	// remove user from project, when the user does not belong to any project, the user disable
	RemoveUserFromProject(cert Credential, projectID string) error
	// get user default project
	GetDefaultProject(cert Credential) (string, error)
	// determine if the user has a super administrator role
	IsSystemAdmin(cert Credential) (bool, error)
	// determine if the user has a cloud administrator role
	IsCloudAdmin(cert Credential) (bool, error)
}
