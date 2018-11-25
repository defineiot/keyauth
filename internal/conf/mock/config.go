package mock

import (
	"github.com/defineiot/keyauth/internal/conf"
)

// NewConfig mock an conf for test
func NewConfig() *conf.Config {
	app := new(conf.AppConf)
	app.Host = "0.0.0.0"
	app.Key = "default"
	app.Name = "keyauth"
	app.Port = "8080"
	mysql := new(conf.MySQLConf)
	mysql.Host = "127.0.0.1"
	mysql.Port = "3306"
	mysql.DB = "keyauth"
	mysql.User = "root"
	mysql.Pass = "passwd"
	log := new(conf.LogConf)
	log.FilePath = "/tmp/keyauth.log"
	log.Level = "debug"
	token := new(conf.TokenConf)
	token.Type = "bearer"
	token.ExpiresIn = 86400

	conf := new(conf.Config)
	conf.APP = app
	conf.MySQL = mysql
	conf.Token = token
	conf.Log = log

	return conf
}
