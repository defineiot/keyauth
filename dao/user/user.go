package user

import (
	"github.com/defineiot/keyauth/dao/models"
)

// Store is user service
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader use to read user information form store
type Reader interface {
	ListDomainUsers(domainID string) ([]*models.User, error)
	ListDepartmentUsers(departmentID string) ([]*models.User, error)
	ListProjectUsers(projectID string) ([]*models.User, error)
	GetUser(index models.FoundUserIndex, value string) (*models.User, error)

	// ListUserRoles(domainID, userID string) ([]string, error)
	// ValidateUser(domainID, userName, password string) (string, error)
	// ValidateGlobalUser(userName, password string) (string, error)
	// CheckUserNameIsExist(domainID, userName string) (bool, error)
	// CheckUserNameIsGlobalExist(userName string) (bool, error)
	// CheckUserIsExistByID(userID string) (bool, error)
	// ListUserProjects(domainID, userID string) ([]string, error)
	// ListUserOtherDomains(userID string) ([]string, error)
	CheckUserIsExistByID(userID string) (bool, error)
	CheckUserNameIsExist(domainID, account string) error
	// ValidateUser(domainID, userName, password string) (string, error)
}

// Writer use to write user information to store
type Writer interface {
	CreateUser(u *models.User) error
	DeleteUser(domainID, userID string) error
	BindRole(domainID, userID, roleID string) error
	UnBindRole(domainID, userID, roleID string) error
	// SaveUserOtherDomain(userID, otherDomainID string) error
	// DeleteUserOtherDomain(userID, otherDomainID string) error

	// SetUserPassword(userID, oldPass, newPass string) error
	// SetDefaultProject(domainID, userID, projectID string) error
	// AddProjectsToUser(domainID, userID string, projectIDs ...string) error
	// RemoveProjectsFromUser(domainID, userID string, projectIDs ...string) error

	// SaveInvitationsRecord(inviterID string, invitedRoles, accessProjects []string) (*Invitation, error)
	// ListInvitationRecord(inviterID string) ([]*Invitation, error)
	// GetInvitationRecord(inviterID, code string) (*Invitation, error)
	// DeleteInvitationRecord(id int64) error
	// UpdateInvitationsRecord(ir *Invitation) error
}
