package mysql

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/internal/exception"
)

func (s *store) CreateUser(u *models.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := s.CheckUserNameIsExist(u.Domain.ID, u.Account); err != nil {
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
	dpid := ""
	if u.DefaultProject != nil {
		dpid = u.DefaultProject.ID
	}
	// 存入
	if _, err = userPre.Exec(u.ID, u.Department.ID, u.Account, u.Mobile, u.Email,
		u.Phone, u.Address, u.RealName, u.NickName, u.Gender, u.Avatar, u.Language,
		u.City, u.Province, u.Locked, u.Domain.ID, u.CreateAt, u.ExpiresActiveDays,
		dpid); err != nil {
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

		pass := models.Password{
			CreateAt: time.Now().Unix(),
			ExpireAt: u.Password.ExpireAt,
			Password: u.Password.Password,
			UserID:   u.ID,
		}
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

func (s *store) ListDomainUsers(domainID string) ([]*models.User, error) {
	rows, err := s.stmts[FindDomainUsers].Query(domainID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user list error, %s", err)
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		u := new(models.User)
		pass := new(models.Password)
		u.Password = pass
		u.DefaultProject = new(models.Project)
		u.Domain = new(models.Domain)
		u.Department = new(models.Department)
		if err := rows.Scan(&u.ID, &u.Department.ID, &u.Account, &u.Mobile, &u.Email, &u.Phone, &u.Address,
			&u.RealName, &u.NickName, &u.Gender, &u.Avatar, &u.Language, &u.City, &u.Province,
			&u.Locked, &u.Domain.ID, &u.CreateAt, &u.ExpiresActiveDays, &u.DefaultProject.ID,
			&pass.Password, &pass.ExpireAt, &pass.CreateAt, &pass.UpdateAt); err != nil {
			return nil, exception.NewInternalServerError("scan project's user id error, %s", err)
		}

		users = append(users, u)
	}

	return users, nil
}

func (s *store) ListDepartmentUsers(departmentID string) ([]*models.User, error) {
	rows, err := s.stmts[FindDepartmentUsers].Query(departmentID)
	if err != nil {
		return nil, exception.NewInternalServerError("query user list error, %s", err)
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		u := new(models.User)
		pass := new(models.Password)
		u.Password = pass
		u.DefaultProject = new(models.Project)
		u.Domain = new(models.Domain)
		u.Department = new(models.Department)
		if err := rows.Scan(&u.ID, &u.Department.ID, &u.Account, &u.Mobile, &u.Email, &u.Phone, &u.Address,
			&u.RealName, &u.NickName, &u.Gender, &u.Avatar, &u.Language, &u.City, &u.Province,
			&u.Locked, &u.Domain.ID, &u.CreateAt, &u.ExpiresActiveDays, &u.DefaultProject.ID,
			&pass.Password, &pass.ExpireAt, &pass.CreateAt, &pass.UpdateAt); err != nil {
			return nil, exception.NewInternalServerError("scan project's user id error, %s", err)
		}

		users = append(users, u)
	}

	return users, nil
}

func (s *store) ListProjectUsers(projectID string) ([]*models.User, error) {
	rows, err := s.stmts[FindProjectUsers].Query(projectID)
	if err != nil {
		return nil, exception.NewInternalServerError("query project user list error, %s", err)
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		u := new(models.User)
		pass := new(models.Password)
		u.Password = pass
		u.DefaultProject = new(models.Project)
		u.Domain = new(models.Domain)
		u.Department = new(models.Department)
		if err := rows.Scan(&u.ID, &u.Department.ID, &u.Account, &u.Mobile, &u.Email, &u.Phone, &u.Address,
			&u.RealName, &u.NickName, &u.Gender, &u.Avatar, &u.Language, &u.City, &u.Province,
			&u.Locked, &u.Domain.ID, &u.CreateAt, &u.ExpiresActiveDays, &u.DefaultProject.ID,
			&pass.Password, &pass.ExpireAt, &pass.CreateAt, &pass.UpdateAt); err != nil {
			return nil, exception.NewInternalServerError("scan project's user id error, %s", err)
		}

		users = append(users, u)
	}

	return users, nil
}

func (s *store) GetUser(index models.FoundUserIndex, value string) (*models.User, error) {
	var row *sql.Row

	u := new(models.User)
	pass := new(models.Password)
	u.Password = pass
	u.DefaultProject = new(models.Project)
	u.Domain = new(models.Domain)
	u.Department = new(models.Department)
	switch index {
	case models.UserIDIndex:
		row = s.stmts[FindUserByID].QueryRow(value)
	case models.AccountIndex:
		row = s.stmts[FindUserByAccount].QueryRow(value)
	case models.MobileIndex:
		row = s.stmts[FindUserByMobile].QueryRow(value)
	case models.EmailIndex:
		row = s.stmts[FindUserByEmail].QueryRow(value)
	default:
		return nil, exception.NewBadRequest("the user's %s index not found", index)
	}
	if err := row.Scan(&u.ID, &u.Department.ID, &u.Account, &u.Mobile,
		&u.Email, &u.Phone, &u.Address, &u.RealName, &u.NickName, &u.Gender, &u.Avatar, &u.Language,
		&u.City, &u.Province, &u.Locked, &u.Domain.ID, &u.CreateAt, &u.ExpiresActiveDays, &u.DefaultProject.ID,
		&pass.Password, &pass.ExpireAt, &pass.CreateAt, &pass.UpdateAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("user %s not find", value)
		}
		return nil, exception.NewInternalServerError("query single user error, %s", err)
	}

	return u, nil
}

func (s *store) DeleteUser(domainID, userID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return exception.NewInternalServerError("start delete user transaction error, %s", err)
	}

	// 删除用户
	deleteUserPre, err := tx.Prepare("DELETE FROM users WHERE id = ?")
	defer deleteUserPre.Close()
	if err != nil {
		return exception.NewInternalServerError("prepare delete user stmt error, %s", err)
	}

	ret, err := deleteUserPre.Exec(userID)
	if err != nil {
		return exception.NewInternalServerError("delete user exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("user %s not exist", userID)
	}

	// 删除用户的密码
	deletePassPre, err := tx.Prepare("DELETE FROM passwords WHERE user_id = ?")
	defer deletePassPre.Close()
	if err != nil {
		return exception.NewInternalServerError("prepare delete user pass stmt error, %s", err)
	}
	ret, err = deletePassPre.Exec(userID)
	if err != nil {
		return exception.NewInternalServerError("delete user pass exec sql error, %s", err)
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return exception.NewInternalServerError("delete user transaction rollback error, %s", err)
		}
		return exception.NewInternalServerError("delete user transaction commit error, but rollback success, %s", err)
	}

	return nil
}

func (s *store) BindRole(domainID, userID, roleID string) error {
	ok, err := s.CheckUserIsExistByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("user %s not exist", userID)
	}

	ok, err = s.checkUserRoleExist(domainID, userID, roleID)
	if err != nil {
		return err
	}
	if ok {
		return exception.NewBadRequest("user %s has bind the role: %s", userID, roleID)
	}

	_, err = s.stmts[BindRole].Exec(domainID, userID, roleID)
	if err != nil {
		return exception.NewInternalServerError("insert role user mapping exec sql err, %s", err)
	}

	return nil
}

