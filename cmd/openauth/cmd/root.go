package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"openauth/api/version"
)

var (
	// database management option
	sqlFile string

	// openauth config option
	confType string
	confFile string
	confEtcd string

	vers bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "openauth",
	Short: "openauth is an multi tenant user management system based on oauth2.0",
	Long:  `openauth ...`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(version.FullVersion())
			return nil
		}

		return errors.New("no flags find")
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&confType, "config-type", "t", "file", "the openauth service config type [file/env/etcd]")
	RootCmd.PersistentFlags().StringVarP(&confFile, "config-file", "f", "conf/openauth.conf", "the openauth service config from file")
	RootCmd.PersistentFlags().StringVarP(&confEtcd, "config-etcd", "e", "127.0.0.1:2379", "the openauth service config from etcd")
	RootCmd.PersistentFlags().BoolVarP(&vers, "version", "v", false, "the openauth service version")
}
