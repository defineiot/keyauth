package mysql_test

import (
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateDomain(t *testing.T) {
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

func TestGetDomain(t *testing.T) {

}

func TestGetDomainByName(t *testing.T) {

}

func TestListDomain(t *testing.T) {

}

func TestUpdateDomain(t *testing.T) {

}

func TestDeleteDomain(t *testing.T) {

}

func TestCheckDomainIsExistByID(t *testing.T) {

}

func TestCheckDomainIsExistByName(t *testing.T) {

}
