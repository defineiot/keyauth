package controller

import (
	"openauth/api/config"
	"openauth/api/controller/domain"
	"openauth/pkg/domain/mysql"
)

// InitAllController use to initial all controllers
func InitAllController(conf *config.Config) error {
	db, err := conf.GetDBConn()
	if err != nil {
		return err
	}
	logger, err := conf.GetLogger()
	if err != nil {
		return err
	}

	dm, err := mysql.NewDomainManager(db)
	if err != nil {
		return err
	}

	if err := domain.InitController(db, logger, dm); err != nil {
		return err
	}

	return nil
}
