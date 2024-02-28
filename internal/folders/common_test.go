package folders

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type TransactionSuite struct {
	suite.Suite
	conn *sql.DB
	mock sqlmock.Sqlmock

	handler handler
	entity  *Folder
}

func (ts *TransactionSuite) SetupTest() {
	var err error

	ts.conn, ts.mock, err = sqlmock.New()
	assert.NoError(ts.T(), err)

	ts.handler = handler{ts.conn}

	ts.entity = &Folder{
		ID:   1,
		Name: "Fotos",
	}
}

func (ts *TransactionSuite) AfterTest(_, _ string) {
	assert.NoError(ts.T(), ts.mock.ExpectationsWereMet())
}

func setMockListFiles(mock sqlmock.Sqlmock, id int64, err bool, withRows bool) {
	exp := mock.ExpectQuery(regexp.QuoteMeta(`select * from files where "folder_id" = $1 and "deleted"=false`)).
		WithArgs(id)

	if err {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		rows := sqlmock.NewRows([]string{"id", "folder_id", "owner_id", "name", "type", "path", "created_at", "modified_at", "deleted"})
		if withRows {
			rows.AddRow(1, id, 1, "Gopher.png", "image/png", "/", time.Now(), time.Now(), false).
				AddRow(2, id, 1, "Golang-LOGO.png", "image/jpg", "/", time.Now(), time.Now(), false)
		}

		exp.WillReturnRows(rows)
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}
