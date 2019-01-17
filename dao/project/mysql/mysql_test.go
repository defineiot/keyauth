package mysql_test

import (
	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/project"
	"github.com/defineiot/keyauth/dao/project/mysql"
	"github.com/defineiot/keyauth/internal/conf/mock"
)

func newTestStore() project.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	opt := &dao.Options{DB: db}
	store, err := mysql.NewProjectStore(opt)
	if err != nil {
		panic(err)
	}

	return store
}

type projectSuit struct {
	store project.Store
	p     *project.Project
	did   string
	uid   string
}

func (s *projectSuit) TearDown() {
	if s.p.ID != "" {
		s.store.DeleteProjectByID(s.p.ID)
	}

	s.store.Close()
}

func (s *projectSuit) SetUp() {
	s.did = "project-unit-test-domain"
	s.uid = "project-unit-test-user1-id"

	s.p = &project.Project{
		Name:        "unit-test-project1",
		Picture:     "unit-test-pic",
		Latitude:    "111",
		Longitude:   "222",
		Enabled:     true,
		Owner:       "unit-test-project-owner",
		Description: "",
		DomainID:    s.did,
	}

	s.store = newTestStore()

}
