package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:   "database [init]",
	Short: "management openauth database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("[init] are required")
		}

		switch args[0] {
		case "init":
			if err := initDatabase(); err != nil {
				return err
			}
			fmt.Println("initial database successful")
		default:
			return errors.New("unknow argument, see usage")
		}

		return nil
	},
}

func initDatabase() error {
	f, err := os.Open(sqlFile)
	if err != nil {
		return err
	}

	sql, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	conf, err := checkConfType(confType)
	if err != nil {
		return err
	}

	db, err := conf.GetDBConn()
	if err != nil {
		return fmt.Errorf("get mysql connection error, %s", err.Error())
	}

	if _, err := db.Exec(string(sql)); err != nil {
		return fmt.Errorf("[ERROR] - [MySQL Create Table] - create table error, %s", err.Error())
	}

	vers, desc, err := checkDBInit(db)
	if err != nil {
		return err
	}
	if vers == 1 {
		return fmt.Errorf("the database has been initialized, sql version: %d, description: %s", vers, desc)
	}

	if err := initialRecord(db); err != nil {
		return fmt.Errorf("initial system info to database error, %s", err.Error())
	}

	return nil
}

// 初始化系统数据:
// 1. role:
//     - root : 超级管理员
//     - admin: 域管理员
// 2. domain:
//     - root : 超级管理员的domain
// 3. project:
//     - default: 超级管理员的默认项目
// 4. user:
//     - root : root/123456

// InitialDB use to initial default information.
func initialRecord(db *sql.DB) error {

	// var userID string

	// domainID, err := manager.InsertDomain("SystemAdmin", "openauth system administrator")
	// if err != nil {
	// 	return err
	// }

	// if _, err = manager.InsertProject(domainID, "default", "openauth system administrator's default project"); err != nil {
	// 	return err
	// }

	// if _, err = manager.InsertRole("SystemAdmin", "openauth system administrator's role"); err != nil {
	// 	return err
	// }

	// if _, err = manager.InsertRole("DomainAdmin", "openauth domain administrator's role"); err != nil {
	// 	return err
	// }

	// if userID, err = manager.InsertUser(domainID, "sysadmin", "123456"); err != nil {
	// 	return err
	// }

	// ok, err := manager.CheckUserExists(userID)
	// if err != nil {
	// 	return err
	// }

	// if !ok {
	// 	return errors.New("user not exist")
	// }

	// if _, err = manager.InsertClient(userID, "app", "appclientsecret", "admin client", "", 0); err != nil {
	// 	return err
	// }

	return nil
}

func init() {
	RootCmd.AddCommand(databaseCmd)

	databaseCmd.Flags().StringVarP(&sqlFile, "sqlfile", "s", "conf/dbsql/openauthv01.sql", "the initial database sql file")
}
