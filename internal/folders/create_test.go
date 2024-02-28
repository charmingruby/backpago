package folders

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestCreate() {
	tcs := []struct {
		WithMock         bool
		Folder           *Folder
		MockWithErr      bool
		ExpectStatusCode int
	}{
		// success
		{true, ts.entity, false, http.StatusCreated},
		// errors
		{false, nil, true, http.StatusInternalServerError},
		{false, &Folder{}, true, http.StatusBadRequest},
		{true, ts.entity, true, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		var b bytes.Buffer
		if tc.Folder != nil {
			err := json.NewEncoder(&b).Encode(tc.Folder)
			assert.NoError(ts.T(), err)
		} else {
			b.Write([]byte(`{"name": false}`))
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", &b)

		if tc.WithMock {
			setMockInsert(ts.mock, tc.Folder, tc.MockWithErr)
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

func setMockInsert(mock sqlmock.Sqlmock, entity *Folder, err bool) {
	exp := mock.ExpectExec(regexp.QuoteMeta(`insert into "folders" ("parent_id", "name", "modified") values ($1, $2, $3)`)).
		WithArgs(entity.ParentID, entity.Name, AnyTime{})

	if err {
		exp.WillReturnError(sql.ErrConnDone)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
