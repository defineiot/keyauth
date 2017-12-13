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
	"openauth/api/logger"
	"openauth/pkg/user"
)

var (
	deletePrepare *sql.Stmt
)

// NewUserService use to new an use service with mysql
func NewUserService(db *sql.DB, key string, logger logger.OpenAuthLogger) user.Service {
	return &manager{db: db, key: key, logger: logger}
}

type manager struct {
	db     *sql.DB
	key    string
	logger logger.OpenAuthLogger
}

func (m *manager) CreateUser(domainID, name, password string, enabled bool, userExpires, passExpires int) (*user.User, error) {
	// check is user exist
	ok, err := m.CheckUserNameIsExist(domainID, name)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, exception.NewBadRequest("user %s is exist", name)
	}

	// start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("start create user transaction error, %s", err)
	}

	now := time.Now()

	// insert user
	userPre, err := tx.Prepare("INSERT INTO `user` (id, name, enabled, domain_id, create_at, expires_active_days) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return nil, exception.NewInternalServerError("prepare insert user stmt error, user: %s, %s", name, err)
	}

	deltaU, err := time.ParseDuration(fmt.Sprintf("%dh", userExpires))
	if err != nil {
		userPre.Close()
		return nil, exception.NewBadRequest("parse user time delta error, expires: %d, %s", userExpires, err)
	}
	expU := now.Add(deltaU)
	userI := user.User{ID: uuid.NewV4().String(), Name: name, Enabled: enabled, DomainID: domainID, CreateAt: time.Now().Unix(), ExpireActiveDays: expU.Unix()}
	_, err = userPre.Exec(userI.ID, userI.Name, userI.Enabled, userI.DomainID, userI.CreateAt, userI.ExpireActiveDays)
	if err != nil {
		userPre.Close()
		return nil, exception.NewInternalServerError("insert user exec sql err, %s", err)
	}
	userPre.Close()

	// insert password
	passPre, err := tx.Prepare("INSERT INTO `password` (password, expires_at, create_at, user_id) VALUES (?,?,?,?)")
	if err != nil {
		return nil, exception.NewInternalServerError("prepare insert user password error, user: %s, %s", name, err)
	}

	delta, err := time.ParseDuration(fmt.Sprintf("%dh", passExpires))
	if err != nil {
		passPre.Close()
		return nil, exception.NewBadRequest("parse password time delta error, expires: %d, %s", passExpires, err)
	}
	exp := now.Add(delta)
	hashPW := m.hmacHash(password)
	pass := user.Password{CreateAt: now.Unix(), ExpireAt: exp.Unix(), Password: hashPW, UserID: userI.ID}
	ret, err := passPre.Exec(pass.Password, pass.ExpireAt, pass.CreateAt, pass.UserID)
	if err != nil {
		passPre.Close()
		return nil, exception.NewInternalServerError("insert password exec sql err, %s", err)
	}
	id, err := ret.LastInsertId()
	if err != nil {
		passPre.Close()
		return nil, exception.NewInternalServerError("get the last insert id error, %s", err)
	}
	pass.ID = id
	passPre.Close()

	// commit transaction
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, exception.NewInternalServerError("insert user transaction rollback error, %s", err)
		}
		return nil, exception.NewInternalServerError("insert user transaction commit error, but rollback success, %s", err)
	}

	return &userI, nil
}

func (m *manager) ListUserProjects(userID string) ([]string, error) {
	ok, err := m.CheckUserIsExistByID(userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("user %s not exist", userID)
	}

	rows, err := m.db.Query("SELECT project_id FROM mapping WHERE user_id = ?", userID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user's project id error, %s", err)
	}
	defer rows.Close()

	projects := []string{}
	for rows.Next() {
		var projectID string
		if err := rows.Scan(&projectID); err != nil {
			return nil, exception.NewInternalServerError("scan user's project id error, %s", err)
		}
		projects = append(projects, projectID)
	}
	return projects, nil
}

