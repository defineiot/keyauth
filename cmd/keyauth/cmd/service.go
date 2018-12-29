package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/defineiot/keyauth/api/http"
	"github.com/defineiot/keyauth/internal/conf"
	"github.com/defineiot/keyauth/internal/conf/env"
	"github.com/defineiot/keyauth/internal/conf/file"
)

var (
	// pusher service config option
	confType string
	confFile string
	confEtcd string
)

// startCmd represents the start command
var serviceCmd = &cobra.Command{
	Use:   "service [start/stop/reload/restart]",
	Short: "management keyauth service",
	Long:  `use to start keyauth service`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return errors.New("[start/bootstrap] are required")
		}

		conf, err := checkConfType(confType)
		if err != nil {
			return err
		}

		switch args[0] {
		case "start":
			// start service
			s, err := http.NewService(conf)
			if err != nil {
				return err
			}

			if err := s.Start(); err != nil {
				return err
			}

		case "bootstrap":
			// start service
			s, err := http.NewService(conf)
			if err != nil {
				return err
			}

			if err := s.BootStrap(); err != nil {
				return err
			}
		default:
			return errors.New("not support argument, support [start/bootstrap]")
		}

		return nil
	},
}

func checkConfType(configType string) (conf *conf.Config, err error) {
	switch configType {
	case "file":
		fileconf := file.NewFileConf(confFile)
		conf, err = fileconf.GetConf()
		if err != nil {
			return nil, fmt.Errorf("new file conf error, %s", err.Error())
		}
	case "env":
		envconf := env.NewEnvConf()
		conf, err = envconf.GetConf()
		if err != nil {
			return nil, fmt.Errorf("new env conf error, %s", err.Error())
		}
	case "etcd":
		return nil, errors.New("not implemented")
	default:
		return nil, errors.New("unknown config type")
	}
	return conf, nil
}

func init() {
	serviceCmd.Flags().StringVarP(&confType, "config-type", "t", "file", "the keyauth service config type [file/env/etcd]")
	serviceCmd.Flags().StringVarP(&confFile, "config-file", "f", "cmd/etc/keyauth.conf", "the keyauth service config from file")
	serviceCmd.Flags().StringVarP(&confEtcd, "config-etcd", "e", "127.0.0.1:2379", "the keyauth service config from etcd")
	RootCmd.AddCommand(serviceCmd)
}
