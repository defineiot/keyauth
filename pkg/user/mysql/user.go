package mysql

import (
	"database/sql"

	"openauth/pkg/user"
)

// NewUserManager use to new an use service with mysql
func NewUserManager(db *sql.DB) (user.Manager, error) {
	return &manager{}, nil
}

type manager struct {
	db *sql.DB
}

func (m *manager) CreateUser(projectID, userName, password string) (*user.User, error) {
	return nil, nil
}

func (m *manager) GetUserByName(domainID, userName, userPassword string) (*user.User, error) {
	return nil, nil
}

func (m *manager) GetUserByID(userID, userPassword string) (*user.User, error) {
	return nil, nil
}

func (m *manager) GetUser(cert user.Credential) (*user.User, error) {
	return nil, nil
}

func (m *manager) DeleteUser(cert user.Credential) error {
	return nil
}

func (m *manager) AddPhone(cert user.Credential, number, phoneType, description string) error {
	return nil
}

func (m *manager) RemovePhone(cert user.Credential, number string) error {
	return nil
}

func (m *manager) QueryPhone(cert user.Credential) (*[]user.Phone, error) {
	return nil, nil
}

func (m *manager) AddEmail(cert user.Credential, address, description string, primary bool) error {
	return nil
}

func (m *manager) RemoveEmail(cert user.Credential, address string) error {
	return nil
}

func (m *manager) QueryEmail(cert user.Credential) (*[]user.Email, error) {
	return nil, nil
}

func (m *manager) AddRoleToUser(cert user.Credential, roleID string) error {
	return nil
}

func (m *manager) RemoveRoleFromUser(cert user.Credential, roleID string) error {
	return nil
}

func (m *manager) QueryRole(cert user.Credential) ([]string, error) {
	return nil, nil
}

func (m *manager) HasFeatures(cert user.Credential, features ...string) (bool, error) {
	return false, nil
}

func (m *manager) SetDefaultProject(cert user.Credential, projectID string) error {
	return nil
}

func (m *manager) AddUserToProject(cert user.Credential, projectID string) error {
	return nil
}

func (m *manager) RemoveUserFromProject(cert user.Credential, projectID string) error {
	return nil
}

func (m *manager) ListUserProject(cert user.Credential) ([]string, error) {
	return nil, nil
}

func (m *manager) GetDefaultProject(cert user.Credential) (string, error) {
	return "", nil
}

func (m *manager) IsSystemAdmin(cert user.Credential) (bool, error) {
	return false, nil
}

func (m *manager) IsCloudAdmin(cert user.Credential) (bool, error) {
	return false, nil
}
