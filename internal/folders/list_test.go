package folders

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetRootSubFolders(t *testing.T) {
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
		AddRow(2, nil, "cvs", time.Now(), time.Now(), false).
		AddRow(4, nil, "imgs", time.Now(), time.Now(), false)

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "folders" where "parent_id" is null and "deleted"=false`)).
		WillReturnRows(rows)

	_, err = getRootSubFolders(db)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
