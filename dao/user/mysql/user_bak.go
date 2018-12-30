package mysql

// func (s *store) SaveUserOtherDomain(userID, otherDomainID string) error {
// 	_, err := s.stmts[SaveUserOtherDomain].Exec(userID, otherDomainID, time.Now().Unix())
// 	if err != nil {
// 		return exception.NewInternalServerError("insert user other domain exec sql err, %s", err)
// 	}

// 	return nil
// }

// func (s *store) ListUserOtherDomains(userID string) ([]string, error) {
// 	ok, err := s.CheckUserIsExistByID(userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !ok {
// 		return nil, exception.NewBadRequest("user %s not exist", userID)
// 	}

// 	rows, err := s.stmts[FindUserOtherDomain].Query(userID)
// 	if err != nil {
// 		return nil, exception.NewInternalServerError("query user's other domain error, %s", err)
// 	}
// 	defer rows.Close()

// 	domains := []string{}
// 	for rows.Next() {
// 		var did string
// 		if err := rows.Scan(&did); err != nil {
// 			return nil, exception.NewInternalServerError("scan user's other domain id error, %s", err)
// 		}
// 		domains = append(domains, did)
// 	}
// 	return domains, nil
// }

// func (s *store) DeleteUserOtherDomain(userID, otherDomainID string) error {
// 	ret, err := s.stmts[DeleteUserOtherDomain].Exec(userID, otherDomainID)
// 	if err != nil {
// 		return exception.NewInternalServerError("delete user other domain exec sql error, %s", err)
// 	}
// 	count, err := ret.RowsAffected()
// 	if err != nil {
// 		return exception.NewInternalServerError("get delete row affected error, %s", err)
// 	}
// 	if count == 0 {
// 		return exception.NewBadRequest("user %s domain %s not exist", userID, otherDomainID)
// 	}

// 	return nil
// }

// func (s *store) ListUserProjects(domainID, userID string) ([]string, error) {
// 	ok, err := s.CheckUserIsExistByID(userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !ok {
// 		return nil, exception.NewBadRequest("user %s not exist", userID)
// 	}

// 	rows, err := s.stmts[FindUserProjects].Query(userID, domainID)
// 	if err != nil {
// 		return nil, exception.NewInternalServerError("query user's project id error, %s", err)
// 	}
// 	defer rows.Close()

// 	projects := []string{}
// 	for rows.Next() {
// 		var projectID string
// 		if err := rows.Scan(&projectID); err != nil {
// 			return nil, exception.NewInternalServerError("scan user's project id error, %s", err)
// 		}
// 		projects = append(projects, projectID)
// 	}
// 	return projects, nil
// }

// func (s *store) SetUserPassword(userID, oldPass, newPass string) error {
// 	ok, err := s.CheckUserIsExistByID(userID)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return exception.NewBadRequest("user %s not exist", userID)
// 	}

// 	pass, err := s.QueryPassword(userID)
// 	if err != nil {
// 		return err
// 	}

// 	if s.hmacHash(oldPass) != pass.Password {
// 		return exception.NewBadRequest("your old password is incorrect")
// 	}

// 	ret, err := s.stmts[SetUserPassword].Exec(s.hmacHash(newPass), userID)
// 	if err != nil {
// 		return exception.NewInternalServerError("set user's  password exec sql error, %s", err)
// 	}

// 	ra, err := ret.RowsAffected()
// 	if err != nil {
// 		s.log.Error(err)
// 	}
// 	if ra == 0 {
// 		return exception.NewInternalServerError("not row updated here.")
// 	}

// 	return nil
// }

// func (s *store) SetDefaultProject(domainID, userID, projectID string) error {
// 	ok, err := s.CheckUserIsExistByID(userID)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return exception.NewBadRequest("user %s not exist", userID)
// 	}

// 	projects, err := s.ListUserProjects(domainID, userID)
// 	if err != nil {
// 		return err
// 	}

// 	// check the project is user's project
// 	for _, pid := range projects {
// 		if pid == projectID {
// 			ok = true
// 		}
// 	}
// 	if !ok {
// 		return exception.NewBadRequest("user %s hasn't project %s", userID, projectID)
// 	}

