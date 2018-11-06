package mysql

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/satori/go.uuid"

	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/internal/exception"
)

func (s *store) SaveVerifyCode(toEmail, phoneNumber string, code int) (*user.VerifyCode, error) {
	after, err := time.ParseDuration("1h")
	if err != nil {
		return nil, err
	}
	expire := time.Now().Add(after)

	vc := &user.VerifyCode{EmailAddress: toEmail, PhoneNumber: phoneNumber, Code: code, CreateAt: time.Now().Unix(), ExpireAt: expire.Unix()}

	_, err = s.stmts[SaveVerifyCode].Exec(vc.EmailAddress, vc.PhoneNumber, vc.Code, vc.CreateAt, vc.ExpireAt)
	if err != nil {
		return nil, exception.NewInternalServerError("insert verify code exec sql err, %s", err)
	}

	return vc, nil
}

func (s *store) GetVerifyCodeByMail(toEmail string, code int) (*user.VerifyCode, error) {
	vc := user.VerifyCode{}
	err := s.stmts[FindVerifyCodeByMail].QueryRow(code, toEmail).Scan(
		&vc.ID, &vc.EmailAddress, &vc.PhoneNumber, &vc.Code, &vc.CreateAt, &vc.ExpireAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("verify code %d not find", code)
		}

		return nil, exception.NewInternalServerError("query single verify code error, %s", err)
	}

	return &vc, nil
}

func (s *store) GetVerifyCodeByPhone(phoneNumber string, code int) (*user.VerifyCode, error) {
	vc := user.VerifyCode{}
	err := s.stmts[FindVerifyCodeByPhone].QueryRow(code, phoneNumber).Scan(
		&vc.ID, &vc.EmailAddress, &vc.PhoneNumber, &vc.Code, &vc.CreateAt, &vc.ExpireAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("verify code %s not find", code)
		}

		return nil, exception.NewInternalServerError("query single verify code error, %s", err)
	}

	return &vc, nil
}

func (s *store) RevolkVerifyCode(id int64) error {
	ret, err := s.stmts[DeleteVerifyCodeByID].Exec(id)
	if err != nil {
		return exception.NewInternalServerError("delete domain exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("domian %s not exist", id)
	}

	return nil
}

func (s *store) SaveInvitationsRecord(inviterID string, invitedRoles, accessProjects []string) (*user.InvitationCode, error) {
	after, err := time.ParseDuration("1h")
	if err != nil {
		return nil, err
	}
	expire := time.Now().Add(after)

	code := hex.EncodeToString(uuid.NewV4().Bytes())

	ir := &user.InvitationCode{Code: code, InviterID: inviterID, InvitedUserRoleNames: invitedRoles, AccessProjects: accessProjects, InvitedTime: time.Now().Unix(), ExpireTime: expire.Unix()}
	_, err = s.stmts[SaveInvitationsRecord].Exec(ir.InviterID, strings.Join(invitedRoles, ","), strings.Join(accessProjects, ","), ir.InvitedTime, ir.ExpireTime, code)
	if err != nil {
		return nil, exception.NewInternalServerError("insert verify code exec sql err, %s", err)
	}

	return ir, nil
}

func (s *store) ListInvitationRecord(inviterID string) ([]*user.InvitationCode, error) {
	ok, err := s.CheckUserIsExistByID(inviterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("inviter user %s not exist", inviterID)
	}

	rows, err := s.stmts[FindUserAllInvitationsRecords].Query(inviterID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user's invitation records error, %s", err)
	}
	defer rows.Close()

	irs := []*user.InvitationCode{}
	for rows.Next() {
		ir := new(user.InvitationCode)
		if err := rows.Scan(ir.ID, ir.InviterID, ir.InvitedUserID, ir.InvitedUserDomainID, ir.InvitedUserRoleNames, ir.InvitedTime, ir.AcceptTime, ir.ExpireTime, ir.Code, ir.AccessProjects); err != nil {
			return nil, exception.NewInternalServerError("scan user's project id error, %s", err)
		}
		irs = append(irs, ir)
	}
	return irs, nil
}

func (s *store) GetInvitationRecord(inviterID, code string) (*user.InvitationCode, error) {
	s.log.Debugf("Get Invitation Record SQL: %s Params: inviter: %s, code: %s", s.unprepared[FindOneInvitationRecord], inviterID, code)
	ir := new(user.InvitationCode)
	roles := ""
	projects := ""
	err := s.stmts[FindOneInvitationRecord].QueryRow(inviterID, code).Scan(&ir.ID, &ir.InviterID, &ir.InvitedUserID, &ir.InvitedUserDomainID, &roles, &ir.InvitedTime, &ir.AcceptTime, &ir.ExpireTime, &ir.Code, &projects)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("registry code %s not find", code)
		}

		return nil, exception.NewInternalServerError("query single registry code error, %s", err)
	}

	if roles != "" {
		ir.InvitedUserRoleNames = strings.Split(roles, ",")
	}
	if projects != "" {
		ir.AccessProjects = strings.Split(projects, ",")
	}

	return ir, nil
}

