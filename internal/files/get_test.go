package files

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestGet() {
	tcs := []struct {
		ID          int64
		MockWithErr bool
	}{
		{1, false},
		{300, true},
	}

	for _, tc := range tcs {
		setMockGet(ts.mock, tc.ID, tc.MockWithErr)

		_, err := Get(ts.conn, tc.ID)
		if tc.MockWithErr {
			assert.Error(ts.T(), err)
		} else {
			assert.NoError(ts.T(), err)
		}
	}
}

func setMockGet(mock sqlmock.Sqlmock, id int64, err bool) {
	rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "modified_at", "deleted"})

	if err {
		rows.AddRow(1, 1, 1, "Gopher.png", "image/png", "/", time.Now(), time.Now(), 10)
	} else {
		rows.AddRow(1, 1, 1, "Gopher.png", "image/png", "/", time.Now(), time.Now(), false)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select * from "files" where id = $1`)).
		WithArgs(id).
		WillReturnRows(rows)
}