// 	_, err = s.stmts[SetUserDefaultProject].Exec(projectID, userID)
// 	if err != nil {
// 		return exception.NewInternalServerError("set user's default project exec sql error, %s", err)
// 	}

// 	return nil
// }

// func (s *store) AddProjectsToUser(domainID, userID string, projectIDs ...string) error {
// 	// check the user is not exist
// 	ok, err := s.CheckUserIsExistByID(userID)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return exception.NewBadRequest("user %s not exist", userID)
// 	}

// 	// check projects is owned by this user
// 	pids, err := s.ListUserProjects(domainID, userID)
// 	if err != nil {
// 		return err
// 	}
// 	existPIDs := []string{}
// 	for _, pid := range pids {
// 		for _, inpid := range projectIDs {
// 			if inpid == pid {
// 				existPIDs = append(existPIDs, inpid)
// 			}
// 		}
// 	}
// 	if len(existPIDs) != 0 {
// 		return exception.NewBadRequest("project %s is in this project", existPIDs)
// 	}

// 	for _, projectID := range projectIDs {
// 		//
// 		_, err = s.stmts[AddProjectToUser].Exec(userID, projectID)
// 		if err != nil {
// 			return fmt.Errorf("insert add projects to user mapping exec sql err, %s", err)
// 		}
// 	}

// 	return nil
// }

// func (s *store) RemoveProjectsFromUser(domainID, userID string, projectIDs ...string) error {
// 	// check the user is not exist
// 	ok, err := s.CheckUserIsExistByID(userID)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return exception.NewBadRequest("user %s not exist", userID)
// 	}

// 	// check projects is owned by this user
// 	pids, err := s.ListUserProjects(domainID, userID)
// 	if err != nil {
// 		return err
// 	}
// 	nExistPIDs := []string{}
// 	for _, inpid := range projectIDs {
// 		var ok bool
// 		for _, pid := range pids {
// 			if pid == inpid {
// 				ok = true
// 			}
// 		}
// 		if !ok {
// 			nExistPIDs = append(nExistPIDs, inpid)
// 		}
// 	}
// 	if len(nExistPIDs) != 0 {
// 		return exception.NewBadRequest("project %s isn't in this project", nExistPIDs)
// 	}

// 	for _, projectID := range projectIDs {
// 		_, err = s.stmts[RemoveProjectsFromUser].Exec(userID, projectID)
// 		if err != nil {
// 			return fmt.Errorf("insert remove projects to user mapping exec sql err, %s", err)
// 		}
// 	}

// 	return nil
// }

// func (s *store) QueryPassword(userID string) (*user.Password, error) {
// 	pass := user.Password{}
// 	err := s.stmts[FindUserPassword].QueryRow(userID).Scan(
// 		&pass.Password, &pass.ExpireAt, &pass.CreateAt, &pass.UpdateAt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, exception.NewNotFound("password %s not find", userID)
// 		}

// 		return nil, exception.NewInternalServerError("query single password error, %s", err)
// 	}

// 	return &pass, nil
// }

// func (s *store) GetUserByName(domainID, userName string) (*user.User, error) {
// 	userI := user.User{}
// 	err := s.stmts[FindUserByName].QueryRow(userName, domainID).Scan(
// 		&userI.ID, &userI.Name, &userI.Enabled, &userI.LastActiveTime, &userI.DomainID, &userI.CreateAt, &userI.ExpireActiveDays, &userI.DefaultProjectID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, exception.NewNotFound("user %s not find", userName)
// 		}

// 		return nil, exception.NewInternalServerError("query single user error, %s", err)
// 	}

// 	// get user's emails
// 	emails, err := s.QueryEmail(userI.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	userI.Emails = emails

// 	// get user's phones
// 	phones, err := s.QueryPhone(userI.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	userI.Phones = phones

// 	// get user's password
// 	userI.Password = &user.Password{}

// 	return &userI, nil
// }

// func (s *store) ValidateUser(domainID, userName, password string) (string, error) {
// 	var uid string
// 	err := s.stmts[FindUserIDByName].QueryRow(userName, domainID).Scan(&uid)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return "", exception.NewNotFound("user %s not find", userName)
// 		}

