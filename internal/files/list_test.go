package files

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestList() {
	tcs := []struct {
		FolderID    int64
		MockWithErr bool
	}{
		{1, false},
		{300, true},
	}

	for _, tc := range tcs {
		setMockList(ts.mock, tc.FolderID, tc.MockWithErr)

		_, err := List(ts.conn, tc.FolderID)
		if tc.MockWithErr {
			assert.Error(ts.T(), err)
		} else {
			assert.NoError(ts.T(), err)
		}
	}
}

func setMockList(mock sqlmock.Sqlmock, folderID int64, err bool) {
	exp := mock.ExpectQuery(regexp.QuoteMeta(`select * from files where "folder_id" = $1 and "deleted"=false`)).
		WithArgs(folderID)

	if err {
		exp.WillReturnError(sql.ErrConnDone)
	} else {
		rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "modified_at", "deleted"}).
			AddRow(1, 1, 1, "Gopher.png", "image/png", "/", time.Now(), time.Now(), false).
			AddRow(2, 1, 1, "Golang-LOGO.png", "image/jpg", "/", time.Now(), time.Now(), false)

		exp.WillReturnRows(rows)
	}
}

func (ts *TransactionSuite) TestListRoot() {
	tcs := []struct {
		MockWithErr bool
	}{
		{false},
		{true},
	}

	for _, tc := range tcs {
		setMockListRoot(ts.mock, tc.MockWithErr)

		_, err := ListRoot(ts.conn)
		if tc.MockWithErr {
			assert.Error(ts.T(), err)
		} else {
			assert.NoError(ts.T(), err)
		}
	}
}

func setMockListRoot(mock sqlmock.Sqlmock, err bool) {
	exp := mock.ExpectQuery(regexp.QuoteMeta(`select * from files where "folder_id" is null and "deleted"=false`))

	if err {
		exp.WillReturnError(sql.ErrConnDone)
	} else {
		rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "modified_at", "deleted"}).
			AddRow(1, nil, 1, "Gopher.png", "image/png", "/", time.Now(), time.Now(), false).
			AddRow(2, nil, 1, "Golang-LOGO.png", "image/jpg", "/", time.Now(), time.Now(), false)

		exp.WillReturnRows(rows)
	}
}
