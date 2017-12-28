package mock

import (
	"openauth/api/config"
)

// NewConfig mock an conf for test
func NewConfig() *config.Config {
	app := new(config.AppConf)
	app.Host = "0.0.0.0"
	app.Key = "default"
	app.Name = "openauth"
	app.Port = "8080"
	mysql := new(config.MySQLConf)
	mysql.Host = "127.0.0.1"
	mysql.Port = "3306"
	mysql.DB = "openauth"
	mysql.User = "openauth"
	mysql.Pass = "openauth"
	log := new(config.LogConf)
	log.FilePath = "/tmp/debug.log"
	log.Level = "debug"
	token := new(config.Token)
	token.ExpiresIn = 3306
	token.Type = "bearer"

	conf := new(config.Config)
	conf.APP = app
	conf.MySQL = mysql
	conf.Log = log
	conf.Token = token

	return conf
}
