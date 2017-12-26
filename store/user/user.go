package user

import (
	"openauth/store/project"
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

	Extra string `json:"-"`
}

// Phone user's phone
type Phone struct {
	ID          int64  `json:"-"`
	UserID      string `json:"-"`
	Number      string `json:"number"`
	Primary     bool   `json:"primary"`
	Description string `json:"descrption"`

	Extra string `json:"-"`
}

// Email use's email
type Email struct {
	ID          int64  `json:"-"`
	UserID      string `json:"-"`
	Address     string `json:"address"`
	Primary     bool   `json:"primary"`
	Description string `json:"description"`

	Extra string `json:"-"`
}

// Password user's password
type Password struct {
	ID       int64  `json:"-"`
	Password string `json:"-"`
	ExpireAt int64  `json:"expire_at"`
	CreateAt int64  `json:"create_at"`
	UpdateAt int64  `json:"update_at,omitempty"`
	UserID   string `json:"-"`

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

// Store is user service
type Store interface {
	StoreReader
	StoreWriter
	Close() error
}

// StoreReader use to read user information form store
type StoreReader interface {
	ListUser(domainID string) ([]*User, error)
	GetUserByID(userID string) (*User, error)
	GetUserByName(domainID, userName string) (*User, error)
	ValidateUser(domainID, userName, password string) (string, error)
	CheckUserNameIsExist(domainID, userName string) (bool, error)
	CheckUserIsExistByID(userID string) (bool, error)

	ListUserProjects(userID string) ([]string, error)
}

// StoreWriter use to write user information to store
type StoreWriter interface {
	CreateUser(domainID, userName, password string, enabled bool, userExpires, passExpires int) (*User, error)
	DeleteUser(userID string) error

	SetDefaultProject(userID, projectID string) error
	AddProjectsToUser(userID string, projectIDs ...string) error
	RemoveProjectsFromUser(userID string, projectIDs ...string) error
}
