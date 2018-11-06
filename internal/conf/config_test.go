package conf_test

import (
	"testing"

	"github.com/defineiot/keyauth/internal/conf/mock"
)

func TestConfigValidate(t *testing.T) {
	conf := mock.NewConfig()
	if err := conf.Validate(); err != nil {
		t.Fatal(err)
	}

}

func TestConfigGetDB(t *testing.T) {
	conf := mock.NewConfig()
	_, err := conf.GetDBConn()
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigGetLogger(t *testing.T) {
	conf := mock.NewConfig()
	_, err := conf.GetLogger()
	if err != nil {
		t.Fatal(err)
	}
}
