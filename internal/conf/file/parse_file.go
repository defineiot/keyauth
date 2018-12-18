package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"

	"github.com/defineiot/keyauth/internal/conf"
)

var (
	config *conf.Config
	once   sync.Once
)

// NewFileConf from an file
func NewFileConf(filePath string) conf.Configer {
	return &fileConfig{filePath: filePath}
}

type fileConfig struct {
	filePath string
}

func (f *fileConfig) GetConf() (*conf.Config, error) {
	var err error

	once.Do(func() {
		err = initConfig(f.filePath)
	})

	if err != nil {
		return nil, err
	}

	if err = config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func initConfig(configpath string) error {
	configPath, err := filepath.Abs(configpath)
	if err != nil {
		return fmt.Errorf("get config file absolute path failed, %s", err.Error())
	}

	file, err := os.Open(configPath)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("open keyauth config file error, %s", err.Error())
	}

	fd, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("open keyauth config file error, %s", err.Error())
	}

	config = new(conf.Config)
	config.APP = new(conf.AppConf)
	config.Log = new(conf.LogConf)
	config.MySQL = new(conf.MySQLConf)
	config.Etcd = new(conf.ETCDConf)
	config.Token = new(conf.TokenConf)
	config.Admin = new(conf.AdminCount)
	config.Mail = new(conf.MailConf)
	config.SMS = new(conf.AliYunSMSConf)

	if err := toml.Unmarshal(fd, config); err != nil {
		return fmt.Errorf("load config file to json error, %s", err.Error())
	}

	return nil
}
