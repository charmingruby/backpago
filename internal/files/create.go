package files

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/charmingruby/backpago/internal/queue"
)

func (h *handler) Create(rw http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	path := fmt.Sprintf("/%s", fileHeader.Filename)

	err = h.bucket.Upload(file, path)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	entity, err := New(
		1,
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
		path,
	)
	if err != nil {
		h.bucket.Delete(path)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	folderId := r.Form.Get("folder_id")
	if folderId != "" {
		fid, err := strconv.Atoi(folderId)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		entity.FolderId = int64(fid)
	}

	id, err := Insert(h.db, entity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	entity.Id = id

	dto := queue.QueueDto{
		ID:       int(id),
		Filename: fileHeader.Filename,
		Path:     path,
	}

	msg, err := dto.Marshal()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.queue.Publish(msg)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(entity)
}

func Insert(db *sql.DB, f *File) (int64, error) {
	stmt := `insert into "files" 
	("folder_id", "owner_id", "name", "type", "path", "modified_at") 
		VALUES ($1, $2, $3, $4, $5, $6)`

	result, err := db.Exec(stmt, f.FolderId, f.OwnerId, f.Name, f.Type, f.Path, f.ModifiedAt)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
