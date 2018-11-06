package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/http/context"
	"github.com/defineiot/keyauth/api/http/request"
	"github.com/defineiot/keyauth/api/http/response"
	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/internal/exception"
)

// RegistryUser use to create domain
func RegistryUser(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	email := val.Get("email").ToString()
	phone := val.Get("phone").ToString()
	code := val.Get("code").ToString()
	name := val.Get("username").ToString()
	pass := val.Get("password").ToString()

	if name == "" || pass == "" {
		response.Failed(w, exception.NewBadRequest("user_name or password is missed"))
		return
	}

	if email == "" && phone == "" {
		response.Failed(w, exception.NewBadRequest("email or phone is missed"))
		return
	}

	if email != "" && phone != "" {
		response.Failed(w, exception.NewBadRequest("email and phone only choice one"))
		return
	}

	if code == "" {
		response.Failed(w, exception.NewBadRequest("verify code missed"))
		return
	}
	codeI, err := strconv.Atoi(code)
	if err != nil {
		response.Failed(w, exception.NewBadRequest("parse code error, %s", err))
		return
	}

	// check verify code
	if email != "" {
		vc, err := global.Store.GetVerifyCode(email, "", codeI)
		if err != nil {
			response.Failed(w, err)
			return
		}

		if time.Now().Unix() > vc.ExpireAt {
			response.Failed(w, exception.NewBadRequest("verify code expired"))
			return
		}

		if err := global.Store.RevolkCode(vc.ID); err != nil {
			response.Failed(w, err)
			return
		}
	}

	if phone != "" {
		vc, err := global.Store.GetVerifyCode("", phone, codeI)
		if err != nil {
			response.Failed(w, err)
			return
		}

		if time.Now().Unix() > vc.ExpireAt {
			response.Failed(w, exception.NewBadRequest("verify code expired"))
			return
		}

		if err := global.Store.RevolkCode(vc.ID); err != nil {
			response.Failed(w, err)
			return
		}
	}

	// check user exist
	ok, err := global.Store.CheckUserIsGlobalExist(name)
	if err != nil {
		response.Failed(w, exception.NewInternalServerError("check username exist error, %s", err.Error()))
		return
	}
	if ok {
		response.Failed(w, exception.NewBadRequest("username: %s is exist", name))
		return
	}

	// create user domain
	dn := fmt.Sprintf("domain-%s", name)
	desc := fmt.Sprintf("user %s domain", name)
	domain, err := global.Store.CreateDomain(dn, desc, dn, true)
	if err != nil {
		response.Failed(w, exception.NewInternalServerError("automatic create user domain error, %s", err.Error()))
		return
	}

	// 交给业务控制层处理
	user, err := global.Store.CreateUser(domain.ID, name, pass, true, 6360, 6360)
	if err != nil {
		response.Failed(w, err)
		return
	}

	if err := global.Store.BindRole(domain.ID, user.ID, "domain_admin"); err != nil {
		response.Failed(w, err)
		return
	}

	if err := global.Store.BindRole(domain.ID, user.ID, "operator"); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, user)
	return
}

// IssueVerifyCode  issue verify code
func IssueVerifyCode(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	email := qs.Get("email")
	phone := qs.Get("phone")
	user := qs.Get("user")

	if email == "" && phone == "" {
		response.Failed(w, exception.NewBadRequest("email or phone required one"))
		return
	}

	vc, err := global.Store.IssueVerifyCode(email, phone)
	if err != nil {
		response.Failed(w, err)
		return
	}

	if vc.EmailAddress != "" {
		// mail, err := global.Conf.GetMailer()
		mail := global.Conf.Mail
		err := mail.Init()
		if err != nil {
			response.Failed(w, err)
			return
		}

		if mail == nil {
			response.Failed(w, exception.NewInternalServerError("mail config error not initial success"))
			return
		}

		if user == "" {
			user = vc.EmailAddress
		}

		if mail.Send(vc.EmailAddress, user, vc.Code); err != nil {
			response.Failed(w, err)
			return
		}
	}

	if vc.PhoneNumber != "" {
		sms := global.Conf.SMS
		if err := sms.SendSms(vc.PhoneNumber, strconv.Itoa(vc.Code)); err != nil {
			response.Failed(w, exception.NewInternalServerError("send verify code by sms error, %s", err))
			return
		}
	}

	response.Success(w, http.StatusCreated, nil)
	return
}

