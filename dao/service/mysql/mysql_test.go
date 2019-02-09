package mysql_test

import (
	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/dao/service/mysql"
	"github.com/defineiot/keyauth/internal/conf/mock"
)

func newTestStore() service.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	log, err := conf.GetLogger()
	if err != nil {
		panic(err)
	}

	opt := &dao.Options{DB: db, LOG: log}
	store, err := mysql.NewServiceStore(opt)
	if err != nil {
		panic(err)
	}

	return store
}

type serviceSuit struct {
	roleID          string
	registryVersion string

	store    service.Store
	svr      *models.Service
	features []*models.Feature
}

func (s *serviceSuit) TearDown() {
	s.store.Close()
}

func (s *serviceSuit) SetUp() {
	s.roleID = "unit-test-role-id"
	s.registryVersion = "unit-test-for-instance-registry"

	s.svr = &models.Service{
		Name:        "unit-test-service-name01",
		Description: "just for unit test",
		Type:        models.Public,
	}

	f1 := &models.Feature{
		Name:         "feature01",
		Tag:          "POST",
		HTTPEndpoint: "/features/01",
		Description:  "only for unit test",
	}
	f2 := &models.Feature{
		Name:         "feature02",
		Tag:          "POST",
		HTTPEndpoint: "/features/02",
		Description:  "only for unit test",
	}
	s.features = []*models.Feature{f1, f2}

	s.store = newTestStore()

}
