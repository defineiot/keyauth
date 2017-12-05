package handler

import (
	"sync"

	"openauth/api/config"
	"openauth/pkg/domain"
	"openauth/pkg/project"
	"openauth/pkg/user"

	domainstorage "openauth/storage/domain/mysql"
	projectstorage "openauth/storage/project/mysql"
	userstorage "openauth/storage/user/mysql"
)

var (
	domainctl  *domain.Controller
	projectctl *project.Controller
	userctl    *user.Controller
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
	us := userstorage.NewUserStorage(db, conf.APP.Key, logger)

	once.Do(func() {
		domainctl = domain.NewController(logger, ds)
		projectctl = project.NewController(logger, ds, ps, us)
		userctl = user.NewController(logger, us, ds, ps)
		logger.Debugf("domain controller: %v", domainctl)
		logger.Debugf("project controoler: %v", projectctl)
		logger.Debugf("user controller: %v", userctl)
	})

	return nil
}
