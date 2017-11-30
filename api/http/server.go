package http

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"openauth/api/config"
	"openauth/api/http/router"
)

var stopSignal = make(chan bool, 1)

// Service is gateway service
type Service struct {
	conf   *config.Config
	logger *logrus.Logger
}

// NewService use to new an gateway service
func NewService(conf *config.Config) (*Service, error) {
	err := conf.Validate()
	if err != nil {
		return nil, fmt.Errorf("config validate failed, %s", err)
	}

	logger, err := conf.GetLogger()
	if err != nil {
		return nil, err
	}

	return &Service{conf: conf, logger: logger}, nil
}

// Start use to start openauth http service
func (s *Service) Start() error {
	n := negroni.New()

	r := router.NewRouter()
	r.RouteToDomain()

	// includes some default middlewares
	corsM := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowedHeaders:   []string{"Origin", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})
	recoverM := negroni.NewRecovery()
	n.Use(corsM)
	n.Use(recoverM)
	s.logger.Info("loading http middleware success")

	// config router
	n.UseHandler(r.Router)
	s.logger.Info("loading router success")

	// run http service
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sg := <-ch
		s.logger.Infof("receive signal '%v'", sg)
		os.Exit(1)
	}()

	addr := fmt.Sprintf("%s:%s", s.conf.APP.Host, s.conf.APP.Port)
	s.logger.Infof("starting openauth service at %s", addr)
	if err := http.ListenAndServe(addr, n); err != nil {
		return fmt.Errorf("start openauth error, %s", err.Error())
	}
	return nil
}
