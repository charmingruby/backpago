package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) Create(rw http.ResponseWriter, r *http.Request) {
	f := new(Folder)

	if err := json.NewDecoder(r.Body).Decode(f); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := f.Validate(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := Insert(h.db, f)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	f.Id = id

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(f)

}

func Insert(db *sql.DB, f *Folder) (int64, error) {
	stmt := `insert into "folders" ("parent_id", "name", "modified_at") values ($1, $2, $3)`

	result, err := db.Exec(stmt)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()

}
