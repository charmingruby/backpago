package folders

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetFolder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"parent_id",
		"name",
		"created_at",
		"modified_at",
		"deleted",
	}).
		AddRow(1, 2, "docs", time.Now(), time.Now(), false)

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "folders" where id=$1`)).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = GetFolder(db, 1)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func TestGetSubFolder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"parent_id",
		"name",
		"created_at",
		"modified_at",
		"deleted",
	}).
		AddRow(2, 3, "cvs", time.Now(), time.Now(), false).
		AddRow(4, 3, "imgs", time.Now(), time.Now(), false)

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "folders" where "parent_id"=$1 and "deleted"=false`)).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = getSubFolders(db, 1)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
