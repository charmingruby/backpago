package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *handler) Modify(rw http.ResponseWriter, r *http.Request) {
	f := new(Folder)

	if err := json.NewDecoder(r.Body).Decode(f); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := f.Validate(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := Update(h.db, int64(id), f); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: GetFolder

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(f)
}

func Update(db *sql.DB, id int64, f *Folder) error {
	f.ModifiedAt = time.Now()

	stmt := `update "folders" set name = $1, modified_at = $2 where id = $3`
	_, err := db.Exec(stmt, f.Name, f.ModifiedAt, id)
	if err != nil {
		return err
	}

	return nil
}
