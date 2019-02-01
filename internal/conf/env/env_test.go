package env_test

import (
	"os"
	"testing"

	"github.com/defineiot/keyauth/internal/conf/env"
)

func TestNewEnvConf(t *testing.T) {
	fakeEnv()
	envconf := env.NewEnvConf()
	conf, err := envconf.GetConf()
	if err != nil {
		t.Fatal(err)
	}

	if conf.APP.Key != "default" {
		t.Fatal("the key not default")
	}
}

func fakeEnv() {
	os.Setenv("APP_KEY", "default")
	os.Setenv("MYSQL_HOST", "127.0.0.1")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DB", "keyauth")
	os.Setenv("MYSQL_USER", "root")
	os.Setenv("MYSQL_PASS", "passwd")
}
