package conf

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/defineiot/keyauth/internal/logger"
	"github.com/defineiot/keyauth/internal/logger/logrus"
)

// Configer use to get conf
type Configer interface {
	GetConf() (*Config, error)
}

// Config is service conf
type Config struct {
	APP        *AppConf   `toml:"app"`
	Log        *LogConf   `toml:"log"`
	MySQL      *MySQLConf `toml:"mysql"`
	Etcd       *ETCDConf  `toml:"etcd"`
	Token      *TokenConf `toml:"token"`
	logger     logger.Logger
	loggerOnce sync.Once

	mailOnce sync.Once
}

// AppConf service config
type AppConf struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
	Key  string `toml:"key"`
	Name string `toml:"name"`
}

// LogConf log config
type LogConf struct {
	Level    string `toml:"level"`
	FilePath string `toml:"path"`
}

// MySQLConf mysql config
type MySQLConf struct {
	Host        string `toml:"host"`
	Port        string `toml:"port"`
	User        string `toml:"user"`
	Pass        string `toml:"pass"`
	DB          string `toml:"db"`
	MaxOpenConn int    `toml:"max_open_conn"`
	MaxIdleConn int    `toml:"max_idle_conn"`
	MaxLifeTime int    `toml:"max_life_time"`
}

// ETCDConf todo
type ETCDConf struct {
	EnableRegisteInstance  bool     `toml:"enable_registe_instance"`
	InstanceTTL            int      `toml:"instance_ttl"`
	EnableRegisteFeatures  bool     `toml:"enable_registe_features"`
	InstanceRegistryPrefix string   `toml:"instance_registry_prefix"`
	Endpoints              []string `toml:"endpoints"`
	UserName               string   `toml:"username"`
	Password               string   `toml:"password"`
}

// TokenConf token config
type TokenConf struct {
	Type      string `toml:"type"`
	ExpiresIn int64  `toml:"expires_in"`
}

// Validate use to check the service config
func (c *Config) Validate() error {
	if c.APP.Name == "" {
		c.APP.Name = "github.com/defineiot/keyauth"
	}
	if c.APP.Host == "" {
		c.APP.Host = "127.0.0.1"
	}
	if c.APP.Port == "" {
		c.APP.Port = "8080"
	}
	if c.APP.Key == "" {
		c.APP.Key = "qjPmB1agz5Au1So5TVIpPiLyEYN3SKOK"
	}

	if c.Log.Level == "" {
		c.Log.Level = "debug"
	}
	if c.Log.FilePath == "" {
		c.Log.FilePath = "./log/data-gateway.log"
	}

	if c.MySQL.MaxIdleConn == 0 {
		c.MySQL.MaxIdleConn = 80
	}

	if c.MySQL.MaxLifeTime == 0 {
		c.MySQL.MaxLifeTime = 60
	}

	if c.MySQL.MaxOpenConn == 0 {
		c.MySQL.MaxOpenConn = 200
	}

	if c.MySQL.Host == "" || c.MySQL.Port == "" || c.MySQL.User == "" || c.MySQL.Pass == "" {
		return errors.New("mysql host,port,user,pass required")
	}

	if len(c.Etcd.Endpoints) == 0 {
		return errors.New("etcd endpoints is required")
	}
	if c.Etcd.InstanceRegistryPrefix == "" {
		c.Etcd.InstanceRegistryPrefix = "/registry/instances"
	}

	if c.Token.Type == "" {
		c.Token.Type = "bearer"
	}

	if c.Token.ExpiresIn == 0 {
		c.Token.ExpiresIn = 86400
	}

	return nil
}

// GetLogger use to get logger instanc
func (c *Config) GetLogger() (logger.Logger, error) {
	var err error

	opts := logger.Opts{Name: c.APP.Name, Level: c.Log.Level, FilePath: c.Log.FilePath}
	c.loggerOnce.Do(func() {
		c.logger, err = logrus.NewLogrusLogger(&opts)
	})

	if err != nil {
		return nil, err
	}

	return c.logger, nil
}

// GetDBConn use to get db connection pool
func (c *Config) GetDBConn() (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", c.MySQL.User, c.MySQL.Pass, c.MySQL.Host, c.MySQL.Port, c.MySQL.DB)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	db.SetMaxOpenConns(c.MySQL.MaxOpenConn)
	db.SetMaxIdleConns(c.MySQL.MaxIdleConn)
	db.SetConnMaxLifetime(time.Minute * time.Duration(c.MySQL.MaxLifeTime))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}

	return db, nil
}
