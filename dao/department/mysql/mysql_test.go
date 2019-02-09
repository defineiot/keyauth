package mysql_test

import (
	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/department"
	"github.com/defineiot/keyauth/dao/department/mysql"
	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/internal/conf/mock"
)

func newTestStore() department.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}
	logger, err := conf.GetLogger()
	if err != nil {
		panic(err)
	}

	opt := &dao.Options{DB: db, LOG: logger}
	store, err := mysql.NewDepartmentStore(opt)
	if err != nil {
		panic(err)
	}

	return store
}

type departmentSuit struct {
	store department.Store
	l1    *models.Department
	l2    *models.Department
	l3    *models.Department
	l4    *models.Department
}

func (s *departmentSuit) TearDown() {

	s.store.Close()
}

func (s *departmentSuit) SetUp() {
	s.store = newTestStore()
	s.l1 = &models.Department{
		Name:     "Root_Department_l1",
		DomainID: "domain_unit_test_for_department",
	}
	s.l2 = &models.Department{
		Name:     "Root_Department_l2",
		DomainID: "domain_unit_test_for_department",
	}
	s.l3 = &models.Department{
		Name:     "Root_Department_l3",
		DomainID: "domain_unit_test_for_department",
	}
	s.l4 = &models.Department{
		Name:     "Root_Department_l4",
		DomainID: "domain_unit_test_for_department",
	}
}
