package folders

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

func (ts *TransactionSuite) TestGet() {
	tcs := []struct {
		ID               string
		WithMock         bool
		MockID           int64
		MockSubfolderErr bool
		MockListErr      bool
		MockWithErr      bool
		ExpectStatusCode int
	}{
		// success
		{"1", true, 1, false, false, false, http.StatusOK},
		// errors
		{"A", false, -1, true, true, true, http.StatusInternalServerError},
		{"10", true, 10, true, true, false, http.StatusInternalServerError},
		{"15", true, 15, false, true, false, http.StatusInternalServerError},
		{"20", true, 20, false, false, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/{id}", nil)

		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.ID)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMock {
			setMockGet(ts.mock, tc.MockID, tc.MockWithErr)
			// caso o mock para get retornar erro, não podemos fazer o mock do content
			if !tc.MockWithErr {
				setMockGetSubFolder(ts.mock, tc.MockID, tc.MockSubfolderErr)

				if !tc.MockSubfolderErr {
					setMockListFiles(ts.mock, tc.MockID, tc.MockListErr, true)
				}
			}
		}

		ts.handler.Get(rr, req)
		assert.Equal(ts.T(), tc.ExpectStatusCode, rr.Code)
	}
}
func (ts *TransactionSuite) TestGetFolder() {
	setMockGet(ts.mock, 1, false)

	_, err := GetFolder(ts.conn, 1)
	assert.NoError(ts.T(), err)
}

func (ts *TransactionSuite) TestGetSubFolder() {
	setMockGetSubFolder(ts.mock, 1, false)

	_, err := getSubFolders(ts.conn, 1)
	assert.NoError(ts.T(), err)
}

func setMockGet(mock sqlmock.Sqlmock, id int64, err bool) {
	exp := mock.ExpectQuery(regexp.QuoteMeta(`select * from "folders" where id=$1`)).
		WithArgs(id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "modified_at", "deleted"}).
			AddRow(id, 2, "Documentos", time.Now(), time.Now(), false)

		exp.WillReturnRows(rows)
	}
}

func setMockGetSubFolder(mock sqlmock.Sqlmock, id int64, err bool) {
	exp := mock.ExpectQuery(regexp.QuoteMeta(`select * from "folders" where "parent_id"=$1 and "deleted"=false`)).
		WithArgs(id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		rows := sqlmock.NewRows([]string{"id", "parent_id", "name", "created_at", "modified_at", "deleted"}).
			AddRow(2, id, "Projetos Pessoais", time.Now(), time.Now(), false).
			AddRow(4, id, "Projetos Em Geral", time.Now(), time.Now(), false)

		exp.WillReturnRows(rows)
	}
}
