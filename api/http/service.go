package http

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"github.com/defineiot/keyauth/api/global"
	"github.com/defineiot/keyauth/api/http/router"
	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/internal/cache"
	"github.com/defineiot/keyauth/internal/conf"
	"github.com/defineiot/keyauth/internal/log"
	"github.com/defineiot/keyauth/store"

	appalication "github.com/defineiot/keyauth/dao/application/mysql"
	client "github.com/defineiot/keyauth/dao/client/mysql"
	domain "github.com/defineiot/keyauth/dao/domain/mysql"
	project "github.com/defineiot/keyauth/dao/project/mysql"
	role "github.com/defineiot/keyauth/dao/role/mysql"
	svr "github.com/defineiot/keyauth/dao/service/mysql"
	token "github.com/defineiot/keyauth/dao/token/mysql"
	user "github.com/defineiot/keyauth/dao/user/mysql"
)

var stopSignal = make(chan bool, 1)

// Service is gateway service
type Service struct {
	http        *http.Server
	conf        *conf.Config
	log         log.IOTAuthLogger
	v1endpoints map[string]map[string]string
	description string
}

// NewService use to new an gateway service
func NewService(config *conf.Config) (*Service, error) {
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("config validate failed, %s", err)
	}

	logger, err := config.GetLogger()
	if err != nil {
		return nil, err
	}

	desc := `iot平台权限管理服务, 提供用户的管理, 认证, 鉴权, 服务发现等功能`

	return &Service{conf: config, log: logger, description: desc}, nil
}

// Start use to start openauth http service
func (s *Service) Start() error {
	// initial global variables
	if err := s.initGlobal(); err != nil {
		return err
	}
	s.log.Debug("initial global variables success")

	// prepare http service
	if err := s.prepare(); err != nil {
		return err
	}
	s.log.Debug("initial http service success")

	// registe service
	if err := s.registryService(); err != nil {
		return err
	}
	s.log.Debug("registry github.com/defineiot/keyauth service features success")

	// start http service
	if err := s.start(); err != nil {
		return err
	}

	return nil
}

func (s *Service) BootStrap() error {
	// initial global variables
	if err := s.initGlobal(); err != nil {
		return err
	}
	s.log.Debug("initial global variables success")

	// prepare http service
	if err := s.prepare(); err != nil {
		return err
	}
	s.log.Debug("initial http service success")

	// registe service
	if err := s.registryService(); err != nil {
		return err
	}
	s.log.Debug("registry github.com/defineiot/keyauth service features success")

	// initial roles
	if err := s.initialRoles(); err != nil {
		return err
	}
	s.log.Debug("initial default roles for system success")

	// initial sysadmin
	admin := s.conf.Admin
	if err := s.initialSysAdmin(admin.Domain, admin.DomainDisplay, admin.UserName, admin.Password); err != nil {
		return err
	}
	s.log.Debug("initial supser admin success")

	// start http service
	if err := s.start(); err != nil {
		return err
	}

	return nil
}

func (s *Service) initGlobal() error {
	db, err := s.conf.GetDBConn()
	if err != nil {
		return err
	}
	dom, err := domain.NewDomainStore(db)
	if err != nil {
		return err
	}
	pro, err := project.NewProjectStore(db)
	if err != nil {
		return err
	}
	usr, err := user.NewUserStore(db, s.conf.APP.Key, s.log)
	if err != nil {
		return err
	}
	app, err := appalication.NewAppStore(db)
	if err != nil {
		return err
	}
	token, err := token.NewTokenStore(db)
	if err != nil {
		return err
	}
	cli, err := client.NewClientStore(db)
	if err != nil {
		return err
	}
	mysqlService, err := svr.NewServiceStore(db, s.log)
	if err != nil {
		return err
	}
	rl, err := role.NewRoleStore(db, s.log)
	if err != nil {
		return err
	}

	opts := store.Options{Domain: dom, Project: pro, User: usr, Log: s.log, App: app, Conf: s.conf, Token: token, Client: cli, Service: mysqlService, Role: rl}

	store := store.NewStore(&opts)
	store.SetCache(cache.Newmemcache(100000), time.Minute*5)

	global.Store = store
	global.Conf = s.conf
	global.Log = s.log

	return nil
}

