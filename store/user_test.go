package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectSuit(t *testing.T) {
	suit := new(storeSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("InitSystemOK", testInitSystemOK(suit))
}

func testInitSystemOK(s *storeSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.InitAdmin("admin", "password")
		should.NoError(err)
	}
}
