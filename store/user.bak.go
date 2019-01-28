package store

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/defineiot/keyauth/dao/verifycode"
)

// IssueVerifyCode todo
func (s *Store) IssueVerifyCode(vf *verifycode.VerifyCode) error {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(8000) + rand.Intn(1000) + 1000
	vf.Code = strconv.Itoa(code)
	return s.dao.VerifyCode.CreateVerifyCode(vf)
}

// GetVerifyCode todo
func (s *Store) GetVerifyCode(purpose verifycode.CodePurpose, target string) (*verifycode.VerifyCode, error) {
	return s.dao.VerifyCode.GetVerifyCode(purpose, target)
}

// RevolkCode todo
func (s *Store) RevolkCode(purpose verifycode.CodePurpose, target string, code string) error {
	return s.dao.VerifyCode.DeleteVerifyCode(purpose, target, code)
}

// InvitationUser todo
// func (s *Store) InvitationUser(inviterDomainID, inviterID string, invitedRoles, accessProjects []string) (*verifycode.VerifyCode, error) {
// 	// 1. check roles in validated
// 	roles, err := s.role.ListRole()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, inR := range invitedRoles {
// 		var isOK bool
// 		for _, r := range roles {
// 			if r.Name == inR {
// 				isOK = true
// 				break
// 			}
// 		}

// 		if !isOK {
// 			return nil, exception.NewBadRequest("role: %s not found", inR)
// 		}
// 	}

// 	// 2. check the projects has own by inviter
// 	projects, err := s.user.ListUserProjects(inviterDomainID, inviterID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, inP := range accessProjects {
// 		var isOK bool
// 		for _, p := range projects {
// 			if p == inP {
// 				isOK = true
// 				break
// 			}
// 		}

// 		if !isOK {
// 			return nil, exception.NewBadRequest("project: %s not found", inP)
// 		}
// 	}

// 	// 3. ok, start generate invitation code
// 	return s.user.SaveInvitationsRecord(inviterID, invitedRoles, accessProjects)
// }

// // ListUserInvitations todo
// func (s *Store) ListUserInvitations(inviterID string) ([]*user.InvitationCode, error) {
// 	return s.user.ListInvitationRecord(inviterID)
// }

// // GetInvitation todo
// func (s *Store) GetInvitation(inviterID, code string) (*user.InvitationCode, error) {
// 	return s.user.GetInvitationRecord(inviterID, code)
// }

// // RevolkInvitation todo
// func (s *Store) RevolkInvitation(id int64) error {
// 	return s.user.DeleteInvitationRecord(id)
// }

// // AcceptInvitation todo
// func (s *Store) AcceptInvitation(inviterID, invitationCode, userID, userDomainID string) (*user.InvitationCode, error) {
// 	inviter, err := s.user.GetUserByID(inviterID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if inviter.DomainID == userDomainID {
// 		return nil, exception.NewForbidden("can't invite in your domain's user")
// 	}

// 	code, err := s.user.GetInvitationRecord(inviterID, invitationCode)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if time.Now().Unix() > code.ExpireTime {
// 		return nil, exception.NewBadRequest("registry code expired, expire time: %d", code.ExpireTime)
// 	}

// 	code.InvitedUserID = userID
// 	code.InvitedUserDomainID = userDomainID

// 	// 1. link to inviter's domain
// 	if err := s.user.SaveUserOtherDomain(userID, inviter.DomainID); err != nil {
// 		return nil, err
// 	}

// 	// 2. bind other domain's roles
// 	for _, rn := range code.InvitedUserRoleNames {
// 		if err := s.user.BindRole(inviter.DomainID, code.InvitedUserID, rn); err != nil {
// 			return nil, err
// 		}
// 	}

// 	// 3. add to projects
// 	if err := s.user.AddProjectsToUser(inviter.DomainID, code.InvitedUserID, code.AccessProjects...); err != nil {
// 		return nil, err
// 	}

// 	// 4. update invitation record
// 	if err := s.user.UpdateInvitationsRecord(code); err != nil {
// 		return nil, err
// 	}

// 	return code, nil
// }

// CheckUserIsGlobalExist use to create user
// func (s *Store) CheckUserIsGlobalExist(username string) (bool, error) {
// 	return s.user.CheckUserNameIsGlobalExist(username)
// }

// SetUserDefaultProject todo
// func (s *Store) SetUserDefaultProject(domainID, userID, projectID string) error {
// 	// 1. check project is exist
// 	ok, err := s.project.CheckProjectIsExistByID(projectID)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return exception.NewBadRequest("project %s not exist", projectID)
// 	}

// 	err = s.user.SetDefaultProject(domainID, userID, projectID)
// 	if err != nil {
// 		return err
// 	}

