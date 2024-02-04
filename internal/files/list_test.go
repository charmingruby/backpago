package files

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id",
		"folder_id",
		"owner_id",
		"name",
		"type",
		"path",
		"created_at",
		"modified_at",
		"deleted",
	}).
		AddRow(1, 1, 1, "photo", ".png", "./", now, now, false).
		AddRow(2, 1, 1, "photo", ".png", "./", now, now, false)

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "files" where folder_id=$1 and deleted = false`)).
		WillReturnRows(rows)

	_, err = List(db, 1)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func TestListRoot(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id",
		"folder_id",
		"owner_id",
		"type",
		"name",
		"path",
		"created_at",
		"modified_at",
		"deleted",
	}).
		AddRow(1, nil, 1, "photo", ".png", "./", now, now, false).
		AddRow(2, nil, 1, "photo", ".png", "./", now, now, false)

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "files" where folder_id is null and deleted = false`)).
		WillReturnRows(rows)

	_, err = ListRoot(db)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
