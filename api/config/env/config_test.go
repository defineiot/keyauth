package env_test

import (
	"os"
	"testing"

	"openauth/api/config/env"
)

func TestNewEnvConf(t *testing.T) {
	fakeEnv()
	env := env.NewEnvConf()
	_, err := env.GetConf()
	if err != nil {
		t.Fatal(err)
	}

}

func TestEnvStruct(t *testing.T) {

}

func fakeEnv() {
	os.Setenv("OA_APP_KEY", "default")
	os.Setenv("OA_MYSQL_USER", "openauth")
	os.Setenv("OA_MYSQL_PASS", "openauth")
	os.Setenv("OA_MYSQL_DB", "openauth")
	os.Setenv("OA_LOG_FILE_PATH", "log/debug.log")
}
