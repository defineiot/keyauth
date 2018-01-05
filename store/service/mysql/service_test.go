package mysql_test

import (
	"testing"
)

func TestSaveService(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	if _, err := s.SaveService("test_service01", "just for unit test"); err != nil {
		t.Fatal(err)
	}
}
