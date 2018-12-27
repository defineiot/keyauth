package mysql

import (
	"fmt"
	"time"

	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/satori/go.uuid"
)

func (s *store) CreateUser(u *user.User) (*user.User, error) {
	// check user exist
	ok, err := s.CheckUserNameIsExist(domainID, name)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, exception.NewBadRequest("user %s is exist", name)
	}

	// start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("start create user transaction error, %s", err)
	}

	now := time.Now()

	// insert user
	userPre, err := tx.Prepare("INSERT INTO `users` (id, name, enabled, domain_id, create_at, expires_active_days) VALUES (?,?,?,?,?,?)")
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
	passPre, err := tx.Prepare("INSERT INTO `passwords` (password, expires_at, create_at, user_id) VALUES (?,?,?,?)")
	if err != nil {
		return nil, exception.NewInternalServerError("prepare insert user password error, user: %s, %s", name, err)
	}

	delta, err := time.ParseDuration(fmt.Sprintf("%dh", passExpires))
	if err != nil {
		passPre.Close()
		return nil, exception.NewBadRequest("parse password time delta error, expires: %d, %s", passExpires, err)
	}
	exp := now.Add(delta)
	hashPW := s.hmacHash(password)
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

func (s *store) CheckUserNameIsExist(domainID, userName string) error {
	rows, err := s.stmts[CheckUserExistByName].Query(userName, domainID)
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
