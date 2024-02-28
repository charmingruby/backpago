package files

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
		MockID            int64
		MockFile          *File
		MockGetWithErr    bool
		MockUpdateWithErr bool
		ExpectStatusCode  int
	}{
		// success
		{"1", true, 1, &File{ID: 1, Name: "aprenda-golang.png"}, false, false, http.StatusOK},
		// errors
		{"5", true, 5, nil, false, true, http.StatusInternalServerError},
		{"10", true, 10, &File{ID: 10}, false, true, http.StatusBadRequest},
		{"A", false, 0, &File{Name: "wrong-id.jpg"}, false, true, http.StatusInternalServerError},
		{"25", true, 25, &File{ID: 25, Name: "aprenda-golang.png"}, true, false, http.StatusInternalServerError},
		{"500", true, 500, &File{ID: 500, Name: "aprenda-golang.png"}, false, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		var b bytes.Buffer
		if tc.MockFile != nil {
			err := json.NewEncoder(&b).Encode(tc.MockFile)
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
			setMockGet(ts.mock, tc.MockID, tc.MockGetWithErr)
			if !tc.MockGetWithErr && tc.MockFile != nil && tc.MockFile.Name != "" {
				setMockUpdate(ts.mock, tc.MockFile.Name, tc.MockID, tc.MockUpdateWithErr)
			}
		}

		ts.handler.Modify(rr, req)
		assert.Equal(ts.T(), tc.ExpectStatusCode, rr.Code)
	}
}

func (ts *TransactionSuite) TestUpdate() {
	setMockUpdate(ts.mock, "aprenda-golang.png", 1, false)

	err := Update(ts.conn, 1, &File{Name: "aprenda-golang.png"})
	assert.NoError(ts.T(), err)
}

func setMockUpdate(mock sqlmock.Sqlmock, name string, id int64, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`update "files" set "name"=$1, "modified_at"=$2, "deleted"=$3 where id=$4`)).
		WithArgs(name, AnyTime{}, false, id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
