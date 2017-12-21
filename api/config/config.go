package config

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"openauth/api/logger"
	"openauth/api/logger/logrus"
)

var (
	db         *sql.DB
	oalogger   logger.OpenAuthLogger
	dbOnce     sync.Once
	loggerOnce sync.Once
)

// Configer use to get conf
type Configer interface {
	GetConf() (*Config, error)
}

// Config is service conf
type Config struct {
	APP   *AppConf   `toml:"app"`
	MySQL *MySQLConf `toml:"mysql"`
	Log   *LogConf   `toml:"log"`
	Token *Token     `toml:"token"`
}

// AppConf service config
type AppConf struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
	Key  string `toml:"key"`
	Name string `toml:"name"`
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

// LogConf log config
type LogConf struct {
	Level    string `toml:"level"`
	FilePath string `toml:"path"`
}

// Token is oauth token
type Token struct {
	Type      string `toml:"type"`
	ExpiresIn int64  `toml:"expires_in"`
}

// Validate use to check the service config
func (c *Config) Validate() error {
	if err := c.validateAPP(); err != nil {
		return err
	}
	if err := c.validateMySQL(); err != nil {
		return err
	}
	if err := c.validateLog(); err != nil {
		return err
	}
	if err := c.validateToken(); err != nil {
		return err
	}

	return nil
}

func (c *Config) validateAPP() error {
	if c.APP == nil {
		c.APP = &AppConf{}
	}

	if c.APP.Host == "" {
		c.APP.Host = "0.0.0.0"
	}
	if c.APP.Port == "" {
		c.APP.Port = "8080"
	}
	if c.APP.Name == "" {
		c.APP.Name = "openauth"
	}

	if c.APP.Key == "" {
		return errors.New("app key isn't config")
	}

	return nil
}

func (c *Config) validateMySQL() error {
	if c.MySQL == nil {
		c.MySQL = &MySQLConf{}
	}

	if c.MySQL.Host == "" {
		c.MySQL.Host = "127.0.0.1"
	}
	if c.MySQL.Port == "" {
		c.MySQL.Port = "3306"
	}

	if c.MySQL.User == "" || c.MySQL.Pass == "" || c.MySQL.DB == "" {
		return errors.New("mysql user or pass or db isn't config")
	}

	return nil
}

func (c *Config) validateLog() error {
	if c.Log == nil {
		c.Log = new(LogConf)
	}

	if c.Log.Level == "" {
		c.Log.Level = "info"
	}

	if c.Log.FilePath == "" {
		return errors.New("log path not config")
	}

	return nil
}

func (c *Config) validateToken() error {
	if c.Token == nil {
		c.Token = new(Token)
	}

	if c.Token.Type == "" {
		c.Token.Type = "bearer"
	}
	if c.Token.ExpiresIn == 0 {
		c.Token.ExpiresIn = 3600
	}

	return nil
}

// GetDBConn use to get mysql database connection
func (c *Config) GetDBConn() (*sql.DB, error) {
	var err error

	dbOnce.Do(func() {
		err = c.initDBConn()
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetLogger use to get logger instance
func (c *Config) GetLogger() (logger.OpenAuthLogger, error) {
	var err error

	opts := logger.Opts{Name: c.APP.Name, Level: c.Log.Level, FilePath: c.Log.FilePath}
	loggerOnce.Do(func() {
		oalogger, err = logrus.NewLogrusLogger(&opts)
	})

	if err != nil {
		return nil, err
	}

	return oalogger, nil
}

func (c *Config) initDBConn() error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", c.MySQL.User, c.MySQL.Pass, c.MySQL.Host, c.MySQL.Port, c.MySQL.DB)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	db.SetMaxOpenConns(c.MySQL.MaxOpenConn)
	db.SetMaxIdleConns(c.MySQL.MaxIdleConn)
	db.SetConnMaxLifetime(time.Minute * time.Duration(c.MySQL.MaxLifeTime))
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}

	return nil
}
