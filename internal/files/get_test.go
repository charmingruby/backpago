package files

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
		"folder_id",
		"owner_id",
		"name",
		"type",
		"path",
		"created_at",
		"modified_at",
		"deleted",
	}).
		AddRow(1, "1", "1", "photo", ".png", "./", time.Now(), time.Now(), false)

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "files" where id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = Get(db, 1)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
