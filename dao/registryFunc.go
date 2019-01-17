package dao

import (
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
// func Registe(registryFunc interface{}) {
// 	fmt.Println(reflect.TypeOf(registryFunc).String())
// 	switch v := registryFunc.(type) {
// 	case registryAPP:
// 		Factory.app = v
// 	case registryDEP:
// 		Factory.dep = v
// 	case registryDomain:
// 		Factory.dom = v
// 	case registryProject:
// 		Factory.pro = v
// 	case registryRole:
// 		Factory.role = v
// 	case registryService:
// 		Factory.svr = v
// 	case registryToken:
// 		Factory.tk = v
// 	case registryUser:
// 		Factory.usr = v
// 	default:
// 		fmt.Printf("unknow registry func: %v\n", v)
// 		panic("unknow registry func")
// 	}
// }

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

// RegistryAPP 创建对象DAO层的方法
func RegistryAPP(rfunc registryAPP) {
	Factory.app = rfunc
}

// RegistryDEP 创建对象DAO层的方法
func RegistryDEP(rfunc registryDEP) {
	Factory.dep = rfunc
}

// RegistryDomain 创建对象DAO层的方法
func RegistryDomain(rfunc registryDomain) {
	Factory.dom = rfunc
}

// RegistryProject 创建对象DAO层的方法
func RegistryProject(rfunc registryProject) {
	Factory.pro = rfunc
}

// RegistryRole 创建对象DAO层的方法
func RegistryRole(rfunc registryRole) {
	Factory.role = rfunc
}

// RegistryService 创建对象DAO层的方法
func RegistryService(rfunc registryService) {
	Factory.svr = rfunc
}

// RegistryToken 创建对象DAO层的方法
func RegistryToken(rfunc registryToken) {
	Factory.tk = rfunc
}

// RegistryUser 创建对象DAO层的方法
func RegistryUser(rfunc registryUser) {
	Factory.usr = rfunc
}

// RegistryVerifyCode 创建对象DAO层的方法
func RegistryVerifyCode(rfunc registryVerifyCode) {
	Factory.vf = rfunc
}
