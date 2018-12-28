package mysql

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/satori/go.uuid"
)

func (s *store) CreateUser(u *user.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := s.CheckUserNameIsExist(u.DomainID, u.Account); err != nil {
		return err
	}

	// 开启存储事物
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("start create user transaction error, %s", err)
	}

	// 保存用户质料 Prepare SQL
	userPre, err := tx.Prepare(s.unprepared[SaveUser])
	defer userPre.Close()
	if err != nil {
		return exception.NewInternalServerError("prepare insert user stmt error, user: %s, %s", u.Account, err)
	}

	u.ID = uuid.NewV4().String()
	u.CreateAt = time.Now().Unix()
	// 存入
	if _, err = userPre.Exec(u.ID, u.DepartmentID, u.Account, u.Mobile, u.Email, u.Phone, u.Address, u.RealName, u.NickName,
		u.Gender, u.Avatar, u.Locked, u.DomainID, u.CreateAt, u.ExpiresActiveDays, u.DefaultProjectID); err != nil {
		return exception.NewInternalServerError("insert user exec sql err, %s", err)
	}

	// 保存密码
	if u.Password != nil {
		// Prepare SQL
		passPre, err := tx.Prepare(s.unprepared[SavePass])
		defer passPre.Close()
		if err != nil {
			return exception.NewInternalServerError("prepare insert user password error, user: %s, %s", u.Account, err)
		}

		hashPW := s.hmacHash(u.Password.Password)
		pass := user.Password{CreateAt: time.Now().Unix(), ExpireAt: u.Password.ExpireAt, Password: hashPW, UserID: u.ID}
		// 存入
		if _, err = passPre.Exec(pass.Password, pass.ExpireAt, pass.CreateAt, pass.UserID); err != nil {
			return exception.NewInternalServerError("insert password exec sql err, %s", err)
		}
	}

	// 提交事物
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return exception.NewInternalServerError("insert user transaction rollback error, %s", err)
		}
		return exception.NewInternalServerError("insert user transaction commit error, but rollback success, %s", err)
	}

	return nil
}

func (s *store) CheckUserNameIsExist(domainID, account string) error {
	rows, err := s.stmts[CheckUserExistByName].Query(account, domainID)
	if err != nil {
		return exception.NewInternalServerError("query user name exist error, %s", err)
	}
	defer rows.Close()

	userNames := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return exception.NewHasExist("scan user name exist record error, %s", err)
		}
		userNames = append(userNames, name)
	}
	if len(userNames) == 0 {
		return nil
	}

	return nil
}

func (s *store) CheckUserIsExistByID(userID string) (bool, error) {
	var uid string
	err := s.stmts[CheckUserExistByID].QueryRow(userID).Scan(&uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, exception.NewInternalServerError("check user exist by id error, %s", err)
	}

	return true, nil
}

func (s *store) hmacHash(msg string) string {
	mac := hmac.New(sha256.New, []byte(s.key))
	io.WriteString(mac, msg)

	return fmt.Sprintf("%x", mac.Sum(nil))
}
