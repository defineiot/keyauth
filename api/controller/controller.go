package controller

import (
	"openauth/api/config"
	"openauth/api/controller/domain"
	"openauth/api/controller/project"

	domsql "openauth/pkg/domain/mysql"
	prosql "openauth/pkg/project/mysql"
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

	dm := domsql.NewDomainManager(db)
	pm := prosql.NewProjectManager(db, dm)

	if err := domain.InitController(db, logger, dm); err != nil {
		return err
	}
	if err := project.InitController(db, logger, pm); err != nil {
		return err
	}

	return nil
}
