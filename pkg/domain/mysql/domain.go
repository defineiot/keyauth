package mongodb

import (
	"database/sql"
)

// DomainManager is use mongodb as storage
type DomainManager struct {
	DB *sql.DB
}