func (s *store) DeleteInvitationRecord(id int64) error {
	ret, err := s.stmts[DeleteInvitationRecord].Exec(id)
	if err != nil {
		return exception.NewInternalServerError("delete invitation record exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("invitation recode %s not exist", id)
	}

	return nil
}

func (s *store) UpdateInvitationsRecord(ir *user.InvitationCode) error {
	ir.AcceptTime = time.Now().Unix()
	_, err := s.stmts[UpdateInvitationsRecord].Exec(ir.InvitedUserID, ir.InvitedUserDomainID, ir.AcceptTime, ir.InviterID, ir.Code)
	if err != nil {
		return exception.NewInternalServerError("insert verify code exec sql err, %s", err)
	}

	return nil
}

func (s *store) SaveUserOtherDomain(userID, otherDomainID string) error {
	_, err := s.stmts[SaveUserOtherDomain].Exec(userID, otherDomainID, time.Now().Unix())
	if err != nil {
		return exception.NewInternalServerError("insert user other domain exec sql err, %s", err)
	}

	return nil
}

func (s *store) ListUserOtherDomains(userID string) ([]string, error) {
	ok, err := s.CheckUserIsExistByID(userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("user %s not exist", userID)
	}

	rows, err := s.stmts[FindUserOtherDomain].Query(userID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user's other domain error, %s", err)
	}
	defer rows.Close()

	domains := []string{}
	for rows.Next() {
		var did string
		if err := rows.Scan(&did); err != nil {
			return nil, exception.NewInternalServerError("scan user's other domain id error, %s", err)
		}
		domains = append(domains, did)
	}
	return domains, nil
}

func (s *store) DeleteUserOtherDomain(userID, otherDomainID string) error {
	ret, err := s.stmts[DeleteUserOtherDomain].Exec(userID, otherDomainID)
	if err != nil {
		return exception.NewInternalServerError("delete user other domain exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("user %s domain %s not exist", userID, otherDomainID)
	}

	return nil
}

func (s *store) CreateUser(domainID, name, password string, enabled bool, userExpires, passExpires int) (*user.User, error) {
	// check is user exist
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

func (s *store) ListUserProjects(domainID, userID string) ([]string, error) {
	ok, err := s.CheckUserIsExistByID(userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, exception.NewBadRequest("user %s not exist", userID)
	}

	rows, err := s.stmts[FindUserProjects].Query(userID, domainID)
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

func (s *store) SetUserPassword(userID, oldPass, newPass string) error {
	ok, err := s.CheckUserIsExistByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("user %s not exist", userID)
	}

	pass, err := s.QueryPassword(userID)
	if err != nil {
		return err
	}

	if s.hmacHash(oldPass) != pass.Password {
		return exception.NewBadRequest("your old password is incorrect")
	}

	ret, err := s.stmts[SetUserPassword].Exec(s.hmacHash(newPass), userID)
	if err != nil {
		return exception.NewInternalServerError("set user's  password exec sql error, %s", err)
	}

	ra, err := ret.RowsAffected()
	if err != nil {
		s.log.Error(err)
	}
	if ra == 0 {
		return exception.NewInternalServerError("not row updated here.")
	}

	return nil
}

func (s *store) SetDefaultProject(domainID, userID, projectID string) error {
	ok, err := s.CheckUserIsExistByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("user %s not exist", userID)
	}

	projects, err := s.ListUserProjects(domainID, userID)
	if err != nil {
		return err
	}

	// check the project is user's project
	for _, pid := range projects {
		if pid == projectID {
			ok = true
		}
	}
	if !ok {
		return exception.NewBadRequest("user %s hasn't project %s", userID, projectID)
	}

	_, err = s.stmts[SetUserDefaultProject].Exec(projectID, userID)
	if err != nil {
		return exception.NewInternalServerError("set user's default project exec sql error, %s", err)
	}

	return nil
}

func (s *store) AddProjectsToUser(domainID, userID string, projectIDs ...string) error {
	// check the user is not exist
	ok, err := s.CheckUserIsExistByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("user %s not exist", userID)
	}

	// check projects is owned by this user
	pids, err := s.ListUserProjects(domainID, userID)
	if err != nil {
		return err
	}
	existPIDs := []string{}
	for _, pid := range pids {
		for _, inpid := range projectIDs {
			if inpid == pid {
				existPIDs = append(existPIDs, inpid)
			}
		}
	}
	if len(existPIDs) != 0 {
		return exception.NewBadRequest("project %s is in this project", existPIDs)
	}

	for _, projectID := range projectIDs {
		//
		_, err = s.stmts[AddProjectToUser].Exec(userID, projectID)
		if err != nil {
			return fmt.Errorf("insert add projects to user mapping exec sql err, %s", err)
		}
	}

	return nil
}

func (s *store) RemoveProjectsFromUser(domainID, userID string, projectIDs ...string) error {
	// check the user is not exist
	ok, err := s.CheckUserIsExistByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("user %s not exist", userID)
	}

	// check projects is owned by this user
	pids, err := s.ListUserProjects(domainID, userID)
	if err != nil {
		return err
	}
	nExistPIDs := []string{}
	for _, inpid := range projectIDs {
		var ok bool
		for _, pid := range pids {
			if pid == inpid {
				ok = true
			}
		}
		if !ok {
			nExistPIDs = append(nExistPIDs, inpid)
		}
	}
	if len(nExistPIDs) != 0 {
		return exception.NewBadRequest("project %s isn't in this project", nExistPIDs)
	}

	for _, projectID := range projectIDs {
		_, err = s.stmts[RemoveProjectsFromUser].Exec(userID, projectID)
		if err != nil {
			return fmt.Errorf("insert remove projects to user mapping exec sql err, %s", err)
		}
	}

	return nil
}

