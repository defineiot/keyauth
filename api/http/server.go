package http

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"openauth/api/config"
	"openauth/api/http/router"
)

var stopSignal = make(chan bool, 1)

// Service is gateway service
type Service struct {
	conf *config.Config
}

// NewService use to new an gateway service
func NewService(conf *config.Config) (*Service, error) {
	err := conf.Validate()
	if err != nil {
		return nil, fmt.Errorf("config validate failed, %s", err)
	}

	return &Service{conf: conf}, nil
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

	// config router
	n.UseHandler(r.Router)

	// run http service
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		fmt.Printf("receive signal '%v'", s)
		os.Exit(1)
	}()

	addr := fmt.Sprintf("%s:%s", s.conf.APP.Host, s.conf.APP.Port)
	fmt.Printf("starting openauth service at %s:%s ...\n", s.conf.APP.Host, s.conf.APP.Port)
	if err := http.ListenAndServe(addr, n); err != nil {
		return fmt.Errorf("start openauth error, %s", err.Error())
	}
	return nil
}