// InvitationsUser todo
func InvitationsUser(w http.ResponseWriter, r *http.Request) {
	tk := context.GetTokenFromContext(r)
	// 读取body数据
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		response.Failed(w, exception.NewBadRequest("read request body error, %s", err))
		return
	}

	data := map[string][]string{}
	if err := json.Unmarshal(body, &data); err != nil {
		response.Failed(w, exception.NewBadRequest("parse body json error, %s", err))
		return
	}

	invitatioRoles := []string{}
	if roles, ok := data["roles"]; ok {
		invitatioRoles = roles
	}

	invitationProjects := []string{}
	if projects, ok := data["project"]; ok {
		invitatioRoles = projects
	}

	code, err := global.Store.InvitationUser(tk.DomainID, tk.UserID, invitatioRoles, invitationProjects)
	if err != nil {
		response.Failed(w, err)
		return
	}

	uri := fmt.Sprintf("http://%s/v1/invitations/users/%s/code/%s/", r.Host, code.InviterID, code.Code)
	code.InvitationURL = uri

	response.Success(w, http.StatusCreated, code)
	return
}

// RevolkInvitation todo
func RevolkInvitation(w http.ResponseWriter, r *http.Request) {

}

// ListInvitationsRecords todo
func ListInvitationsRecords(w http.ResponseWriter, r *http.Request) {

}

// GetInvitationsRecord todo
func GetInvitationsRecord(w http.ResponseWriter, r *http.Request) {

}

// AcceptInvitation todo
func AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	tk := context.GetTokenFromContext(r)
	uid := tk.UserID
	did := tk.DomainID

	if uid == "" {
		response.Failed(w, exception.NewBadRequest("invited user id missed"))
		return
	}

	inviterID := ps.ByName("uid")
	invitationCode := ps.ByName("code")

	ir, err := global.Store.AcceptInvitation(inviterID, invitationCode, uid, did)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, ir)
	return
}

// CreateUser use to create domain
func CreateUser(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	domainName := val.Get("domain_name").ToString()
	name := val.Get("user_name").ToString()
	pass := val.Get("password").ToString()

	if name == "" || pass == "" {
		response.Failed(w, exception.NewBadRequest("user_name or password is missed"))
		return
	}

	// only system admin can create user in other domain
	tk := context.GetTokenFromContext(r)

	did := ""
	if domainName != "" {
		dom, err := global.Store.GetDomain("name", domainName)
		if err != nil {
			response.Failed(w, err)
			return
		}
		if dom == nil {
			response.Failed(w, exception.NewBadRequest("domain %s not found", domainName))
			return
		}

		if !tk.IsSystemAdmin && tk.DomainID != dom.ID {
			response.Failed(w, exception.NewForbidden("your create user must be in your own domain"))
			return
		}
		did = dom.ID
	} else {
		did = tk.DomainID
	}

	// check user exist
	ok, err := global.Store.CheckUserIsGlobalExist(name)
	if err != nil {
		response.Failed(w, exception.NewInternalServerError("check username exist error, %s", err.Error()))
		return
	}
	if ok {
		response.Failed(w, exception.NewBadRequest("username: %s is exist", name))
		return
	}

	// 交给业务控制层处理
	user, err := global.Store.CreateUser(did, name, pass, true, 6360, 6360)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, user)
	return
}

// CreateDomainUser todo
func CreateDomainUser(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	name := val.Get("user_name").ToString()
	pass := val.Get("password").ToString()

	if name == "" || pass == "" {
		response.Failed(w, exception.NewBadRequest("user_name or password is missed"))
		return
	}

	// only system admin can create user in other domain
	tk := context.GetTokenFromContext(r)

	if !tk.IsSystemAdmin {
		response.Failed(w, exception.NewForbidden("only system admin can create domain admin user"))
		return
	}

	// check user exist
	ok, err := global.Store.CheckUserIsGlobalExist(name)
	if err != nil {
		response.Failed(w, exception.NewInternalServerError("check username exist error, %s", err.Error()))
		return
	}
	if ok {
		response.Failed(w, exception.NewBadRequest("username: %s is exist", name))
		return
	}

	// create user domain
	dn := fmt.Sprintf("domain-%s", name)
	desc := fmt.Sprintf("user %s domain", name)
	domain, err := global.Store.CreateDomain(dn, desc, dn, true)
	if err != nil {
		response.Failed(w, exception.NewInternalServerError("automatic create user domain error, %s", err.Error()))
		return
	}

	// 交给业务控制层处理
	user, err := global.Store.CreateUser(domain.ID, name, pass, true, 6360, 6360)
	if err != nil {
		response.Failed(w, err)
		return
	}

	if err := global.Store.BindRole(domain.ID, user.ID, "domain_admin"); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, user)
}