func (s *store) UnBindRole(domainID, userID, roleID string) error {
	ok, err := s.CheckUserIsExistByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("user %s exist", userID)
	}

	ret, err := s.stmts[UnbindRole].Exec(domainID, userID, roleID)
	if err != nil {
		return exception.NewInternalServerError("delete role user mapping exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete role user mapping affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("the role: %s is not bind to user: %s", roleID, userID)
	}

	return nil
}

func (s *store) checkUserRoleExist(domainID, userID, roleID string) (bool, error) {
	var name string
	if err := s.stmts[CheckUserRoleIsBind].QueryRow(domainID, userID, roleID).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("query user role exist error, %s", err)
	}

	return true, nil
}

func (s *store) AddProjectsToUser(domainID, userID string, projects []*models.Project) error {
	// check the user is not exist
	ok, err := s.CheckUserIsExistByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("user %s not exist", userID)
	}

	// check projects is owned by this user
	pids, err := s.listUserProjectIDs(domainID, userID)
	if err != nil {
		return err
	}
	existPIDs := []string{}
	for _, pid := range pids {
		for _, inpid := range projects {
			if inpid.ID == pid {
				existPIDs = append(existPIDs, inpid.ID)
			}
		}
	}
	if len(existPIDs) != 0 {
		return exception.NewBadRequest("project %s is in this project", existPIDs)
	}

	for _, p := range projects {
		_, err = s.stmts[AddProjectToUser].Exec(userID, p.ID)
		if err != nil {
			return fmt.Errorf("insert add projects to user mapping exec sql err, %s", err)
		}
	}

	return nil
}

func (s *store) RemoveProjectsFromUser(domainID, userID string, projects []*models.Project) error {
	// check the user is not exist
	ok, err := s.CheckUserIsExistByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewBadRequest("user %s not exist", userID)
	}

	// check projects is owned by this user
	pids, err := s.listUserProjectIDs(domainID, userID)
	if err != nil {
		return err
	}
	nExistPIDs := []string{}
	for _, inpid := range projects {
		var ok bool
		for _, pid := range pids {
			if pid == inpid.ID {
				ok = true
			}
		}
		if !ok {
			nExistPIDs = append(nExistPIDs, inpid.ID)
		}
	}
	if len(nExistPIDs) != 0 {
		return exception.NewBadRequest("project %s isn't in this project", nExistPIDs)
	}

	for _, p := range projects {
		_, err = s.stmts[RemoveProjectsFromUser].Exec(userID, p.ID)
		if err != nil {
			return fmt.Errorf("insert remove projects to user mapping exec sql err, %s", err)
		}
	}

	return nil
}

func (s *store) listUserProjectIDs(domainID, userID string) ([]string, error) {
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
