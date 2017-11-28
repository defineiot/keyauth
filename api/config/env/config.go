package env

import (
	"openauth/api/config"
	"os"
	"sync"
)

var (
	conf *config.Config
	once sync.Once
)

// NewConfigManager use to get env config
func NewConfigManager() config.Configer {
	return &envConfig{}
}

type envConfig struct {
}

func (e *envConfig) GetConf() (*config.Config, error) {
	var err error

	once.Do(func() {
		err = initConfig()
	})

	if err != nil {
		return nil, err
	}

	return conf, nil
}

func initConfig() error {
	conf = &config.Config{}

	if err := getAPP(conf); err != nil {
		return err
	}
	if err := getMySQL(conf); err != nil {
		return err
	}

	if err := conf.Validate(); err != nil {
		return err
	}

	return nil
}

func getAPP(conf *config.Config) error {
	conf.APP.Host = os.Getenv("OA_APP_HOST")
	conf.APP.Port = os.Getenv("OA_APP_PORT")

	return nil
}

func getMySQL(conf *config.Config) error {
	conf.MySQL.Host = os.Getenv("OA_MYSQL_HOST")
	conf.MySQL.Port = os.Getenv("OA_MYSQL_PORT")
	conf.MySQL.User = os.Getenv("OA_MYSQL_USER")
	conf.MySQL.Pass = os.Getenv("OA_MYSQL_PASS")
	conf.MySQL.DB = os.Getenv("OA_MYSQL_DB")

	return nil
}
