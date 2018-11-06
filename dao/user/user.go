package user

import (
	"github.com/defineiot/keyauth/dao/project"
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
	RoleNames        []string         `json:"roles"`
	DomainID         string           `json:"domain_id"`
	DefaultProjectID string           `json:"-"`
	DefaultProject   *project.Project `json:"home_project,omitempty"`
	ExpireActiveDays int64            `json:"expires_active_days"`

	Extra string `json:"-"`
}

// VerifyCode code
type VerifyCode struct {
	ID           int64  `json:"id"`
	EmailAddress string `json:"email_address"`
	PhoneNumber  string `json:"phone_number"`
	Code         int    `json:"code"`
	CreateAt     int64  `json:"create_at"`
	ExpireAt     int64  `json:"expire_at"`
	Status       int    `json:"status"`
	Type         int    `json:"type"`
}

// InvitationCode code
type InvitationCode struct {
	ID                   int64    `json:"-"`
	InviterID            string   `json:"inviter_id"`
	InvitedUserID        string   `json:"invited_user_id,omitempty"`
	InvitedUserDomainID  string   `json:"invited_user_domain_id,omitempty"`
	InvitedTime          int64    `json:"invited_time"`
	AcceptTime           int64    `json:"accept_time,omitempty"`
	ExpireTime           int64    `json:"expire_time,omitempty"`
	Code                 string   `json:"code"`
	InvitationURL        string   `json:"invitation_url"`
	InvitedUserRoleNames []string `json:"invited_user_role_names"`
	AccessProjects       []string `json:"access_project_ids"`
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
	Reader
	Writer
	Close() error
}

// Reader use to read user information form store
type Reader interface {
	ListUser(domainID string) ([]*User, error)
	ListUserRoles(domainID, userID string) ([]string, error)
	GetUserByID(userID string) (*User, error)
	GetUserByName(domainID, userName string) (*User, error)
	ValidateUser(domainID, userName, password string) (string, error)
	ValidateGlobalUser(userName, password string) (string, error)
	CheckUserNameIsExist(domainID, userName string) (bool, error)
	CheckUserNameIsGlobalExist(userName string) (bool, error)
	CheckUserIsExistByID(userID string) (bool, error)

	ListUserProjects(domainID, userID string) ([]string, error)
	ListUserOtherDomains(userID string) ([]string, error)
}

// Writer use to write user information to store
type Writer interface {
	SaveVerifyCode(toEmail, phoneNumber string, code int) (*VerifyCode, error)
	GetVerifyCodeByMail(toEmail string, code int) (*VerifyCode, error)
	GetVerifyCodeByPhone(phoneNumber string, code int) (*VerifyCode, error)
	RevolkVerifyCode(id int64) error

	SaveUserOtherDomain(userID, otherDomainID string) error
	DeleteUserOtherDomain(userID, otherDomainID string) error
	CreateUser(domainID, userName, password string, enabled bool, userExpires, passExpires int) (*User, error)
	DeleteUser(domainID, userID string) error

	SetUserPassword(userID, oldPass, newPass string) error
	SetDefaultProject(domainID, userID, projectID string) error
	AddProjectsToUser(domainID, userID string, projectIDs ...string) error
	RemoveProjectsFromUser(domainID, userID string, projectIDs ...string) error
	BindRole(domainID, userID, roleName string) error
	UnBindRole(domainID, userID, roleName string) error

	SaveInvitationsRecord(inviterID string, invitedRoles, accessProjects []string) (*InvitationCode, error)
	ListInvitationRecord(inviterID string) ([]*InvitationCode, error)
	GetInvitationRecord(inviterID, code string) (*InvitationCode, error)
	DeleteInvitationRecord(id int64) error
	UpdateInvitationsRecord(ir *InvitationCode) error
}
