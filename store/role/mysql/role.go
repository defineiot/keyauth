package mysql

import (
	"database/sql"
)

// RoleManager is use mongodb as storage
type RoleManager struct {
	DB *sql.DB
}
