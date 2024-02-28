package users

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
		MockWithErr      bool
		ExpectStatusCode int
	}{
		// success
		{true, false, http.StatusOK},
		// errors
		{false, true, http.StatusInternalServerError},
		{true, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		if tc.WithMock {
			setMockList(ts.mock, tc.MockWithErr)
		}

		ts.handler.List(rr, req)
		assert.Equal(ts.T(), tc.ExpectStatusCode, rr.Code)
	}
}

func (ts *TransactionSuite) TestSelectAll() {
	setMockList(ts.mock, false)

	_, err := SelectAll(ts.conn)
	assert.NoError(ts.T(), err)
}

func setMockList(mock sqlmock.Sqlmock, err bool) {
	exp := mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where deleted = false`))

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
			AddRow(1, "Tiago", "tiago@aprendagolang.com.br", "123456", time.Now(), time.Now(), false, time.Now()).
			AddRow(2, "Maria", "maria@exemplo.com", "123456", time.Now(), time.Now(), false, time.Now())

		exp.WillReturnRows(rows)
	}
}