// SetUserDefaultProject yes
func SetUserDefaultProject(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	pid := ps.ByName("pid")

	tk := context.GetTokenFromContext(r)
	uid := tk.UserID

	err := global.Store.SetUserDefaultProject(tk.DomainID, uid, pid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, nil)
	return
}

// GetUser use to get an user information
func GetUser(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	tk := context.GetTokenFromContext(r)
	uid := ps.ByName("uid")

	u, err := global.Store.GetUser(tk.DomainID, uid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// check permisson
	if tk.DomainID != u.DomainID {
		response.Failed(w, exception.NewForbidden("this user not belong to you"))
	}

	response.Success(w, http.StatusOK, u)
	return
}

// SetUserPassword use to set user's password
func SetUserPassword(w http.ResponseWriter, r *http.Request) {
	val, err := request.CheckObjectBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	oldPass := val.Get("old_password").ToString()
	newPass := val.Get("new_password").ToString()

	if oldPass == "" || newPass == "" {
		response.Failed(w, exception.NewBadRequest("old_password or new_password must be \"\""))
		return
	}

	if oldPass == newPass {
		response.Failed(w, exception.NewBadRequest("new_password must be different from old_password"))
		return
	}

	// check permisson
	uid := context.GetParamsFromContext(r).ByName("uid")
	tk := context.GetTokenFromContext(r)

	if !tk.IsDomainAdmin || !tk.IsSystemAdmin {
		if uid != tk.UserID {
			response.Failed(w, exception.NewForbidden("only domain_admin or system_admin can set other user password"))
			return
		}
	}

	if err := global.Store.SetUserPassword(uid, oldPass, newPass); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "")
	return
}

// ListProjectUser list all users can filter by domain or project
func ListProjectUser(w http.ResponseWriter, r *http.Request) {
	var err error

	ps := context.GetParamsFromContext(r)
	tk := context.GetTokenFromContext(r)
	pid := ps.ByName("pid")

	users := []*user.User{}

	proj, err := global.Store.GetProject(pid)
	if err != nil {
		response.Failed(w, err)
		return
	}
	if proj.DomainID != tk.DomainID {
		response.Failed(w, exception.NewForbidden("the project not belonge to you"))
		return
	}
	users, err = global.Store.ListProjectUser(pid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, users)
	return
}

// ListDomainUser list all users can filter by domain or project
func ListDomainUser(w http.ResponseWriter, r *http.Request) {
	var err error

	tk := context.GetTokenFromContext(r)
	did := tk.DomainID

	users := []*user.User{}

	users, err = global.Store.ListDomainUser(did)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusOK, users)
	return
}

// DeleteUser delete an user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	tk := context.GetTokenFromContext(r)
	uid := ps.ByName("uid")

	// user can't delete self, when delete self your ungristry api
	if tk.UserID == uid {
		response.Failed(w, exception.NewForbidden("your can't delete your self, but your can unregistry your"))
		return
	}

	u, err := global.Store.GetUser(tk.DomainID, uid)
	if err != nil {
		response.Failed(w, err)
		return
	}
	// 1. check user is in your domain
	if u.DomainID != tk.DomainID {
		response.Failed(w, exception.NewForbidden("this user: %s not in your domain", uid))
		return
	}

	// 2. the initial system admin can't delete
	if u.Name == global.Conf.Admin.UserName && tk.IsSystemAdmin {
		response.Failed(w, exception.NewForbidden("can't delete system initial admin user"))
		return
	}

	if err := global.Store.DeleteUser(tk.DomainID, uid); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}

// UnRegistry delete an user
func UnRegistry(w http.ResponseWriter, r *http.Request) {
	tk := context.GetTokenFromContext(r)

	u, err := global.Store.GetUser(tk.DomainID, tk.UserID)
	if err != nil {
		response.Failed(w, err)
		return
	}

	if u.Name == global.Conf.Admin.UserName && tk.IsSystemAdmin {
		response.Failed(w, exception.NewForbidden("can't unregistry system initial admin user"))
		return
	}

	if !tk.IsDomainAdmin {
		response.Failed(w, exception.NewForbidden("only domain admin user can unregistry"))
		return
	}

	if err := global.Store.DeleteUser(tk.DomainID, tk.UserID); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, "")
	return
}

