package users

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestGetByID() {
	tcs := []struct {
		ID               string
		WithMock         bool
		MockID           int64
		MockWithErr      bool
		ExpectStatusCode int
	}{
		// success
		{"1", true, 1, false, http.StatusOK},
		// errors
		{"A", false, -1, true, http.StatusInternalServerError},
		{"25", true, 25, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/{id}", nil)

		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.ID)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMock {
			setMockGet(ts.mock, tc.MockID, tc.MockWithErr)
		}

		ts.handler.GetByID(rr, req)
		assert.Equal(ts.T(), tc.ExpectStatusCode, rr.Code)
	}
}

func (ts *TransactionSuite) TestGet() {
	setMockGet(ts.mock, 1, false)

	_, err := Get(ts.conn, 1)
	assert.NoError(ts.T(), err)
}

func setMockGet(mock sqlmock.Sqlmock, id int64, err bool) {
	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(id, "Tiago", "tiago@aprendagolang.com.br", "123456", time.Now(), time.Now(), false, time.Now())

	exp := mock.ExpectQuery(regexp.QuoteMeta(`select * from "users" where id=$1`)).
		WithArgs(id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnRows(rows)
	}
}
