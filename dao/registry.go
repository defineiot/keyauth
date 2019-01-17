package dao

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao/application"
	"github.com/defineiot/keyauth/dao/department"
	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/dao/project"
	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/dao/token"
	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/dao/verifycode"
	"github.com/defineiot/keyauth/internal/logger"
)

// DaoFactory 默认的DAO层
var DaoFactory *factory

// factory 所有对象的DAO层
type factory struct {
	app  registryAPP
	dep  registryDEP
	dom  registryDomain
	pro  registryProject
	role registryRole
	svr  registryService
	tk   registryToken
	usr  registryUser
	vf   registryVerifyCode
}

// Options 创建DAO的选择参数
type Options struct {
	DB        *sql.DB
	LOG       logger.Logger
	ConfigMap map[string]string
}

// RegistryAPP 创建对象DAO层的方法
type registryAPP func(opt *Options) (application.Store, error)
type registryDEP func(opt *Options) (department.Store, error)
type registryDomain func(opt *Options) (domain.Store, error)
type registryProject func(opt *Options) (project.Store, error)
type registryRole func(opt *Options) (role.Store, error)
type registryService func(opt *Options) (service.Store, error)
type registryToken func(opt *Options) (token.Store, error)
type registryUser func(opt *Options) (user.Store, error)
type registryVerifyCode func(opt *Options) (verifycode.Store, error)

// Registe 注册一个对象的DAO层
func Registe(registryFunc interface{}) {
	switch v := registryFunc.(type) {
	case registryAPP:
		DaoFactory.app = v
	case registryDEP:
		DaoFactory.dep = v
	case registryDomain:
		DaoFactory.dom = v
	case registryProject:
		DaoFactory.pro = v
	case registryRole:
		DaoFactory.role = v
	case registryService:
		DaoFactory.svr = v
	case registryToken:
		DaoFactory.tk = v
	case registryUser:
		DaoFactory.usr = v
	}
}

// Dao 所有对象的DAO层
type Dao struct {
	Application application.Store
	Department  department.Store
	Domain      domain.Store
	Project     project.Store
	Role        role.Store
	Service     service.Store
	Token       token.Store
	User        user.Store
	VerifyCode  verifycode.Store
}

// Init 初始化dao层
func (f *factory) Init(opt *Options) (*Dao, error) {
	dao := new(Dao)

	app, err := f.app(opt)
	if err != nil {
		return nil, err
	}
	dep, err := f.dep(opt)
	if err != nil {
		return nil, err
	}
	dom, err := f.dom(opt)
	if err != nil {
		return nil, err
	}
	pro, err := f.pro(opt)
	if err != nil {
		return nil, err
	}
	role, err := f.role(opt)
	if err != nil {
		return nil, err
	}
	svr, err := f.svr(opt)
	if err != nil {
		return nil, err
	}
	tk, err := f.tk(opt)
	if err != nil {
		return nil, err
	}
	usr, err := f.usr(opt)
	if err != nil {
		return nil, err
	}
	vf, err := f.vf(opt)
	if err != nil {
		return nil, err
	}

	dao.Application = app
	dao.Department = dep
	dao.Domain = dom
	dao.Project = pro
	dao.Role = role
	dao.Service = svr
	dao.Token = tk
	dao.User = usr
	dao.VerifyCode = vf

	return nil, nil
}
