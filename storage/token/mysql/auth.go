package mysql

import "database/sql"

type manager struct {
	db *sql.DB
}
