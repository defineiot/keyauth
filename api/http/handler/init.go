package handler

import (
	"openauth/api/config"
	"openauth/pkg/domain"
	"openauth/pkg/project"
	"sync"

	domainstorage "openauth/storage/domain/mysql"
	projectstorage "openauth/storage/project/mysql"
)

var (
	domainctl  *domain.Controller
	projectctl *project.Controller
	once       sync.Once
)

// InitController use to initial all controllers
func InitController(conf *config.Config) error {
	db, err := conf.GetDBConn()
	if err != nil {
		return err
	}
	logger, err := conf.GetLogger()
	if err != nil {
		return err
	}

	ps := projectstorage.NewProjectStroage(db)
	ds := domainstorage.NewDomainStorage(db)

	once.Do(func() {
		domainctl = domain.NewController(logger, ds)
		projectctl = project.NewController(logger, ds, ps)
	})

	return nil
}
