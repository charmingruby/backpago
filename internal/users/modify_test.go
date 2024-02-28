package users

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
		MockUser          *User
		MockUpdateWithErr bool
		MockGetWithErr    bool
		ExpectStatusCode  int
	}{
		// success
		{"1", true, &User{ID: 1, Name: "Tiago Temporin"}, false, false, http.StatusOK},
		// errors
		{"5", false, nil, true, true, http.StatusInternalServerError},
		{"10", false, &User{ID: 10}, true, false, http.StatusBadRequest},
		{"A", false, &User{Name: "Tiago Temporin"}, true, false, http.StatusInternalServerError},
		{"25", true, &User{ID: 25, Name: "Tiago Temporin"}, true, false, http.StatusInternalServerError},
		{"500", true, &User{ID: 500, Name: "Tiago Temporin"}, false, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		var b bytes.Buffer
		if tc.MockUser != nil {
			err := json.NewEncoder(&b).Encode(tc.MockUser)
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
			setMockUpdate(ts.mock, tc.MockUser.ID, tc.MockUser.Name, tc.MockUpdateWithErr)
			// caso o mock para update retornar erro, não podemos fazer o mock do get
			if !tc.MockUpdateWithErr {
				setMockGet(ts.mock, tc.MockUser.ID, tc.MockGetWithErr)
			}
		}

		ts.handler.Modify(rr, req)
		assert.Equal(ts.T(), tc.ExpectStatusCode, rr.Code)
	}
}

func (ts *TransactionSuite) TestUpdate() {
	setMockUpdate(ts.mock, 1, "Tiago Temporin", false)

	err := Update(ts.conn, 1, &User{Name: "Tiago Temporin"})
	assert.NoError(ts.T(), err)
}

func setMockUpdate(mock sqlmock.Sqlmock, id int64, name string, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`update "users" set "name"=$1, "modified_at"=$2 where id=$3`)).
		WithArgs(name, AnyTime{}, id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
