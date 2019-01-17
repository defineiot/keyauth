package dao

import (
	"fmt"

	"github.com/defineiot/keyauth/dao/application"
	"github.com/defineiot/keyauth/dao/department"
	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/dao/project"
	"github.com/defineiot/keyauth/dao/role"
	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/dao/token"
	"github.com/defineiot/keyauth/dao/user"
	"github.com/defineiot/keyauth/dao/verifycode"
)

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

// Registe 注册一个对象的DAO层
func Registe(registryFunc interface{}) {
	switch v := registryFunc.(type) {
	case func(*Options) (application.Store, error):
		Factory.app = v
	case func(*Options) (department.Store, error):
		Factory.dep = v
	case func(*Options) (domain.Store, error):
		Factory.dom = v
	case func(*Options) (project.Store, error):
		Factory.pro = v
	case func(*Options) (role.Store, error):
		Factory.role = v
	case func(*Options) (service.Store, error):
		Factory.svr = v
	case func(*Options) (token.Store, error):
		Factory.tk = v
	case func(*Options) (user.Store, error):
		Factory.usr = v
	case func(*Options) (verifycode.Store, error):
		Factory.vf = v
	default:
		fmt.Printf("unknow registry func: %v\n", v)
		panic("unknow registry func")
	}
}

// RegistryAPP 创建对象DAO层的方法
type registryAPP func(*Options) (application.Store, error)
type registryDEP func(*Options) (department.Store, error)
type registryDomain func(*Options) (domain.Store, error)
type registryProject func(*Options) (project.Store, error)
type registryRole func(*Options) (role.Store, error)
type registryService func(*Options) (service.Store, error)
type registryToken func(*Options) (token.Store, error)
type registryUser func(*Options) (user.Store, error)
type registryVerifyCode func(*Options) (verifycode.Store, error)