func (s *Service) prepare() error {
	n := negroni.New()
	r := router.NewRouter()
	RouteToV1(r)

	// includes some default middlewares
	corsM := cors.AllowAll()
	corsM.Log = stdlog.New(os.Stdout, "Info: ", stdlog.Ltime|stdlog.Lshortfile)
	recoverM := negroni.NewRecovery()
	accessL := negroni.NewLogger()
	n.Use(corsM)
	n.Use(accessL)
	n.Use(recoverM)
	s.log.Info("loading http middleware success")

	// config router
	n.UseHandler(r.Router)
	s.log.Info("loading router success")

	// run http service
	addr := fmt.Sprintf("%s:%s", s.conf.APP.Host, s.conf.APP.Port)
	srv := &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		Addr:              addr,
		Handler:           n,
	}

	s.http = srv
	s.v1endpoints = r.GetEndpoints()

	return nil
}

func (s *Service) registryService() error {
	if len(s.v1endpoints) == 0 {
		return errors.New("there is no feature to registry")
	}

	name := s.conf.APP.Name
	ok, err := global.Store.CheckService(name)
	if err != nil {
		return err
	}
	if !ok {
		if _, err := global.Store.CreateService(name, s.description); err != nil {
			return err
		}
	}

	features := []service.Feature{}
	for method, v := range s.v1endpoints {
		for feature, ep := range v {
			f := service.Feature{Name: feature, Endpoint: ep, Method: method}
			features = append(features, f)
		}
	}

	if err := global.Store.RegistryServiceFeatures(name, features...); err != nil {
		return err
	}

	return nil
}

// 初始化3个角色:
// 系统管理员: system_admin
// 域管理员:   domain_admin
// 普通用户:   member
func (s *Service) initialRoles() error {
	gs := global.Store

	ok, err := gs.CheckRoleExist("system_admin")
	if err != nil {
		return err
	}
	if !ok {
		if _, err := gs.CreateRole("system_admin", "系统超级管理员"); err != nil {
			return err
		}
	}

	ok, err = gs.CheckRoleExist("domain_admin")
	if err != nil {
		return err
	}
	if !ok {
		if _, err := gs.CreateRole("domain_admin", "域管理员"); err != nil {
			return err
		}
	}

	ok, err = gs.CheckRoleExist("member")
	if err != nil {
		return err
	}
	if !ok {
		if _, err := gs.CreateRole("member", "普通用户"); err != nil {
			return err
		}
	}

	dfs, err := gs.ListDomainFeatures()
	if err != nil {
		return err
	}
	mfs, err := gs.ListMemberFeatures()
	if err != nil {
		return err
	}

	dfids := []int64{}
	for _, df := range dfs {
		dfids = append(dfids, df.ID)
	}
	err = gs.AddFeaturesToRole("domain_admin", dfids...)
	if err != nil {
		s.log.Info(err)
	}

	mfids := []int64{}
	for _, mf := range mfs {
		mfids = append(mfids, mf.ID)
	}
	err = gs.AddFeaturesToRole("member", mfids...)
	if err != nil {
		s.log.Info(err)
	}

	return nil
}

func (s *Service) initialSysAdmin(domainName, domainDisplay, username, password string) error {
	gs := global.Store

	ok, err := gs.CheckDomainExistByName(domainName)
	if err != nil {
		return err
	}
	if !ok {
		dom, err := gs.CreateDomain(domainName, "超级管理员域", domainDisplay, true)
		if err != nil {
			return err
		}
		us, err := gs.CreateUser(dom.ID, username, password, true, 8760, 8760)
		if err != nil {
			return err
		}

		app, err := gs.CreateApplication(us.ID, "dashboard", "", "超级管理员的默认APP", "")
		if err != nil {
			return err
		}

		if err := gs.BindRole(dom.ID, us.ID, "system_admin"); err != nil {
			return err
		}

		s.log.Infof("super admin client_id: %s", app.Client.ID)
		s.log.Infof("super admin client_secret: %s", app.Client.Secret)
	}

	return nil
}

func (s *Service) start() error {
	if s.http == nil {
		return errors.New("no http service find, http service not prepare")
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sg := <-ch
		s.log.Infof("receive signal '%v', start graceful shutdown", sg)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.http.Shutdown(ctx); err != nil {
			s.log.Error("graceful shutdown timeout, force exit")
		}
		os.Exit(1)
	}()

	s.log.Infof("starting keyauth service at %s", s.http.Addr)
	if err := s.http.ListenAndServe(); err != nil {
		return fmt.Errorf("start keyauth error, %s", err.Error())
	}
	return nil
}