func (s *store) GetUserByID(userID string) (*user.User, error) {
	// get user by id
	userI := user.User{}
	err := s.stmts[FindUserByID].QueryRow(userID).Scan(
		&userI.ID, &userI.Name, &userI.Enabled, &userI.LastActiveTime, &userI.DomainID, &userI.CreateAt, &userI.ExpireActiveDays, &userI.DefaultProjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("user %s not find", userID)
		}

		return nil, exception.NewInternalServerError("query single user error, %s", err)
	}

	// get user's emails
	emails, err := s.QueryEmail(userID)
	if err != nil {
		return nil, err
	}
	userI.Emails = emails

	// get user's phones
	phones, err := s.QueryPhone(userID)
	if err != nil {
		return nil, err
	}
	userI.Phones = phones

	// get user's password
	userI.Password = &user.Password{}

	return &userI, nil
}

func (s *store) QueryPhone(userID string) ([]*user.Phone, error) {
	rows, err := s.stmts[FindUserPhones].Query(userID)
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

func (s *store) QueryEmail(userID string) ([]*user.Email, error) {
	rows, err := s.stmts[FindUserEmails].Query(userID)
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

func (s *store) QueryPassword(userID string) (*user.Password, error) {
	pass := user.Password{}
	err := s.stmts[FindUserPassword].QueryRow(userID).Scan(
		&pass.Password, &pass.ExpireAt, &pass.CreateAt, &pass.UpdateAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("password %s not find", userID)
		}

		return nil, exception.NewInternalServerError("query single password error, %s", err)
	}

	return &pass, nil
}

func (s *store) ListUser(domainID string) ([]*user.User, error) {
	rows, err := s.stmts[FindAllUsers].Query(domainID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user list error, %s", err)
	}
	defer rows.Close()

	return s.getUserFromRows(rows)
}

func (s *store) getUserFromRows(rows *sql.Rows) ([]*user.User, error) {
	users := []*user.User{}
	for rows.Next() {
		u := user.User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.Enabled, &u.LastActiveTime, &u.CreateAt, &u.ExpireActiveDays, &u.DefaultProjectID, &u.DomainID); err != nil {
			return nil, exception.NewInternalServerError("scan project's user id error, %s", err)
		}

		// get user's emails
		emails, err := s.QueryEmail(u.ID)
		if err != nil {
			return nil, err
		}
		u.Emails = emails
		// get user's phone
		phones, err := s.QueryPhone(u.ID)
		if err != nil {
			return nil, err
		}
		u.Phones = phones
		// get user's password
		pass, err := s.QueryPassword(u.ID)
		if err != nil {
			return nil, err
		}
		u.Password = pass

		users = append(users, &u)
	}
	return users, nil
}

