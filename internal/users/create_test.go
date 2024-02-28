package users

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestCreate() {
	tcs := []struct {
		WithMock         bool
		User             *User
		MockWithErr      bool
		ExpectStatusCode int
	}{
		// success
		{true, ts.entity, false, http.StatusCreated},
		// errors
		{false, nil, true, http.StatusInternalServerError},
		{false, &User{Name: "Tiago Temporin"}, true, http.StatusBadRequest},
		{false, &User{Name: "Tiago Temporin", Password: "123"}, true, http.StatusBadRequest},
		{false, &User{Password: "1234567"}, true, http.StatusBadRequest},
		{true, ts.entity, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		var b bytes.Buffer
		if tc.User != nil {
			err := json.NewEncoder(&b).Encode(tc.User)
			assert.NoError(ts.T(), err)
		} else {
			b.Write([]byte(`{"name": false}`))
		}

		ts.entity.SetPassword(ts.entity.Password)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", &b)

		if tc.WithMock {
			setMockInsert(ts.mock, tc.User, tc.MockWithErr)
		}

		ts.handler.Create(rr, req)
		assert.Equal(ts.T(), tc.ExpectStatusCode, rr.Code)
	}
}

func (ts *TransactionSuite) TestInsert() {
	setMockInsert(ts.mock, ts.entity, false)

	_, err := Insert(ts.conn, ts.entity)
	assert.NoError(ts.T(), err)
}

func setMockInsert(mock sqlmock.Sqlmock, entity *User, err bool) {
	exp := mock.ExpectExec(`insert into "users" ("name", "login", "password", "modified_at")*`).
		WithArgs(entity.Name, entity.Login, entity.Password, entity.ModifiedAt)

	if err {
		exp.WillReturnError(sql.ErrConnDone)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
