package cmd

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"openauth/api/config"
	"openauth/api/config/env"
	"openauth/api/config/file"
	"openauth/api/http"
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

		// get db conn and logger
		db, err := conf.GetDBConn()
		if err != nil {
			return err
		}
		logger, err := conf.GetLogger()
		if err != nil {
			return err
		}

		// check the database is initial
		vers, desc, err := checkDBInit(db)
		if err != nil {
			return err
		}
		if vers < 1 {
			return fmt.Errorf("the database hasn't initialized")
		}
		logger.Debug(fmt.Sprintf("the database version: %d, desc: %s", vers, desc))

		// start service
		s, err := http.NewService(conf)
		if err != nil {
			return err
		}

		if err := s.Start(); err != nil {
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
		return nil, errors.New("unknown config type")
	}
	return conf, nil
}

// CheckDBInit use to check the mysql db is initial
func checkDBInit(db *sql.DB) (version int, desc string, err error) {
	// if table not exists, version: 0
	rows, err := db.Query("SHOW TABLES LIKE 'dbmanager'")
	if err != nil {
		return 0, "", fmt.Errorf("check table exists failed, %s", err)
	}
	tc := 0
	for rows.Next() {
		tc++
	}
	if tc == 0 {
		return 0, "", nil
	}

	// query the version
	rows, err = db.Query("SELECT version,description FROM dbmanager ORDER BY version DESC LIMIT 1")
	if err != nil {
		return 0, "", fmt.Errorf("check database version failed, %s", err)
	}
	count := 0
	for rows.Next() {
		if err := rows.Scan(&version, &desc); err != nil {
			return 0, "", fmt.Errorf("check database initial table failed, %s", err)
		}
		count++
	}

	return version, desc, nil
}

func init() {
	RootCmd.AddCommand(serviceCmd)
}
