package files

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestCreate() {
	tcs := []struct {
		WithMock         bool
		WithUpload       bool
		File             *File
		FolderID         int64
		MockFolderID     string
		MockWithErr      bool
		ExpectStatusCode int
	}{
		// success
		{true, true, ts.entity, 0, "", false, http.StatusCreated},
		{true, true, ts.entity, 15, "15", false, http.StatusCreated},
		// errors
		{false, false, nil, 0, "", true, http.StatusInternalServerError},
		{false, true, &File{Name: ".empty"}, 0, "", true, http.StatusInternalServerError},
		{false, true, ts.entity, 0, "A", false, http.StatusInternalServerError},
	}

	for _, tc := range tcs {
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)

		// START UPLOAD
		if tc.WithUpload {
			file, err := os.Open(fmt.Sprintf("./testdata/%s", tc.File.Name))
			assert.NoError(ts.T(), err)

			w, err := mw.CreateFormFile("file", tc.File.Name)
			assert.NoError(ts.T(), err)

			_, err = io.Copy(w, file)
			assert.NoError(ts.T(), err)
		}

		if tc.MockFolderID != "" {
			w, err := mw.CreateFormField("folder_id")
			assert.NoError(ts.T(), err)

			w.Write([]byte(tc.MockFolderID))
		}

		mw.Close()
		// END UPLOAD

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Add("Content-Type", mw.FormDataContentType())

		if tc.WithMock {
			setMockInsert(ts.mock, tc.File, tc.FolderID)
		}

		ts.handler.Create(rr, req)
		assert.Equal(ts.T(), tc.ExpectStatusCode, rr.Code)
	}
}

func (ts *TransactionSuite) TestInsert() {
	setMockInsert(ts.mock, ts.entity, 0)

	_, err := Insert(ts.conn, ts.entity)
	assert.NoError(ts.T(), err)
}

func setMockInsert(mock sqlmock.Sqlmock, entity *File, folderID int64) {
	mock.ExpectExec(regexp.QuoteMeta(`insert into "files" ("folder_id", "owner_id", "name", "type", "path", "modified_at") VALUES ($1, $2, $3, $4, $5, $6)`)).
		WithArgs(folderID, entity.OwnerID, entity.Name, entity.Type, entity.Path, AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
