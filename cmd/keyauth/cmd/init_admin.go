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
var initAdminCmd = &cobra.Command{
	Use:   "admin [init]",
	Short: "初始化服务的系统管理员账号",
	Long:  `初始化服务的系统管理员账号`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("缺少子命令")
		}

		switch args[0] {
		case "init":
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
		default:
			return errors.New("未知子命令")
		}

		return nil
	},
}

func init() {
	initAdminCmd.Flags().StringVarP(&adminAccount, "account", "u", "", "系统管理员账号名称")
	initAdminCmd.Flags().StringVarP(&adminPass, "password", "p", "", "系统管理员秘密")
	initAdminCmd.Flags().StringVarP(&confType, "config-type", "t", "file", "服务配置类型 [file/env/etcd]")
	initAdminCmd.Flags().StringVarP(&confFile, "config-file", "f", "cmd/etc/keyauth.conf", "如果服务采用配置文件配置时, 配置文件的具体路径")
	initAdminCmd.Flags().StringVarP(&confEtcd, "config-etcd", "e", "127.0.0.1:2379", "如果服务采用Etcd来存储配置时, Etcd的服务地址")
	RootCmd.AddCommand(initAdminCmd)
}
