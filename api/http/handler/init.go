package handler

import (
	"sync"

	"openauth/api/config"
	"openauth/api/logger"
	"openauth/pkg/application"
	"openauth/pkg/domain"
	"openauth/pkg/project"
	"openauth/pkg/user"

	appmysql "openauth/storage/application/mysql"
	domainmysql "openauth/storage/domain/mysql"
	projectmysql "openauth/storage/project/mysql"
	usermysql "openauth/storage/user/mysql"
)

var (
	domainsrv  *domain.Controller
	projectsrv *project.Controller
	usersrv    *user.Controller
	appsrv     *application.Controller
	log        logger.OpenAuthLogger
	once       sync.Once
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
		domainstr := domainmysql.NewDomainService(db)
		projectstr := projectmysql.NewProjectService(db)
		userstr := usermysql.NewUserService(db, conf.APP.Key, log)
		appstr := appmysql.NewApplicationService(db)

		domain.InitController(domainstr, log)
		project.InitController(log, domainstr, projectstr, userstr)
		user.InitController(log, userstr, domainstr, projectstr)
		application.InitController(log, appstr, userstr)
	})

	domainsrv, err = domain.GetController()
	if err != nil {
		return err
	}
	projectsrv, err = project.GetController()
	if err != nil {
		return err
	}
	usersrv, err = user.GetController()
	if err != nil {
		return err
	}
	appsrv, err = application.GetController()
	if err != nil {
		return err
	}

	return nil
}
