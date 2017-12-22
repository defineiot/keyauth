package handler

import (
	"openauth/api/exception"
	"sync"

	"openauth/api/config"
	"openauth/api/logger"
	"openauth/pkg/application"
	"openauth/pkg/domain"
	"openauth/pkg/oauth2"
	"openauth/pkg/project"
	"openauth/pkg/user"

	appmysql "openauth/store/application/mysql"
	domainmysql "openauth/store/domain/mysql"
	projectmysql "openauth/store/project/mysql"
	tokenmysql "openauth/store/token/mysql"
	usermysql "openauth/store/user/mysql"
)

var (
	domainsrv  *domain.Controller
	projectsrv *project.Controller
	usersrv    *user.Controller
	appsrv     *application.Controller
	authsrc    *oauth2.Controller
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

	storeErr := []error{}
	once.Do(func() {
		domainstr, err := domainmysql.NewDomainStore(db)
		if err != nil {
			storeErr = append(storeErr, err)
		}
		projectstr := projectmysql.NewProjectStorage(db)
		userstr := usermysql.NewUserStorage(db, conf.APP.Key, log)
		appstr := appmysql.NewApplicationStorage(db)
		tokenstr := tokenmysql.NewTokenStorage(db)

		domain.InitController(domainstr, log)
		project.InitController(log, domainstr, projectstr, userstr)
		user.InitController(log, userstr, domainstr, projectstr)
		application.InitController(log, appstr, userstr)
		oauth2.InitController(tokenstr, userstr, domainstr, appstr, log, conf.Token.Type, conf.Token.ExpiresIn)
	})
	if len(storeErr) != 0 {
		return exception.NewInternalServerError("get store error, %s", storeErr)
	}

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
	authsrc, err = oauth2.GetController()
	if err != nil {
		return err
	}

	return nil
}