// 	if s.isCache {
// 		cacheKey := "user_" + userID
// 		if !s.cache.Delete(cacheKey) {
// 			s.log.Debugf("delete user from cache failed, key: %s", cacheKey)
// 		}
// 		s.log.Debugf("delete user from cache success, key: %s", cacheKey)
// 	}

// 	return nil
// }

// // SetUserPassword use to change user password
// func (s *Store) SetUserPassword(userID, oldPass, newPass string) error {
// 	cacheKey := "user_" + userID

// 	if err := s.user.SetUserPassword(userID, oldPass, newPass); err != nil {
// 		return err
// 	}

// 	if s.isCache {
// 		if !s.cache.Delete(cacheKey) {
// 			s.log.Debugf("delete user from cache failed, key: %s", cacheKey)
// 		}
// 		s.log.Debugf("delete user from cache success, key: %s", cacheKey)
// 	}
// 	return nil
// }

// // AddProjectsToUser add project to user
// func (s *Store) AddProjectsToUser(domainID, userID string, projectIDs ...string) error {
// 	if err := s.checkProjectExist(projectIDs...); err != nil {
// 		return err
// 	}

// 	return s.user.AddProjectsToUser(domainID, userID, projectIDs...)
// }

// // RemoveProjectsFromUser remove project ot user
// func (s *Store) RemoveProjectsFromUser(domainID, userID string, projectIDs ...string) error {
// 	if err := s.checkProjectExist(projectIDs...); err != nil {
// 		return err
// 	}

// 	return s.user.RemoveProjectsFromUser(domainID, userID, projectIDs...)
// }

// // BindRole todo
// func (s *Store) BindRole(domainID, userID, roleName string) error {
// 	ok, err := s.role.CheckRoleExist(roleName)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return exception.NewBadRequest("role: %s not exist", roleName)
// 	}

// 	cacheKey := "user_" + userID
// 	if s.isCache {
// 		if !s.cache.Delete(cacheKey) {
// 			s.log.Debugf("delete user from cache failed, key: %s", cacheKey)
// 		}
// 		s.log.Debugf("delete user from cache success, key: %s", cacheKey)
// 	}

// 	return s.user.BindRole(domainID, userID, roleName)
// }

// // UnBindRole todo
// func (s *Store) UnBindRole(domainID, userID, roleName string) error {
// 	ok, err := s.role.CheckRoleExist(roleName)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return exception.NewBadRequest("role: %s not exist", roleName)
// 	}

// 	cacheKey := "user_" + userID
// 	if s.isCache {
// 		if !s.cache.Delete(cacheKey) {
// 			s.log.Debugf("delete user from cache failed, key: %s", cacheKey)
// 		}
// 		s.log.Debugf("delete user from cache success, key: %s", cacheKey)
// 	}

// 	return s.user.UnBindRole(domainID, userID, roleName)
// }

// // ListUserDomain get an user
// func (s *Store) ListUserDomain(userID string) ([]*domain.Domain, error) {
// 	var err error

// 	domains := []*domain.Domain{}
// 	cacheKey := "user_" + userID

// 	if s.isCache {
// 		if s.cache.Get(cacheKey, domains) {
// 			s.log.Debugf("get project from cache key: %s", cacheKey)
// 			return domains, nil
// 		}
// 		s.log.Debugf("get project from cache failed, key: %s", cacheKey)
// 	}

// 	u, err := s.user.GetUserByID(userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if u == nil {
// 		return nil, exception.NewBadRequest("user %s not found", userID)
// 	}

// 	// query user's domain
// 	d, err := s.domain.GetDomainByID(u.DomainID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	domains = append(domains, d)

// 	// query user's other domain
// 	dids, err := s.user.ListUserOtherDomains(userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, did := range dids {
// 		d, err := s.domain.GetDomainByID(did)
// 		if err != nil {
// 			return nil, err
// 		}
// 		domains = append(domains, d)
// 	}

// 	if s.isCache {
// 		if !s.cache.Set(cacheKey, u, s.ttl) {
// 			s.log.Debugf("set user cache failed, key: %s", cacheKey)
// 		}
// 		s.log.Debugf("set user cache ok, key: %s", cacheKey)
// 	}

// 	return domains, nil
// }

// func (s *Store) checkProjectExist(projectIDs ...string) error {
// 	errs := make([]string, 0)
// 	noexist := make([]string, 0)

// 	for _, pid := range projectIDs {
// 		ok, err := s.project.CheckProjectIsExistByID(pid)
// 		if err != nil {
// 			errs = append(errs, err.Error())
// 		}
// 		if !ok {
// 			noexist = append(noexist, pid)
// 		}
// 	}

// 	if len(errs) != 0 {
// 		err := strings.Join(errs, ",")
// 		return errors.New(err)
// 	}

// 	if len(noexist) != 0 {
// 		neu := strings.Join(noexist, ",")
// 		return exception.NewBadRequest("project %s not exist", neu)
// 	}

// 	return nil
// }
