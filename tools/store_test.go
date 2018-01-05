package tools_test

import (
	"testing"

	"openauth/api/config/mock"
	"openauth/tools"
)

func TestPrepareStmts(t *testing.T) {
	t.Run("OK", testPrepareOK)
	t.Run("Failed", testPrepareFailed)
}

func testPrepareOK(t *testing.T) {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	unprepared := map[string]string{
		"Test": `select * from dbmanager`,
	}

	if _, err := tools.PrepareStmts(db, unprepared); err != nil {
		t.Fatal(err)
	}
}

func testPrepareFailed(t *testing.T) {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	unprepared := map[string]string{
		"Test": `select unknow from dbmanager`,
	}

	if _, err := tools.PrepareStmts(db, unprepared); err == nil {
		t.Fatal("need an prepare failed error")
	}
}
