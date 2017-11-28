package mysql

import (
	"database/sql"
)

// ProjectManager is use mongodb as storage
type ProjectManager struct {
	DB *sql.DB
}
