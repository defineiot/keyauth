package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"

	"openauth/api/config"
)

var (
	conf *config.Config
	once sync.Once
)

// NewFileConf from an file
func NewFileConf(filePath string) config.Configer {
	return &fileConfig{filePath: filePath}
}

type fileConfig struct {
	filePath string
}

func (f *fileConfig) GetConf() (*config.Config, error) {
	var err error

	once.Do(func() {
		err = initConfig(f.filePath)
	})

	if err != nil {
		return nil, err
	}

	return conf, nil
}

func initConfig(configpath string) error {
	configPath, err := filepath.Abs(configpath)
	if err != nil {
		return fmt.Errorf("get config file absolute path failed, %s", err.Error())
	}

	file, err := os.Open(configPath)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("open openauth config file error, %s", err.Error())
	}

	fd, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("open openauth config file error, %s", err.Error())
	}

	conf := &config.Config{}
	if err := toml.Unmarshal(fd, conf); err != nil {
		return fmt.Errorf("load config file to json error, %s", err.Error())
	}

	return nil
}
