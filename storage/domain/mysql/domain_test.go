package mysql_test

import (
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func Test_CreateDomain(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func Test_GetDomain(t *testing.T) {

}

func Test_GetDomainByName(t *testing.T) {

}

func Test_ListDomain(t *testing.T) {

}

func Test_UpdateDomain(t *testing.T) {

}

func Test_DeleteDomain(t *testing.T) {

}

func Test_CheckDomainIsExistByID(t *testing.T) {

}

func Test_CheckDomainIsExistByName(t *testing.T) {

}