func (m *manager) SetDefaultProject(userID, projectID string) error {
	projects, err := m.ListUserProjects(userID)
	if err != nil {
		return err
	}

	// check the project is user's project
	var ok bool
	for _, pid := range projects {
		if pid == projectID {
			ok = true
		}
	}
	if !ok {
		return exception.NewBadRequest("user %s hasn't project %s", userID, projectID)
	}

	mappPre, err := m.db.Prepare("UPDATE `user` SET default_project_id = ? WHERE id = ?")
	if err != nil {
		return exception.NewInternalServerError("prepare set user's default project error, user: %s, project: %s, %s", userID, projectID, err)
	}
	defer mappPre.Close()

	ret, err := mappPre.Exec(projectID, userID)
	if err != nil {
		return exception.NewInternalServerError("set user's default project exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get affect rows count error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("user %s not exist")
	}

	return nil
}

func (m *manager) AddProjectsToUser(userID string, projectIDs ...string) error {
	// check the user is not exist
	ok, err := m.CheckUserIsExistByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("user %s not exist", userID)
	}

	// insert
	mappPre, err := m.db.Prepare("INSERT INTO `mapping` (user_id, project_id) VALUES (?,?)")
	if err != nil {
		return fmt.Errorf("prepare add projects to user mapping stmt error, user: %s, project: %s, %s", userID, projectIDs, err)
	}
	defer mappPre.Close()

	for _, projectID := range projectIDs {
		//
		_, err = mappPre.Exec(userID, projectID)
		if err != nil {
			return fmt.Errorf("insert add projects to user mapping exec sql err, %s", err)
		}
	}

	return nil
}

func (m *manager) RemoveProjectsFromUser(userID string, projectIDs ...string) error {
	// check the user is not exist
	ok, err := m.CheckUserIsExistByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("user %s not exist", userID)
	}

	// insert
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

func (m *manager) GetUserByID(userID string) (*user.User, error) {
	// get user by id
	m.logger.Debug("get user by id...")
	userI := user.User{}
	err := m.db.QueryRow("SELECT id,name,enabled,last_active_time,domain_id,create_at,expires_active_days,default_project_id FROM user WHERE id = ?", userID).Scan(
		&userI.ID, &userI.Name, &userI.Enabled, &userI.LastActiveTime, &userI.DomainID, &userI.CreateAt, &userI.ExpireActiveDays, &userI.DefaultProjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("user %s not find", userID)
		}

		return nil, exception.NewInternalServerError("query single user error, %s", err)
	}

	// get user's emails
	emails, err := m.QueryEmail(userID)
	if err != nil {
		return nil, err
	}
	userI.Emails = emails

	// get user's phones
	phones, err := m.QueryPhone(userID)
	if err != nil {
		return nil, err
	}
	userI.Phones = phones

	// get user's password
	userI.Password = &user.Password{}

	m.logger.Debug("get user: ", userI)
	return &userI, nil
}

func (m *manager) QueryPhone(userID string) ([]*user.Phone, error) {
	rows, err := m.db.Query("SELECT id,numbers,'primary',description FROM phone WHERE user_id = ?", userID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user's phone error, %s", err)
	}
	defer rows.Close()

	phones := []*user.Phone{}
	for rows.Next() {
		phone := user.Phone{}
		if err := rows.Scan(&phone.ID, &phone.Number, &phone.Primary, &phone.Description); err != nil {
			return nil, exception.NewInternalServerError("scan user's phone record error, %s", err)
		}
		phones = append(phones, &phone)
	}

	return phones, nil
}

func (m *manager) QueryEmail(userID string) ([]*user.Email, error) {
	rows, err := m.db.Query("SELECT id,address,'primary',description FROM email WHERE user_id = ?", userID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user's email error, %s", err)
	}
	defer rows.Close()

	emails := []*user.Email{}
	for rows.Next() {
		email := user.Email{}
		if err := rows.Scan(&email.ID, &email.Address, &email.Primary, &email.Description); err != nil {
			return nil, exception.NewInternalServerError("scan user's email record error, %s", err)
		}
		emails = append(emails, &email)
	}

	return emails, nil
}

func (m *manager) ListUser(domainID string) ([]*user.User, error) {

	rows, err := m.db.Query("SELECT id,name,enabled,last_active_time,create_at,expires_active_days,default_project_id FROM user WHERE domain_id = ?", domainID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user list error, %s", err)
	}
	defer rows.Close()

	users := []*user.User{}
	for rows.Next() {
		u := user.User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.Enabled, &u.LastActiveTime, &u.CreateAt, &u.ExpireActiveDays, &u.DefaultProjectID); err != nil {
			return nil, exception.NewInternalServerError("scan project's user id error, %s", err)
		}
		// get user's emails
		emails, err := m.QueryEmail(u.ID)
		if err != nil {
			return nil, err
		}
		u.Emails = emails
		// get user's phone
		phones, err := m.QueryPhone(u.ID)
		if err != nil {
			return nil, err
		}
		u.Phones = phones
		// get user's password

		u.DomainID = domainID

		users = append(users, &u)
	}

	m.logger.Debug("get user: ", users)

	return users, nil

}

