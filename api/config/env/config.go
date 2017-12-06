package env

import (
	"openauth/api/config"
	"os"
	"strconv"
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
	if err := getLog(conf); err != nil {
		return err
	}

	if err := conf.Validate(); err != nil {
		return err
	}

	return nil
}

func getAPP(conf *config.Config) error {
	if conf.APP == nil {
		conf.APP = &config.AppConf{}
	}

	conf.APP.Host = os.Getenv("OA_APP_HOST")
	conf.APP.Port = os.Getenv("OA_APP_PORT")
	conf.APP.Key = os.Getenv("OA_APP_KEY")
	conf.APP.Name = os.Getenv("OA_APP_NAME")

	return nil
}

func getMySQL(conf *config.Config) error {
	if conf.MySQL == nil {
		conf.MySQL = &config.MySQLConf{}
	}

	conf.MySQL.Host = os.Getenv("OA_MYSQL_HOST")
	conf.MySQL.Port = os.Getenv("OA_MYSQL_PORT")
	conf.MySQL.User = os.Getenv("OA_MYSQL_USER")
	conf.MySQL.Pass = os.Getenv("OA_MYSQL_PASS")
	conf.MySQL.DB = os.Getenv("OA_MYSQL_DB")

	var err error
	if openconn := os.Getenv("OA_MYSQL_MAX_OPEN_CONN"); openconn != "" {
		if conf.MySQL.MaxOpenConn, err = strconv.Atoi(openconn); err != nil {
			return err
		}
	} else {
		conf.MySQL.MaxOpenConn = 1000
	}

	if idelconn := os.Getenv("OA_MYSQL_MAX_IDEL_CONN"); idelconn != "" {
		if conf.MySQL.MaxIdleConn, err = strconv.Atoi(idelconn); err != nil {
			return err
		}
	} else {
		conf.MySQL.MaxIdleConn = 200
	}

	if maxlife := os.Getenv("OA_MYSQL_MAX_LIFE_TIME"); maxlife != "" {
		if conf.MySQL.MaxLifeTime, err = strconv.Atoi(maxlife); err != nil {
			return err
		}
	} else {
		conf.MySQL.MaxLifeTime = 60
	}

	return nil
}

func getLog(conf *config.Config) error {
	if conf.Log == nil {
		conf.Log = &config.LogConf{}
	}

	conf.Log.FilePath = os.Getenv("OA_LOG_FILE_PATH")
	conf.Log.Level = os.Getenv("OA_LOG_LEVEL")

	return nil
}
