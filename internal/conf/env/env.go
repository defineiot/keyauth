package env

import (
	"os"
	"strconv"
	"strings"
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

func initConfig() error {
	config = new(conf.Config)
	config.APP = new(conf.AppConf)
	config.Log = new(conf.LogConf)
	config.MySQL = new(conf.MySQLConf)
	config.Token = new(conf.TokenConf)
	config.Etcd = new(conf.ETCDConf)

	config.APP.Host = os.Getenv("APP_HOST")
	config.APP.Port = os.Getenv("APP_PORT")
	config.APP.Key = os.Getenv("APP_KEY")
	config.APP.Name = os.Getenv("APP_NAME")

	config.Log.FilePath = os.Getenv("LOG_FILE_PATH")
	config.Log.Level = os.Getenv("LOG_LEVEL")

	config.MySQL.Host = os.Getenv("MYSQL_HOST")
	config.MySQL.Port = os.Getenv("MYSQL_PORT")
	config.MySQL.DB = os.Getenv("MYSQL_DB")
	config.MySQL.User = os.Getenv("MYSQL_USER")
	config.MySQL.Pass = os.Getenv("MYSQL_PASS")
	idle, _ := strconv.Atoi(os.Getenv("MYSQL_MAX_IDEL_CONN"))
	open, _ := strconv.Atoi(os.Getenv("MYSQL_MAX_OPEN_CONN"))
	mttl, _ := strconv.Atoi(os.Getenv("MYSQL_MAX_LIFE_TIME"))
	config.MySQL.MaxIdleConn = idle
	config.MySQL.MaxOpenConn = open
	config.MySQL.MaxLifeTime = mttl

	expireIN, _ := strconv.ParseInt(os.Getenv("TOKEN_EXPIRES_IN"), 10, 64)
	config.Token.ExpiresIn = expireIN
	config.Token.Type = os.Getenv("TOKEN_TYPE")

	if os.Getenv("ETCD_REGIST_FEATURES_ENABLE") != "" {
		config.Etcd.EnableRegisteFeatures = true
	}
	if os.Getenv("ETCD_REGIST_INSTANCE_ENABLE") != "" {
		config.Etcd.EnableRegisteInstance = true
	}
	ittl, _ := strconv.Atoi(os.Getenv("INSTANCE_TTL"))
	config.Etcd.InstanceTTL = ittl
	config.Etcd.InstanceRegistryPrefix = os.Getenv("REGIST_INSTANCE_PREFIX")
	config.Etcd.Endpoints = strings.Split(os.Getenv("ETCD_ENDPOINTS"), ",")
	config.Etcd.UserName = os.Getenv("ETCD_USERNAME")
	config.Etcd.Password = os.Getenv("ETCD_PASSWORD")

	if err := config.Validate(); err != nil {
		return err
	}

	return nil
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
