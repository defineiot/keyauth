package store

import (
	"errors"
	"fmt"

	"github.com/defineiot/keyauth/dao/application"
	"github.com/defineiot/keyauth/dao/department"
	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/dao/user"
)

const (
	systemAdminRoleName   string = "system_admin"
	domainAdminRoleName          = "domain_admin"
	memberUserRoleName           = "member"
	adminDomainName              = "admin_domain"
	adminDepartmentName          = "admin_department"
	defaultDepartmentName        = "default_department"
)

// InitAdmin 初始化系统管理员相关信息
func (s *Store) InitAdmin(username, password string) error {
	ok, err := s.dao.Domain.CheckDomainIsExistByName(adminDomainName)
	if err != nil {
		return err
	}

	if ok {
		return errors.New("系统管理员已经初始化完成")
	}

	fmt.Println("[INIT] 开始初始化 系统需要的角色 ...")
	if err := s.initRoles(); err != nil {
		return err
	}

	fmt.Println("[INIT] 开始初始化 系统管理员账户 ...")
	if err := s.initAdminUser(username, password); err != nil {
		return err
	}

	fmt.Println("[INIT] 开始初始化 系统管理员应用 ...")
	if err := s.initAdminAPPs(username); err != nil {
		return err
	}
	fmt.Println("[INIT] 系统管理员初始化完成")

	fmt.Println("")
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
	fmt.Printf("[INIT] 创建系统管理员角色成功: %s\n", systemAdmin.Name)

	if err := s.dao.Role.CreateRole(domainAdmin); err != nil {
		return err
	}
	fmt.Printf("[INIT] 创建租户管理员成功: %s\n", domainAdmin.Name)

	if err := s.dao.Role.CreateRole(common); err != nil {
		return err
	}
	fmt.Printf("[INIT] 创建普通成员角色成功: %s\n", common.Name)

	fmt.Println("")
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
	fmt.Printf("[INIT] 创建系统管理员部门成功: %s\n", adminDep.Name)

	if err := s.dao.Department.CreateDepartment(defaultDep); err != nil {
		return err
	}
	fmt.Printf("[INIT] 创建系统管理员默认部门成功: %s\n", defaultDep.Name)

	adminUser := &user.User{
		Account:    username,
		Password:   &user.Password{Password: s.hmacHash(password)},
		Domain:     adminDomain,
		Department: adminDep,
	}

	if err := s.dao.User.CreateUser(adminUser); err != nil {
		return err
	}
	fmt.Printf("[INIT] 创建系统管理员成功: %s\n", adminUser.Account)

	sysrole, err := s.dao.Role.GetRoleByName(systemAdminRoleName)
	if err != nil {
		return err
	}
	if err := s.dao.User.BindRole(adminUser.Domain.ID, adminUser.ID, sysrole.ID); err != nil {
		return err
	}
	fmt.Println("[INIT] 绑定系统管理员角色成功")

	domainrole, err := s.dao.Role.GetRoleByName(domainAdminRoleName)
	if err != nil {
		return err
	}
	if err := s.dao.User.BindRole(adminUser.Domain.ID, adminUser.ID, domainrole.ID); err != nil {
		return err
	}
	fmt.Println("[INIT] 绑定租户管理员角色成功")

	fmt.Println("")
	return nil
}

func (s *Store) initAdminAPPs(account string) error {
	u, err := s.dao.User.GetUser(user.Account, account)
	if err != nil {
		return err
	}

	web := &application.Application{
		Name:        "web_app",
		UserID:      u.ID,
		Description: "用于web端服务使用",
	}
	android := &application.Application{
		Name:        "android_app",
		UserID:      u.ID,
		Description: "用于构建安卓应用时使用",
	}
	ios := &application.Application{
		Name:        "ios_app",
		UserID:      u.ID,
		Description: "用于构建IOS应用时使用",
	}
	sdk := &application.Application{
		Name:        "sdk_app",
		UserID:      u.ID,
		Description: "用户构建SDK端时使用",
	}

	if err := s.dao.Application.CreateApplication(web); err != nil {
		return err
	}
	fmt.Printf("[INIT] 创建Web端应用应用成功: client_id -> %s, client_secret -> %s\n",
		web.ClientID, web.ClientSecret)

	if err := s.dao.Application.CreateApplication(android); err != nil {
		return err
	}
	fmt.Printf("[INIT] 创建安卓端应用应用成功: client_id -> %s, client_secret -> %s\n",
		android.ClientID, android.ClientSecret)

	if err := s.dao.Application.CreateApplication(ios); err != nil {
		return err
	}
	fmt.Printf("[INIT] 创建IOS端应用应用成功: client_id -> %s, client_secret -> %s\n",
		ios.ClientID, ios.ClientSecret)

	if err := s.dao.Application.CreateApplication(sdk); err != nil {
		return err
	}
	fmt.Printf("[INIT] 创建SDK端应用应用成功: client_id -> %s, client_secret -> %s\n",
		sdk.ClientID, sdk.ClientSecret)

	fmt.Println("")
	return nil
}