// AddProjectsToUser todo
func AddProjectsToUser(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")

	iter, err := request.CheckArrayBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	pids := make([]string, 0)
	for iter.ReadArray() {
		if pid := iter.ReadString(); pid != "" {
			pids = append(pids, pid)
		}
	}
	if iter.Error != nil {
		response.Failed(w, exception.NewBadRequest("get userid from body array error, %s", iter.Error))
		return
	}

	if len(pids) == 0 {
		response.Failed(w, exception.NewBadRequest(`body array hase no project id`))
		return
	}

	// check permisson
	tk := context.GetTokenFromContext(r)
	u, err := global.Store.GetUser(tk.DomainID, tk.UserID)
	if err != nil {
		response.Failed(w, err)
		return
	}
	if u.DomainID != tk.DomainID {
		response.Failed(w, exception.NewForbidden("the user not belong to you"))
		return
	}

	// check needed add project is belong to you
	projects, err := global.Store.ListDomainProjects(tk.DomainID)
	notInDomainProjects := []string{}
	if err != nil {
		response.Failed(w, err)
		return
	}
	for _, inpid := range pids {
		var isIn bool
		for _, p := range projects {
			if inpid == p.ID {
				isIn = true
				break
			}
		}

		if !isIn {
			notInDomainProjects = append(notInDomainProjects, inpid)
		}
	}
	if len(notInDomainProjects) > 0 {
		response.Failed(w, exception.NewForbidden("the projects: %v not in you domain", notInDomainProjects))
		return
	}

	if err := global.Store.AddProjectsToUser(tk.DomainID, uid, pids...); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, nil)
	return
}

// RemoveProjectsFromUser todo
func RemoveProjectsFromUser(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	uid := ps.ByName("uid")

	iter, err := request.CheckArrayBody(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	pids := make([]string, 0)
	for iter.ReadArray() {
		if pid := iter.ReadString(); pid != "" {
			pids = append(pids, pid)
		}
	}
	if iter.Error != nil {
		response.Failed(w, exception.NewBadRequest("get userid from body array error, %s", iter.Error))
		return
	}

	if len(pids) == 0 {
		response.Failed(w, exception.NewBadRequest(`body array hase no project id`))
		return
	}

	// check permisson
	tk := context.GetTokenFromContext(r)
	u, err := global.Store.GetUser(tk.DomainID, tk.UserID)
	if err != nil {
		response.Failed(w, err)
		return
	}
	if u.DomainID != tk.DomainID {
		response.Failed(w, exception.NewForbidden("the user not belong to you"))
		return
	}

	if err := global.Store.RemoveProjectsFromUser(tk.DomainID, uid, pids...); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, nil)
	return

}

// BindRole todo
func BindRole(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	tk := context.GetTokenFromContext(r)
	uid := ps.ByName("uid")
	rn := ps.ByName("rn")

	u, err := global.Store.GetUser(tk.DomainID, uid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	if !tk.IsSystemAdmin {
		// 1. check user is in your domain
		if u.DomainID != tk.DomainID {
			response.Failed(w, exception.NewForbidden("this user: %s not in your domain", uid))
			return
		}

		// 2. forbidden other role bind system_admin
		if rn == "system_admin" {
			response.Failed(w, exception.NewForbidden("only system admin can bind system_admin role"))
			return
		}
	}

	if err := global.Store.BindRole(tk.DomainID, uid, rn); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "")
	return

}

// UnBindRole todo
func UnBindRole(w http.ResponseWriter, r *http.Request) {
	ps := context.GetParamsFromContext(r)
	tk := context.GetTokenFromContext(r)
	uid := ps.ByName("uid")
	rn := ps.ByName("rn")

	u, err := global.Store.GetUser(tk.DomainID, uid)
	if err != nil {
		response.Failed(w, err)
		return
	}

	if !tk.IsSystemAdmin {
		// 1. check user is in your domain
		if u.DomainID != tk.DomainID {
			response.Failed(w, exception.NewForbidden("this user: %s not in your domain", uid))
			return
		}
	} else {
		// 2. the initial system admin can't ubind
		if u.Name == global.Conf.Admin.UserName {
			response.Failed(w, exception.NewForbidden("can't unbind system initial admin user role"))
			return
		}
	}

	if err := global.Store.UnBindRole(tk.DomainID, uid, rn); err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "")
	return
}
