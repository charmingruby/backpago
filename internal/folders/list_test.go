package folders

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestList() {
	tcs := []struct {
		WithMock         bool
		MockListErr      bool
		MockFilesErr     bool
		ExpectStatusCode int
	}{
		// success
		{true, false, false, http.StatusOK},
		// errors
		{false, true, true, http.StatusInternalServerError},
		{true, true, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		if tc.WithMock {
			setMockList(ts.mock, tc.MockListErr)
			if !tc.MockListErr {
				setMockListRootFiles(ts.mock, tc.MockFilesErr)
			}
		}

		ts.handler.List(rr, req)
		assert.Equal(ts.T(), tc.ExpectStatusCode, rr.Code)
	}
}

func (ts *TransactionSuite) TestGetRootSubFolders() {
	setMockList(ts.mock, false)

	_, err := getRootSubFolders(ts.conn)
	assert.NoError(ts.T(), err)
}

func setMockList(mock sqlmock.Sqlmock, err bool) {
	exp := mock.ExpectQuery(regexp.QuoteMeta(`select * from "folders" where "parent_id" is null and "deleted"=false`))

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "modified_at", "deleted"}).
			AddRow(1, 0, "Documentos", time.Now(), time.Now(), false).
			AddRow(5, 0, "Imagens", time.Now(), time.Now(), false)

		exp.WillReturnRows(rows)
	}
}

func setMockListRootFiles(mock sqlmock.Sqlmock, err bool) {
	exp := mock.ExpectQuery(regexp.QuoteMeta(`select * from files where "folder_id" is null and "deleted"=false`))

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "modified_at", "deleted"}).
			AddRow(1, 0, 1, "Gopher.png", "image/png", "/", time.Now(), time.Now(), false).
			AddRow(2, 0, 1, "Golang-LOGO.png", "image/jpg", "/", time.Now(), time.Now(), false)

		exp.WillReturnRows(rows)
	}
}
