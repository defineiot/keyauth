package mysql

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/satori/go.uuid"

	"openauth/api/exception"
	"openauth/storage/domain"
	"openauth/storage/project"
	"openauth/storage/user"
)

var (
	deletePrepare *sql.Stmt
)

// NewUserManager use to new an use service with mysql
func NewUserManager(db *sql.DB, key string, dm domain.Manager, pm project.Manager) (user.Manager, error) {
	return &manager{db: db, key: key, dm: dm, pm: pm}, nil
}

type manager struct {
	db  *sql.DB
	key string
	dm  domain.Manager
	pm  project.Manager
}

func (m *manager) CreateUser(domainID, name, password string, enabled bool, userExpires, passExpires int) (*user.User, error) {
	// check domain exist
	ok, err := m.dm.CheckDomainIsExistByID(domainID)
	if err != nil {
		return nil, exception.NewInternalServerError("check domain exist error, ", err)
	}
	if !ok {
		return nil, exception.NewBadRequest("domain %s not exist", domainID)
	}

	// check the domain user is exist

	tx, err := m.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("start create user transaction error, %s", err)
	}

	// insert password
	passPre, err := tx.Prepare("INSERT INTO `password` (password, expires_at, create_at) VALUES (?,?,?)")
	if err != nil {
		return nil, fmt.Errorf("prepare insert user password error, user: %s, %s", name, err)
	}

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

	// insert user
	userPre, err := tx.Prepare("INSERT INTO `user` (id, name, enabled, domain_id, create_at, password_id, expires_active_days) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		return nil, fmt.Errorf("prepare insert user stmt error, user: %s, %s", name, err)
	}

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

	// commit transaction
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, fmt.Errorf("insert user transaction rollback error, %s", err)
		}
		return nil, fmt.Errorf("insert user transaction commit error, but rollback success, %s", err)
	}

	return &user, nil
}

func (m *manager) AddProjectsToUser(userID string, projectIDs ...string) error {
	mappPre, err := m.db.Prepare("INSERT INTO `mapping` (user_id, project_id) VALUES (?,?)")
	if err != nil {
		return fmt.Errorf("prepare add projects to user mapping stmt error, user: %s, project: %s, %s", userID, projectIDs, err)
	}

	for _, projectID := range projectIDs {
		_, err = mappPre.Exec(userID, projectID)
		if err != nil {
			mappPre.Close()
			return fmt.Errorf("insert add projects to user mapping exec sql err, %s", err)
		}
	}
	mappPre.Close()

	return nil
}

func (m *manager) RemoveProjectsFromUser(userID string, projectIDs ...string) error {
	mappPre, err := m.db.Prepare("DELETE FROM `mapping` WHERE user_id = ? AND project_id = ?")
	if err != nil {
		return fmt.Errorf("prepare remove projects to user mapping stmt error, user: %s, project: %s, %s", userID, projectIDs, err)
	}

	for _, projectID := range projectIDs {
		_, err = mappPre.Exec(userID, projectID)
		if err != nil {
			mappPre.Close()
			return fmt.Errorf("insert remove projects to user mapping exec sql err, %s", err)
		}
	}
	mappPre.Close()

	return nil
}

func (m *manager) GetUserByName(domainID, userName, userPassword string) (*user.User, error) {
	return nil, nil
}

func (m *manager) GetUserByID(userID, userPassword string) (*user.User, error) {
	// get user by id

	// get user's project by mapping

	// get user's emails

	// get user's phones

	// get user's password
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
