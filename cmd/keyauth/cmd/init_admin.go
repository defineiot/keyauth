package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/defineiot/keyauth/store"
)

var (
	adminAccount string
	adminPass    string
)

// startCmd represents the start command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化服务的系统管理员账号",
	Long:  `初始化服务的系统管理员账号`,
	RunE: func(cmd *cobra.Command, args []string) error {

		conf, err := checkConfType(confType)
		if err != nil {
			return err
		}

		if adminAccount == "" || adminPass == "" {
			return errors.New("系统管理员的账号名称或者账号秘密未设置")
		}

		s, err := store.NewStore(conf)
		if err != nil {
			return err
		}

		if err := s.InitAdmin(adminAccount, adminPass); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&adminAccount, "account", "u", "", "系统管理员账号名称")
	initCmd.Flags().StringVarP(&adminPass, "password", "p", "", "系统管理员秘密")
	initCmd.Flags().StringVarP(&confType, "config-type", "t", "file", "服务配置类型 [file/env/etcd]")
	initCmd.Flags().StringVarP(&confFile, "config-file", "f", "cmd/etc/keyauth.conf", "如果服务采用配置文件配置时, 配置文件的具体路径")
	initCmd.Flags().StringVarP(&confEtcd, "config-etcd", "e", "127.0.0.1:2379", "如果服务采用Etcd来存储配置时, Etcd的服务地址")
	RootCmd.AddCommand(initCmd)
}
