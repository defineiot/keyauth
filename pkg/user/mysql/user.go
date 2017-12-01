package mysql

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"time"

	"openauth/pkg/user"

	uuid "github.com/satori/go.uuid"
)

var (
	deletePrepare *sql.Stmt
)

// NewUserManager use to new an use service with mysql
func NewUserManager(db *sql.DB, key string) (user.Manager, error) {
	return &manager{db: db, key: key}, nil
}

type manager struct {
	db  *sql.DB
	key string
}

func (m *manager) CreateUser(domainID, projectID, name, password string, enabled bool, userExpires, passExpires int) (*user.User, error) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("start create user transaction error, %s", err)
	}

	passPre, err := tx.Prepare("INSERT INTO `password` (password, expires_at, create_at) VALUES (?,?,?)")
	if err != nil {
		return nil, fmt.Errorf("prepare insert user password error, user: %s, %s", name, err)
	}

	// insert password
	now := time.Now()
	delta, err := time.ParseDuration(fmt.Sprintf("%dh", passExpires))
	if err != nil {
		passPre.Close()
		return nil, fmt.Errorf("parse password time delta error, expires: %d, %s", passExpires, err)
	}
	exp := now.Add(delta)
	hashPW := m.hmacHash(password)
	pass := user.Password{CreateAt: now.Unix(), ExpireAt: exp.Unix(), Password: hashPW}
	ret, err := passPre.Exec(pass.Password, pass.ExpireAt, pass.CreateAt)
	if err != nil {
		passPre.Close()
		return nil, fmt.Errorf("insert password exec sql err, %s", err)
	}
	id, err := ret.LastInsertId()
	if err != nil {
		passPre.Close()
		return nil, fmt.Errorf("get the last insert id error, %s", err)
	}
	pass.ID = id
	passPre.Close()

	userPre, err := tx.Prepare("INSERT INTO `user` (id, name, enabled, domain_id, create_at, password_id, expires_active_days) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		return nil, fmt.Errorf("prepare insert user stmt error, user: %s, %s", name, err)
	}

	// insert user
	deltaU, err := time.ParseDuration(fmt.Sprintf("%dh", userExpires))
	if err != nil {
		userPre.Close()
		return nil, fmt.Errorf("parse user time delta error, expires: %d, %s", userExpires, err)
	}
	expU := now.Add(deltaU)
	user := user.User{ID: uuid.NewV4().String(), Name: name, Enabled: enabled, DomainID: domainID, CreateAt: time.Now().Unix(), ExpireActiveDays: expU.Unix()}
	_, err = userPre.Exec(user.ID, user.Name, user.Enabled, user.DomainID, user.CreateAt, pass.ID, user.ExpireActiveDays)
	if err != nil {
		userPre.Close()
		return nil, fmt.Errorf("insert user exec sql err, %s", err)
	}
	userPre.Close()

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, fmt.Errorf("insert user transaction rollback error, %s", err)
		}
		return nil, fmt.Errorf("insert user transaction commit error, but rollback success, %s", err)
	}

	return &user, nil
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

func (m *manager) hmacHash(msg string) string {
	mac := hmac.New(sha256.New, []byte(m.key))
	io.WriteString(mac, msg)

	return fmt.Sprintf("%x", mac.Sum(nil))
}
