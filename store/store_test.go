package store_test

import (
	"github.com/defineiot/keyauth/internal/conf/mock"
	"github.com/defineiot/keyauth/store"
)

func newTestStore() *store.Store {
	conf := mock.NewConfig()
	mystore, err := store.NewStore(conf)
	if err != nil {
		panic(err)
	}

	return mystore
}

type storeSuit struct {
	store *store.Store
}

func (s *storeSuit) TearDown() {

}

func (s *storeSuit) SetUp() {

	s.store = newTestStore()

}
