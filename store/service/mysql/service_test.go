package mysql_test

import (
	"testing"
)

func TestSaveService(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	svr, err := s.SaveService("test_service01", "just for unit test")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := s.FindAllService(); err != nil {
		t.Fatal(err)
	}

	if _, err := s.FindServiceByID(svr.ID); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteService(svr.ID); err != nil {
		t.Fatal(err)
	}

}
