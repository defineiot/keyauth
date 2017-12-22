package mysql_test

import (
	"testing"
)

func TestNewDomainStore(t *testing.T) {

	// _, err = mysql.NewDomainStore(db)
	// if err != nil {
	// 	t.Error(err)
	// }
}

func TestPrepareStatements(t *testing.T) {
	// Connect to the database
	// dsn := fmt.Sprintf(
	// 	"%s:%s@tcp(%s:%s)/?parseTime=true",
	// 	os.Getenv("DB_USERNAME"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_PORT"),
	// )

	// db, err := sql.Open("mysql", dsn)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// unprepared := map[string]string{
	// 	"test": "SELECT 1",
	// }

	// stmts, err := prepareStmts(db, unprepared)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if len(stmts) != 1 {
	// 	t.Fatalf("incorrect number of statements prepared; got %v want %v\n", len(stmts), 1)
	// }
}
