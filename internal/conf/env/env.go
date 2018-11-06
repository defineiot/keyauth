package env

import (
	"os"
	"sync"

	"github.com/defineiot/keyauth/internal/conf"
)

var (
	config *conf.Config
	once   sync.Once
)

// NewEnvConf use to get env config
func NewEnvConf() conf.Configer {
	return &envConfig{}
}

type envConfig struct {
}

func (e *envConfig) GetConf() (*conf.Config, error) {
	var err error

	once.Do(func() {
		err = initConfig()
	})

	if err != nil {
		return nil, err
	}

	if err = config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func initConfig() error {
	config = new(conf.Config)
	config.APP = new(conf.AppConf)
	config.Log = new(conf.LogConf)
	config.MySQL = new(conf.MySQLConf)
	config.Token = new(conf.TokenConf)

	if err := getAPP(config); err != nil {
		return err
	}
	if err := getLog(config); err != nil {
		return err
	}

	if err := config.Validate(); err != nil {
		return err
	}

	return nil
}

func getAPP(conf *conf.Config) error {

	conf.APP.Host = os.Getenv("APP_HOST")
	conf.APP.Port = os.Getenv("APP_PORT")
	conf.APP.Key = os.Getenv("APP_KEY")
	conf.APP.Name = os.Getenv("APP_NAME")

	return nil
}

func getLog(conf *conf.Config) error {

	conf.Log.FilePath = os.Getenv("LOG_FILE_PATH")
	conf.Log.Level = os.Getenv("LOG_LEVEL")

	return nil
}
