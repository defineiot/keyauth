package tools

import (
	"database/sql"
	"fmt"
)

// PrepareStmts will attempt to prepare each unprepared
// query on the database. If one fails, the function returns
// with an error.
func PrepareStmts(db *sql.DB, unprepared map[string]string) (map[string]*sql.Stmt, error) {
	prepared := map[string]*sql.Stmt{}
	for k, v := range unprepared {
		stmt, err := db.Prepare(v)
		if err != nil {
			return nil, fmt.Errorf("prepare statment: %s, %s", k, err)
		}
		prepared[k] = stmt
	}

	return prepared, nil
}
