package users

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"login",
		"password",
		"created_at",
		"modified_at",
		"deleted",
		"last_login",
	}).
		AddRow(1, "john doe", "john_doe", "john1234", time.Now(), time.Now(), false, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where id=$1`)).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = GetById(db, 1)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