// 		return "", exception.NewInternalServerError("query single user error, %s", err)
// 	}

// 	pass, err := s.QueryPassword(uid)
// 	if err != nil {
// 		return "", err
// 	}

// 	if s.hmacHash(password) == pass.Password {
// 		return uid, nil
// 	}

// 	return "", nil
// }

// func (s *store) ValidateGlobalUser(userName, password string) (string, error) {
// 	var uid string
// 	err := s.stmts[FindGlobalUserIDByName].QueryRow(userName).Scan(&uid)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return "", exception.NewNotFound("user %s not find", userName)
// 		}

// 		return "", exception.NewInternalServerError("query single user error, %s", err)
// 	}

// 	pass, err := s.QueryPassword(uid)
// 	if err != nil {
// 		return "", err
// 	}

// 	if s.hmacHash(password) == pass.Password {
// 		return uid, nil
// 	}

// 	return "", nil
// }

// func (s *store) CheckUserNameIsGlobalExist(userName string) (bool, error) {
// 	var uid string
// 	err := s.stmts[CheckUserGlobalExistByName].QueryRow(userName).Scan(&uid)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return false, nil
// 		}
// 		return false, exception.NewInternalServerError("check user exist by global name error, %s", err)
// 	}

// 	return true, nil
// }

// func (s *store) BindRole(domainID, userID, roleName string) error {
// 	ok, err := s.CheckUserIsExistByID(userID)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return exception.NewBadRequest("user %s not exist", userID)
// 	}

// 	ok, err = s.checkUserRoleExist(domainID, userID, roleName)
// 	if err != nil {
// 		return err
// 	}
// 	if ok {
// 		return exception.NewBadRequest("user %s has bind the role: %s", userID, roleName)
// 	}

// 	_, err = s.stmts[BindRole].Exec(domainID, userID, roleName)
// 	if err != nil {
// 		return exception.NewInternalServerError("insert role user mapping exec sql err, %s", err)
// 	}

// 	return nil
// }

// func (s *store) UnBindRole(domainID, userID, roleName string) error {
// 	ok, err := s.CheckUserIsExistByID(userID)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return exception.NewBadRequest("user %s exist", userID)
// 	}

// 	ret, err := s.stmts[UnbindRole].Exec(domainID, userID, roleName)
// 	if err != nil {
// 		return exception.NewInternalServerError("delete role user mapping exec sql error, %s", err)
// 	}
// 	count, err := ret.RowsAffected()
// 	if err != nil {
// 		return exception.NewInternalServerError("get delete role user mapping affected error, %s", err)
// 	}
// 	if count == 0 {
// 		return exception.NewBadRequest("the role: %s is not bind to user: %s", roleName, userID)
// 	}

// 	return nil
// }

// func (s *store) ListUserRoles(domainID, userID string) ([]string, error) {
// 	s.log.Debugf("List User Roles SQL: %s Params: domain_id: %s, user_id: %s", s.unprepared[FindUserRoles], domainID, userID)
// 	rows, err := s.stmts[FindUserRoles].Query(domainID, userID)
// 	if err != nil {
// 		return nil, exception.NewInternalServerError("query project list error, %s", err)
// 	}
// 	defer rows.Close()

// 	roles := []string{}
// 	for rows.Next() {
// 		var roleName string
// 		if err := rows.Scan(&roleName); err != nil {
// 			return nil, exception.NewInternalServerError("scan project record error, %s", err)
// 		}
// 		roles = append(roles, roleName)
// 	}

// 	return roles, nil
// }

// func (s *store) checkUserRoleExist(domainID, userID, roleName string) (bool, error) {
// 	var name string
// 	if err := s.stmts[CheckUserRoleIsBind].QueryRow(domainID, userID, roleName).Scan(&name); err != nil {
// 		if err == sql.ErrNoRows {
// 			return false, nil
// 		}

// 		return false, exception.NewInternalServerError("query user role exist error, %s", err)
// 	}

// 	return true, nil
// }