func (m *manager) DeleteUser(userID string) error {
	tx, err := m.db.Begin()
	if err != nil {
		return exception.NewInternalServerError("start delete user transaction error, %s", err)
	}

	// delete user
	deleteUserPre, err := tx.Prepare("DELETE FROM user WHERE id = ?")
	if err != nil {
		deleteUserPre.Close()
		return exception.NewInternalServerError("prepare delete user stmt error, %s", err)
	}

	ret, err := deleteUserPre.Exec(userID)
	if err != nil {
		deleteUserPre.Close()
		return exception.NewInternalServerError("delete user exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		deleteUserPre.Close()
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		deleteUserPre.Close()
		return exception.NewBadRequest("user %s not exist", userID)
	}
	deleteUserPre.Close()

	// delete user's password
	deletePassPre, err := tx.Prepare("DELETE FROM password WHERE user_id = ?")
	if err != nil {
		deletePassPre.Close()
		return exception.NewInternalServerError("prepare delete user pass stmt error, %s", err)
	}

	ret, err = deletePassPre.Exec(userID)
	if err != nil {
		deletePassPre.Close()
		return exception.NewInternalServerError("delete user pass exec sql error, %s", err)
	}
	deletePassPre.Close()

	// delete user's email
	deleteEmailPre, err := tx.Prepare("DELETE FROM email WHERE user_id = ?")
	if err != nil {
		deleteEmailPre.Close()
		return exception.NewInternalServerError("prepare delete user email stmt error, %s", err)
	}

	ret, err = deleteEmailPre.Exec(userID)
	if err != nil {
		deleteEmailPre.Close()
		return exception.NewInternalServerError("delete user pass exec sql error, %s", err)
	}
	deleteEmailPre.Close()

	// delete user's phone
	deletePhonePre, err := tx.Prepare("DELETE FROM phone WHERE user_id = ?")
	if err != nil {
		deletePhonePre.Close()
		return exception.NewInternalServerError("prepare delete user phone stmt error, %s", err)
	}

	ret, err = deletePhonePre.Exec(userID)
	if err != nil {
		deletePhonePre.Close()
		return exception.NewInternalServerError("delete user phone exec sql error, %s", err)
	}
	deletePhonePre.Close()

	// commit transaction
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return exception.NewInternalServerError("delete user transaction rollback error, %s", err)
		}
		return exception.NewInternalServerError("delete user transaction commit error, but rollback success, %s", err)
	}

	return nil
}

func (m *manager) GetUserByName(domainID, userName, userPassword string) (*user.User, error) {
	return nil, nil
}

func (m *manager) GetUser(cert user.Credential) (*user.User, error) {
	return nil, nil
}

func (m *manager) AddPhone(cert user.Credential, number, phoneType, description string) error {
	return nil
}

func (m *manager) RemovePhone(cert user.Credential, number string) error {
	return nil
}

func (m *manager) AddEmail(cert user.Credential, address, description string, primary bool) error {
	return nil
}

func (m *manager) RemoveEmail(cert user.Credential, address string) error {
	return nil
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

func (m *manager) AddUserToProject(cert user.Credential, projectID string) error {
	return nil
}

func (m *manager) RemoveUserFromProject(cert user.Credential, projectID string) error {
	return nil
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

func (m *manager) CheckUserNameIsExist(domainID, userName string) (bool, error) {
	rows, err := m.db.Query("SELECT name FROM user WHERE name = ? AND domain_id = ?", userName, domainID)
	if err != nil {
		return false, exception.NewInternalServerError("query user name exist error, %s", err)
	}
	defer rows.Close()

	userNames := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return false, exception.NewInternalServerError("scan user name exist record error, %s", err)
		}
		userNames = append(userNames, name)
	}
	if len(userNames) == 0 {
		return false, nil
	}

	return true, nil
}

func (m *manager) CheckUserIsExistByID(userID string) (bool, error) {
	var uid string
	err := m.db.QueryRow("SELECT id FROM user WHERE id = ?", userID).Scan(&uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, exception.NewInternalServerError("check user exist by id error, %s", err)
	}

	return true, nil
}

func (m *manager) hmacHash(msg string) string {
	mac := hmac.New(sha256.New, []byte(m.key))
	io.WriteString(mac, msg)

	return fmt.Sprintf("%x", mac.Sum(nil))
}
