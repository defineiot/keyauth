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
	"openauth/pkg/service"
	"openauth/pkg/user"

	appmysql "openauth/store/application/mysql"
	domainmysql "openauth/store/domain/mysql"
	projectmysql "openauth/store/project/mysql"
	svrmysql "openauth/store/service/mysql"
	tokenmysql "openauth/store/token/mysql"
	usermysql "openauth/store/user/mysql"
)

var (
	domainsrv  *domain.Controller
	projectsrv *project.Controller
	usersrv    *user.Controller
	appsrv     *application.Controller
	authsrc    *oauth2.Controller
	svr        *service.Controller
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
		projectstr, err := projectmysql.NewProjectStore(db)
		if err != nil {
			storeErr = append(storeErr, err)
		}
		userstr, err := usermysql.NewUserStore(db, conf.APP.Key, log)
		if err != nil {
			storeErr = append(storeErr, err)
		}
		appstr, err := appmysql.NewAppStore(db)
		if err != nil {
			storeErr = append(storeErr, err)
		}
		tokenstr, err := tokenmysql.NewTokenStore(db)
		if err != nil {
			storeErr = append(storeErr, err)
		}
		svrstr, err := svrmysql.NewServiceStore(db)
		if err != nil {
			storeErr = append(storeErr, err)
		}

		domainsrv = domain.NewController(domainstr, log)
		projectsrv = project.NewController(log, domainstr, projectstr, userstr)
		usersrv = user.NewController(log, userstr, domainstr, projectstr)
		appsrv = application.NewController(log, appstr, userstr)
		authsrc = oauth2.NewController(tokenstr, userstr, domainstr, appstr, log, conf.Token.Type, conf.Token.ExpiresIn)
		svr = service.NewController(svrstr, log)
	})
	if len(storeErr) != 0 {
		return exception.NewInternalServerError("get store error, %s", storeErr)
	}

	return nil
}
