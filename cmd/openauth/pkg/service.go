package pkg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"openauth/api/config"
	"openauth/api/config/env"
	"openauth/api/config/file"
)

// startCmd represents the start command
var serviceCmd = &cobra.Command{
	Use:   "service [start/stop/reload/restart]",
	Short: "management openauth api",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return errors.New("[start/stop/reload/restart] are required")
		}

		conf, err := checkConfType(confType)
		if err != nil {
			return err
		}

		return nil
	},
}

func checkConfType(configType string) (conf *config.Config, err error) {

	switch configType {
	case "file":
		envconf := file.NewFileConf(confFile)
		conf, err = envconf.GetConf()
		if err != nil {
			return nil, fmt.Errorf("new file conf error, %s", err.Error())
		}
	case "env":
		fileconf := env.NewConfigManager()
		conf, err = fileconf.GetConf()
		if err != nil {
			return nil, fmt.Errorf("new env conf error, %s", err.Error())
		}
	case "etcd":
		return nil, errors.New("not implemented")
	default:
		return nil, errors.New("unknow config type")
	}

	return
}

// CheckDBInit use to check the mysql db is initial
func checkDBInit(db *sql.DB) (version int, desc string, err error) {
	rows, errSE := db.Query("SELECT version,description FROM dbmanager")
	if err != nil {
		err = fmt.Errorf("check database initial failed, %s", errSE.Error())
		return
	}

	count := 0
	for rows.Next() {

		rows.Columns()
		if errSC := rows.Scan(&version, &desc); err != nil {
			err = fmt.Errorf("check database initial table failed, %s", errSC.Error())
			return
		}
		count++

	}

	if count == 0 {
		err = errors.New("database initial failed, there is no dbversion in dbmanager")
		return
	}

	return
}

func init() {

	RootCmd.AddCommand(serviceCmd)

}
