package store

import (
	"github.com/defineiot/keyauth/dao/application"
	"github.com/defineiot/keyauth/dao/department"
	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/dao/user"
)

const (
	systemAdminRoleName   = "system_admin"
	domainAdminRoleName   = "domain_admin"
	memberUserRoleName    = "member"
	adminDomainName       = "admin_domain"
	adminDepartmentName   = "admin_department"
	defaultDepartmentName = "default_department"
)

// InitAdmin 初始化系统管理员相关信息
func (s *Store) InitAdmin(username, password string) error {
	if err := s.initRoles(); err != nil {
		return err
	}

	if err := s.initAdminUser(username, password); err != nil {
		return err
	}

	if err := s.initAdminAPPs(username); err != nil {
		return err
	}

	return nil
}

// 初始化3个角色:
// 系统管理员: system_admin
// 域管理员:   domain_admin
// 普通用户:   member
func (s *Store) initRoles() error {
	systemAdmin := &role.Role{
		Name:        systemAdminRoleName,
		Description: "系统管理员",
	}
	domainAdmin := &role.Role{
		Name:        domainAdminRoleName,
		Description: "域管理员/公司管理员/组织管理员",
	}
	common := &role.Role{
		Name:        memberUserRoleName,
		Description: "成员用户",
	}

	if err := s.dao.Role.CreateRole(systemAdmin); err != nil {
		return err
	}

	if err := s.dao.Role.CreateRole(domainAdmin); err != nil {
		return err
	}

	if err := s.dao.Role.CreateRole(common); err != nil {
		return err
	}

	return nil
}

// 创建管理员的账号和域
func (s *Store) initAdminUser(username, password string) error {
	adminDomain := &domain.Domain{
		Type:        domain.Enterprise,
		Name:        adminDomainName,
		DisplayName: "系统管理员域空间",
		Description: "系统管理员域空间",
	}

	if err := s.dao.Domain.CreateDomain(adminDomain); err != nil {
		return err
	}

	adminDep := &department.Department{
		Name:     adminDepartmentName,
		DomainID: adminDomain.ID,
	}

	defaultDep := &department.Department{
		Name:     defaultDepartmentName,
		DomainID: adminDomain.ID,
	}

	if err := s.dao.Department.CreateDepartment(adminDep); err != nil {
		return err
	}

	if err := s.dao.Department.CreateDepartment(defaultDep); err != nil {
		return err
	}

	adminUser := &user.User{
		Account:    username,
		Password:   &user.Password{Password: s.hmacHash(password)},
		Domain:     adminDomain,
		Department: adminDep,
	}

	if err := s.dao.User.CreateUser(adminUser); err != nil {
		return err
	}

	sysrole, err := s.dao.Role.GetRoleByName(systemAdminRoleName)
	if err != nil {
		return err
	}
	if err := s.dao.User.BindRole(adminUser.Domain.ID, adminUser.ID, sysrole.ID); err != nil {
		return err
	}

	domainrole, err := s.dao.Role.GetRoleByName(domainAdminRoleName)
	if err != nil {
		return err
	}
	if err := s.dao.User.BindRole(adminUser.Domain.ID, adminUser.ID, domainrole.ID); err != nil {
		return err
	}

	return nil
}

func (s *Store) initAdminAPPs(account string) error {
	u, err := s.dao.User.GetUser(user.Account, account)
	if err != nil {
		return err
	}

	web := &application.Application{
		Name:   "web_app",
		UserID: u.ID,
	}
	android := &application.Application{
		Name:   "android_app",
		UserID: u.ID,
	}
	ios := &application.Application{
		Name:   "ios_app",
		UserID: u.ID,
	}
	sdk := &application.Application{
		Name:   "sdk_app",
		UserID: u.ID,
	}

	if err := s.dao.Application.CreateApplication(web); err != nil {
		return err
	}

	if err := s.dao.Application.CreateApplication(android); err != nil {
		return err
	}

	if err := s.dao.Application.CreateApplication(ios); err != nil {
		return err
	}

	if err := s.dao.Application.CreateApplication(sdk); err != nil {
		return err
	}

	return nil
}
