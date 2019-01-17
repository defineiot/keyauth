package dao

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/defineiot/keyauth/internal/logger"
)

// Factory 默认的DAO层
var Factory *factory

var once sync.Once

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

// Init 初始化dao层
func Init(opt *Options) (*Dao, error) {
	f := Factory
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

	return dao, nil
}

func init() {
	once.Do(func() {
		Factory = new(factory)
	})
	fmt.Println("init dao")
}
