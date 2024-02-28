package folders

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestDeleteHTTP() {
	tcs := []struct {
		ID                string
		WithMock          bool
		MockID            int64
		MockListErr       bool
		MockDeleteFileErr bool
		MockSubfolderErr  bool
		MockWithErr       bool
		ExpectStatusCode  int
	}{
		// success
		{"1", true, 1, false, false, false, false, http.StatusNoContent},
		// errors
		{"A", false, -1, true, true, true, true, http.StatusInternalServerError},
		{"10", true, 10, true, true, true, true, http.StatusInternalServerError},
		{"15", true, 15, false, true, true, true, http.StatusInternalServerError},
		{"20", true, 20, false, false, true, true, http.StatusInternalServerError},
		{"25", true, 25, false, false, false, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/{id}", nil)

		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.ID)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMock {
			// start dependencies queries
			setMockListFiles(ts.mock, tc.MockID, tc.MockListErr, true)
			if !tc.MockListErr {
				setMockDeleteFile(ts.mock, "Gopher.png", 1, tc.MockDeleteFileErr)
				if !tc.MockDeleteFileErr {
					setMockDeleteFile(ts.mock, "Golang-LOGO.png", 2, tc.MockDeleteFileErr)
				}

				if !tc.MockDeleteFileErr {
					setMockGetSubFolder(ts.mock, tc.MockID, tc.MockSubfolderErr)
					// finish dependencies queries

					if !tc.MockSubfolderErr {
						setMockDelete(ts.mock, tc.MockID, tc.MockWithErr)
					}
				}
			}
		}

		ts.handler.Delete(rr, req)
		assert.Equal(ts.T(), tc.ExpectStatusCode, rr.Code)
	}
}

func (ts *TransactionSuite) TestDelete() {
	setMockDelete(ts.mock, 1, false)

	err := Delete(ts.conn, 1)
	assert.NoError(ts.T(), err)
}

func setMockDelete(mock sqlmock.Sqlmock, id int64, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`update "folders" set "modified_at"=$1, "deleted"=true where id=$2`)).
		WithArgs(AnyTime{}, id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func setMockDeleteFile(mock sqlmock.Sqlmock, filename string, id int64, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`update "files" set "name"=$1, "modified_at"=$2, "deleted"=$3 where id=$4`)).
		WithArgs(filename, AnyTime{}, true, id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))

	}
}
