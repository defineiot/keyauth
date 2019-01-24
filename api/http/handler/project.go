package handler

// import (
// 	"net/http"

// 	"github.com/defineiot/keyauth/api/global"
// 	"github.com/defineiot/keyauth/api/http/context"
// 	"github.com/defineiot/keyauth/api/http/request"
// 	"github.com/defineiot/keyauth/api/http/response"
// 	"github.com/defineiot/keyauth/dao/domain"
// 	"github.com/defineiot/keyauth/dao/project"
// 	"github.com/defineiot/keyauth/internal/exception"
// )

// // CreateProject use to create domain
// func CreateProject(w http.ResponseWriter, r *http.Request) {
// 	val, err := request.CheckObjectBody(r)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	name := val.Get("name").ToString()
// 	desc := val.Get("description").ToString()

// 	tk := context.GetTokenFromContext(r)
// 	did := tk.DomainID

// 	if name == "" {
// 		response.Failed(w, exception.NewBadRequest("name missed"))
// 		return
// 	}

// 	// 交给业务控制层处理
// 	project, err := global.Store.CreateProject(did, name, desc, true)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	response.Success(w, http.StatusCreated, project)
// 	return
// }

// // ListProject list users under his domain
// func ListProject(w http.ResponseWriter, r *http.Request) {
// 	var err error

// 	tk := context.GetTokenFromContext(r)
// 	did := tk.DomainID

// 	projects := []*project.Project{}

// 	projects, err = global.Store.ListDomainProjects(did)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	response.Success(w, http.StatusOK, projects)
// 	return
// }

// // ListUserDomain list users under his domain
// func ListUserDomain(w http.ResponseWriter, r *http.Request) {
// 	var err error

// 	tk := context.GetTokenFromContext(r)
// 	ps := context.GetParamsFromContext(r)
// 	uid := ps.ByName("uid")

// 	domains := []*domain.Domain{}

// 	// 1. get token user
// 	u, err := global.Store.GetUser(tk.DomainID, uid)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	// 2. check user is in your domain
// 	if u.DomainID != tk.DomainID {
// 		response.Failed(w, exception.NewForbidden("this user: %s not in your domain", uid))
// 		return
// 	}

// 	// only domain admin or system admin can list domain projects
// 	if !(tk.IsSystemAdmin || tk.IsDomainAdmin) && u.ID != tk.UserID {
// 		response.Failed(w, exception.NewForbidden("only domain_admin and system_admin can list others user's projects"))
// 		return
// 	}

// 	domains, err = global.Store.ListUserDomain(uid)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	response.Success(w, http.StatusOK, domains)
// 	return
// }

// // ListUserProject list users under his domain
// func ListUserProject(w http.ResponseWriter, r *http.Request) {
// 	var err error

// 	tk := context.GetTokenFromContext(r)
// 	ps := context.GetParamsFromContext(r)
// 	uid := ps.ByName("uid")

// 	projects := []*project.Project{}

// 	// 1. get token user
// 	u, err := global.Store.GetUser(tk.DomainID, uid)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	// 2. check user is in your domain
// 	if u.DomainID != tk.DomainID {
// 		response.Failed(w, exception.NewForbidden("this user: %s not in your domain", uid))
// 		return
// 	}

// 	// only domain admin or system admin can list domain projects
// 	if !(tk.IsSystemAdmin || tk.IsDomainAdmin) && u.ID != tk.UserID {
// 		response.Failed(w, exception.NewForbidden("only domain_admin and system_admin can list others user's projects"))
// 		return
// 	}

// 	projects, err = global.Store.ListUserProjects(tk.DomainID, uid)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	response.Success(w, http.StatusOK, projects)
// 	return
// }

// // GetProject use to create domain
// func GetProject(w http.ResponseWriter, r *http.Request) {

// 	ps := context.GetParamsFromContext(r)
// 	pid := ps.ByName("pid")

// 	proj, err := global.Store.GetProject(pid)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	// check permisson
// 	tk := context.GetTokenFromContext(r)
// 	if proj.DomainID != tk.DomainID {
// 		response.Failed(w, exception.NewForbidden("the project not belong to your"))
// 		return
// 	}

// 	response.Success(w, http.StatusOK, proj)
// 	return
// }

// // DeleteProject use to create domain
// func DeleteProject(w http.ResponseWriter, r *http.Request) {
// 	ps := context.GetParamsFromContext(r)
// 	pid := ps.ByName("pid")

// 	proj, err := global.Store.GetProject(pid)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	// check permisson
// 	tk := context.GetTokenFromContext(r)
// 	if proj.DomainID != tk.DomainID {
// 		response.Failed(w, exception.NewForbidden("the project not belong to your"))
// 		return
// 	}

// 	// TODO: get token from context, and check permission
// 	if err := global.Store.DeleteProjectByID(pid); err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	response.Success(w, http.StatusNoContent, "")
// 	return
// }

// // AddUsersToProject use to create domain
// func AddUsersToProject(w http.ResponseWriter, r *http.Request) {
// 	ps := context.GetParamsFromContext(r)
// 	pid := ps.ByName("pid")

// 	iter, err := request.CheckArrayBody(r)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	uids := make([]string, 0)
// 	for iter.ReadArray() {
// 		if uid := iter.ReadString(); uid != "" {
// 			uids = append(uids, uid)
// 		}
// 	}
// 	if iter.Error != nil {
// 		response.Failed(w, exception.NewBadRequest("get userid from body array error, %s", iter.Error))
// 		return
// 	}

// 	if len(uids) == 0 {
// 		response.Failed(w, exception.NewBadRequest(`body array hase no userid`))
// 		return
// 	}

// 	proj, err := global.Store.GetProject(pid)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	// check permisson
// 	tk := context.GetTokenFromContext(r)
// 	if proj.DomainID != tk.DomainID {
// 		response.Failed(w, exception.NewForbidden("the project not belong to your"))
// 		return
// 	}

// 	if err := global.Store.AddUsersToProject(tk.AccessToken, pid, uids...); err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	response.Success(w, http.StatusCreated, nil)
// 	return
// }

// // RemoveUsersFromProject use to create domain
// func RemoveUsersFromProject(w http.ResponseWriter, r *http.Request) {
// 	ps := context.GetParamsFromContext(r)
// 	pid := ps.ByName("pid")

// 	iter, err := request.CheckArrayBody(r)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	uids := make([]string, 0)
// 	for iter.ReadArray() {
// 		if uid := iter.ReadString(); uid != "" {
// 			uids = append(uids, uid)
// 		}
// 	}
// 	if iter.Error != nil {
// 		response.Failed(w, exception.NewBadRequest("get userid from body array error, %s", iter.Error))
// 		return
// 	}

// 	if len(uids) == 0 {
// 		response.Failed(w, exception.NewBadRequest(`body array hase no userid`))
// 		return
// 	}

// 	proj, err := global.Store.GetProject(pid)
// 	if err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	// check permisson
// 	tk := context.GetTokenFromContext(r)
// 	if proj.DomainID != tk.DomainID {
// 		response.Failed(w, exception.NewForbidden("the project not belong to your"))
// 		return
// 	}

// 	if err := global.Store.RemoveUsersFromProject(tk.AccessToken, pid, uids...); err != nil {
// 		response.Failed(w, err)
// 		return
// 	}

// 	response.Success(w, http.StatusCreated, nil)
// 	return
// }