func (s *store) DeleteUser(domainID, userID string) error {
	// check user has project
	pids, err := s.ListUserProjects(domainID, userID)
	if err != nil {
		return err
	}
	if len(pids) != 0 {
		return exception.NewBadRequest("this user: %s has proejcts: %v", userID, pids)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return exception.NewInternalServerError("start delete user transaction error, %s", err)
	}

	// delete user
	deleteUserPre, err := tx.Prepare("DELETE FROM users WHERE id = ?")
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
	deletePassPre, err := tx.Prepare("DELETE FROM passwords WHERE user_id = ?")
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
	deleteEmailPre, err := tx.Prepare("DELETE FROM emails WHERE user_id = ?")
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
	deletePhonePre, err := tx.Prepare("DELETE FROM phones WHERE user_id = ?")
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

func (s *store) GetUserByName(domainID, userName string) (*user.User, error) {
	userI := user.User{}
	err := s.stmts[FindUserByName].QueryRow(userName, domainID).Scan(
		&userI.ID, &userI.Name, &userI.Enabled, &userI.LastActiveTime, &userI.DomainID, &userI.CreateAt, &userI.ExpireActiveDays, &userI.DefaultProjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("user %s not find", userName)
		}

		return nil, exception.NewInternalServerError("query single user error, %s", err)
	}

	// get user's emails
	emails, err := s.QueryEmail(userI.ID)
	if err != nil {
		return nil, err
	}
	userI.Emails = emails

	// get user's phones
	phones, err := s.QueryPhone(userI.ID)
	if err != nil {
		return nil, err
	}
	userI.Phones = phones

	// get user's password
	userI.Password = &user.Password{}

	return &userI, nil
}

func (s *store) ValidateUser(domainID, userName, password string) (string, error) {
	var uid string
	err := s.stmts[FindUserIDByName].QueryRow(userName, domainID).Scan(&uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", exception.NewNotFound("user %s not find", userName)
		}

		return "", exception.NewInternalServerError("query single user error, %s", err)
	}

	pass, err := s.QueryPassword(uid)
	if err != nil {
		return "", err
	}

	if s.hmacHash(password) == pass.Password {
		return uid, nil
	}

	return "", nil
}

func (s *store) ValidateGlobalUser(userName, password string) (string, error) {
	var uid string
	err := s.stmts[FindGlobalUserIDByName].QueryRow(userName).Scan(&uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", exception.NewNotFound("user %s not find", userName)
		}

		return "", exception.NewInternalServerError("query single user error, %s", err)
	}

	pass, err := s.QueryPassword(uid)
	if err != nil {
		return "", err
	}

	if s.hmacHash(password) == pass.Password {
		return uid, nil
	}

	return "", nil
}

func (s *store) CheckUserNameIsExist(domainID, userName string) (bool, error) {
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

func (s *store) CheckUserNameIsGlobalExist(userName string) (bool, error) {
	var uid string
	err := s.stmts[CheckUserGlobalExistByName].QueryRow(userName).Scan(&uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, exception.NewInternalServerError("check user exist by global name error, %s", err)
	}

	return true, nil
}

func (s *store) BindRole(domainID, userID, roleName string) error {
	ok, err := s.CheckUserIsExistByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("user %s not exist", userID)
	}

	ok, err = s.checkUserRoleExist(domainID, userID, roleName)
	if err != nil {
		return err
	}
	if ok {
		return exception.NewBadRequest("user %s has bind the role: %s", userID, roleName)
	}

	_, err = s.stmts[BindRole].Exec(domainID, userID, roleName)
	if err != nil {
		return exception.NewInternalServerError("insert role user mapping exec sql err, %s", err)
	}

	return nil
}

func (s *store) UnBindRole(domainID, userID, roleName string) error {
	ok, err := s.CheckUserIsExistByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("user %s exist", userID)
	}

	ret, err := s.stmts[UnbindRole].Exec(domainID, userID, roleName)
	if err != nil {
		return exception.NewInternalServerError("delete role user mapping exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete role user mapping affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("the role: %s is not bind to user: %s", roleName, userID)
	}

	return nil
}

func (s *store) ListUserRoles(domainID, userID string) ([]string, error) {
	s.log.Debugf("List User Roles SQL: %s Params: domain_id: %s, user_id: %s", s.unprepared[FindUserRoles], domainID, userID)
	rows, err := s.stmts[FindUserRoles].Query(domainID, userID)
	if err != nil {
		return nil, exception.NewInternalServerError("query project list error, %s", err)
	}
	defer rows.Close()

	roles := []string{}
	for rows.Next() {
		var roleName string
		if err := rows.Scan(&roleName); err != nil {
			return nil, exception.NewInternalServerError("scan project record error, %s", err)
		}
		roles = append(roles, roleName)
	}

	return roles, nil
}

func (s *store) checkUserRoleExist(domainID, userID, roleName string) (bool, error) {
	var name string
	if err := s.stmts[CheckUserRoleIsBind].QueryRow(domainID, userID, roleName).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("query user role exist error, %s", err)
	}

	return true, nil
}

func (s *store) hmacHash(msg string) string {
	mac := hmac.New(sha256.New, []byte(s.key))
	io.WriteString(mac, msg)

	return fmt.Sprintf("%x", mac.Sum(nil))
}
