package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/defineiot/keyauth/store"
)

var (
	sqlFilePath string
)

// startCmd represents the start command
var initDBCmd = &cobra.Command{
	Use:   "db [init]",
	Short: "初始化服务数据库",
	Long:  `初始化服务数据库`,
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

			if sqlFilePath == "" {
				return errors.New("请指定SQL文件的路径")
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
	initDBCmd.Flags().StringVarP(&sqlFilePath, "path", "p", "cmd/ddl/schema_v1.sql", "sql文件的路径")
	initDBCmd.Flags().StringVarP(&confType, "config-type", "t", "file", "服务配置类型 [file/env/etcd]")
	initDBCmd.Flags().StringVarP(&confFile, "config-file", "f", "cmd/etc/keyauth.conf", "如果服务采用配置文件配置时, 配置文件的具体路径")
	initDBCmd.Flags().StringVarP(&confEtcd, "config-etcd", "e", "127.0.0.1:2379", "如果服务采用Etcd来存储配置时, Etcd的服务地址")
	RootCmd.AddCommand(initCmd)
}
