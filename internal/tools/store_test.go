package tools_test

import (
	"testing"

	"github.com/defineiot/keyauth/internal/conf"
	"github.com/defineiot/keyauth/internal/tools"
)

func newConfig() *conf.Config {
	app := new(conf.AppConf)
	app.Host = "0.0.0.0"
	app.Key = "default"
	app.Name = "keyauth"
	app.Port = "8080"
	mysql := new(conf.MySQLConf)
	mysql.Host = "192.168.0.203"
	mysql.Port = "3306"
	mysql.DB = "keyauth"
	mysql.User = "keyauth"
	mysql.Pass = "keyauth"
	log := new(conf.LogConf)
	log.FilePath = "/tmp/keyauth.log"
	log.Level = "debug"

	conf := new(conf.Config)
	conf.APP = app
	conf.MySQL = mysql
	conf.Log = log

	return conf
}

func TestPrepareStmts(t *testing.T) {
	t.Run("OK", testPrepareOK)
	t.Run("Failed", testPrepareFailed)
}

func testPrepareOK(t *testing.T) {
	conf := newConfig()
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
	conf := newConfig()
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
