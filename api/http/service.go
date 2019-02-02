package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/logger"
	"github.com/defineiot/keyauth/store"
	"github.com/defineiot/keyauth/version"
)

var stopSignal = make(chan bool, 1)

// Service is auth service
type Service struct {
	http  *http.Server
	conf  *conf.Config
	log   logger.Logger
	store *store.Store

	v1endpoints map[string]map[string]string
	description string
}

// NewService  new http service
func NewService(config *conf.Config) (*Service, error) {
	// 检查服务配置是否正确
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("配置校验失败, %s", err)
	}

	// 初始化路由
	n := negroni.New()

	// 添加路由
	r := router.NewRouter()
	r.SetURLPrefix("/keyauth/v1")
	RouteToV1(r)
	n.UseHandler(r.Router)

	// 设置路由中间件
	corsM := cors.AllowAll()
	recoverM := negroni.NewRecovery()
	accessL := negroni.NewLogger()
	n.Use(corsM)
	n.Use(accessL)
	n.Use(recoverM)

	// 生成HTTP服务对象
	addr := fmt.Sprintf("%s:%s", config.APP.Host, config.APP.Port)
	httpSvr := &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		Addr:              addr,
		Handler:           n,
	}

	// 获取服务日志对象
	logger, err := config.GetLogger()
	if err != nil {
		return nil, err
	}

	return &Service{
		http: httpSvr,
		conf: config,
		log:  logger,

		v1endpoints: r.GetEndpoints(),
		description: `微服务权限管理中心, 提供用户的管理, 认证, 鉴权, 服务发现等功能`,
	}, nil
}

// Start use to start openauth http service
func (s *Service) Start() error {
	// initial global variables
	if err := s.init(); err != nil {
		return err
	}
	s.log.Debug("initial global variables success")

	// registe service
	if s.conf.Etcd.EnableRegisteFeatures {
		if err := s.registryService(); err != nil {
			s.log.Error("registe keyauth service failed, %s", err)
		} else {
			s.log.Debug("registe keyauth service features success")
		}
	}

	// start http service
	if err := s.start(); err != nil {
		return err
	}

	return nil
}

func (s *Service) init() error {
	// 初始化全局变量
	store, err := store.NewStore(s.conf)
	if err != nil {
		return err
	}
	s.store = store

	store.SetCache(cache.Newmemcache(100000), time.Minute*5)
	global.Store = store
	global.Conf = s.conf
	global.Log = s.log

	return nil
}

func (s *Service) registryService() error {
	if len(s.v1endpoints) == 0 {
		return errors.New("there is no feature to registry")
	}

	var (
		keysvr *service.Service
		err    error
	)

	name := s.conf.APP.Name
	keysvr, err = global.Store.GetServiceByName(name)
	if err != nil {
		if _, ok := err.(*exception.NotFound); ok {
			keysvr = &service.Service{
				Name:        name,
				Description: "微服务权限管理中心",
				Type:        service.Public,
				Enabled:     true,
			}
			if err = global.Store.CreateService(keysvr); err != nil {
				return err
			}
		}
	}

	if keysvr == nil || keysvr.ID == "" {
		return exception.NewInternalServerError("服务自创建异常, 创建失败")
	}

	features := []*service.Feature{}
	for method, v := range s.v1endpoints {
		for feature, ep := range v {
			f := &service.Feature{Name: feature, HTTPEndpoint: ep, Tag: method}
			features = append(features, f)
		}
	}

	if err := global.Store.RegistryServiceFeatures(keysvr.ID, version.GIT_TAG,
		features...); err != nil {
		return err
	}

	byteF, err := json.Marshal(features)
	if err != nil {
		global.Log.Info("marshal features to bytes error, %s", err)
	}
	global.Log.Info("registry keyauth service features success, version: %s, features: %s",
		version.GIT_TAG, string(byteF))

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
		s.log.Info("receive signal '%v', start graceful shutdown", sg)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.http.Shutdown(ctx); err != nil {
			s.log.Error("graceful shutdown timeout, force exit")
		}
		os.Exit(1)
	}()

	s.log.Info("starting keyauth service at %s", s.http.Addr)
	if err := s.http.ListenAndServe(); err != nil {
		return fmt.Errorf("start keyauth error, %s", err.Error())
	}
	return nil
}
