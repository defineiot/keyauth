package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplicationSuit(t *testing.T) {
	suit := new(applicationSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("CreateApplicationOK", testCreateApplicationOK(suit))
	t.Run("ListUserApplicationsOK", testListUserApplicationsOK(suit))
	t.Run("GetApplicationOK", testGetUserApplicationOK(suit))
	t.Run("DeleteApplicationOK", testDeleteApplicationOK(suit))
}

func testCreateApplicationOK(s *applicationSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.CreateApplication(s.app)
		should.NoError(err)

		t.Logf("create application(%s) success: %s", s.app.Name, s.app)
	}
}

func testListUserApplicationsOK(s *applicationSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		apps, err := s.store.ListUserApplications(s.userID)
		should.NoError(err)

		t.Logf("list user applications(%s) success: %s", s.userID, apps)
	}
}

func testGetUserApplicationOK(s *applicationSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		app, err := s.store.GetApplication(s.app.ID)
		should.NoError(err)

		t.Logf("get application(%s) success: %s", s.app.ID, app)
	}
}

func testDeleteApplicationOK(s *applicationSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		err := s.store.DeleteApplication(s.app.ID)
		should.NoError(err)

		t.Logf("delete application(%s) success: %s", s.app.ID, s.app)
	}
}
