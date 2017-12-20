package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"openauth/api/config"
	"openauth/api/http/handler"
	"openauth/api/http/router"
	"openauth/api/logger"
)

var stopSignal = make(chan bool, 1)

// Service is gateway service
type Service struct {
	conf   *config.Config
	logger logger.OpenAuthLogger
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
	r.RouteToV1()

	// includes some default middlewares
	corsM := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowedHeaders:   []string{"Origin", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})
	recoverM := negroni.NewRecovery()
	accessL := negroni.NewLogger()
	n.Use(corsM)
	n.Use(accessL)
	n.Use(recoverM)
	s.logger.Info("loading http middleware success")

	// initial controller
	if err := handler.InitController(s.conf); err != nil {
		return err
	}

	// config router
	n.UseHandler(r.Router)
	s.logger.Info("loading router success")

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

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sg := <-ch
		s.logger.Infof("receive signal '%v', start graceful shutdown", sg)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			s.logger.Error("graceful shutdown timeout, force exit")
		}
		os.Exit(1)
	}()

	s.logger.Infof("starting openauth service at %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("start openauth error, %s", err.Error())
	}

	return nil
}
