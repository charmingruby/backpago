package folders

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestModify() {
	tcs := []struct {
		ID                string
		WithMock          bool
		MockFolder        *Folder
		MockUpdateWithErr bool
		MockGetWithErr    bool
		ExpectStatusCode  int
	}{
		// success
		{"1", true, ts.entity, false, false, http.StatusOK},
		// errors
		{"5", false, nil, true, true, http.StatusInternalServerError},
		{"10", false, &Folder{}, true, true, http.StatusBadRequest},
		{"A", false, ts.entity, true, true, http.StatusInternalServerError},
		{"25", true, &Folder{ID: 25, Name: "Documentos"}, true, true, http.StatusInternalServerError},
		{"500", true, &Folder{ID: 500, Name: "Documentos"}, false, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		var b bytes.Buffer
		if tc.MockFolder != nil {
			err := json.NewEncoder(&b).Encode(tc.MockFolder)
			assert.NoError(ts.T(), err)
		} else {
			b.Write([]byte(`{"name": false}`))
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/{id}", &b)

		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.ID)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMock {
			setMockUpdate(ts.mock, tc.MockFolder.ID, tc.MockFolder.Name, tc.MockUpdateWithErr)
			// caso o mock para update retornar erro, não podemos fazer o mock do get
			if !tc.MockUpdateWithErr {
				setMockGet(ts.mock, tc.MockFolder.ID, tc.MockGetWithErr)
			}
		}

		ts.handler.Modify(rr, req)
		assert.Equal(ts.T(), tc.ExpectStatusCode, rr.Code)
	}
}

func (ts *TransactionSuite) TestUpdate() {
	setMockUpdate(ts.mock, 1, "Documentos", false)

	err := Update(ts.conn, 1, &Folder{Name: "Documentos"})
	assert.NoError(ts.T(), err)
}

func setMockUpdate(mock sqlmock.Sqlmock, id int64, name string, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`update "folders" set name = $1, modified_at = $2 where id = $3`)).
		WithArgs(name, AnyTime{}, id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
