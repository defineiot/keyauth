package handler

import (
	"sync"

	"openauth/api/config"
	"openauth/api/logger"
	"openauth/pkg/application"
	"openauth/pkg/domain"
	"openauth/pkg/project"
	"openauth/pkg/user"

	appmysql "openauth/pkg/application/mysql"
	domainmysql "openauth/pkg/domain/mysql"
	projectmysql "openauth/pkg/project/mysql"
	usermysql "openauth/pkg/user/mysql"
)

var (
	domainsrv  domain.Service
	projectsrv project.Service
	usersrv    user.Service
	appsrv     application.Service
	once       sync.Once
	log        logger.OpenAuthLogger
)

// InitController use to initial all controllers
func InitController(conf *config.Config) error {
	db, err := conf.GetDBConn()
	if err != nil {
		return err
	}
	log, err = conf.GetLogger()
	if err != nil {
		return err
	}

	once.Do(func() {
		domainsrv = domainmysql.NewDomainService(db)
		projectsrv = projectmysql.NewProjectService(db)
		usersrv = usermysql.NewUserService(db, conf.APP.Key, log)
		appsrv = appmysql.NewApplicationService(db)

		log.Debugf("domain service: %v", domainsrv)
		log.Debugf("project service: %v", projectsrv)
		log.Debugf("user service: %v", usersrv)
		log.Debugf("application service: %v", appsrv)
	})

	return nil
}
