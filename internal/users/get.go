package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *handler) GetById(rw http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	u, err := GetById(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Content-Type", "application/")
	json.NewEncoder(rw).Encode(u)
}

func GetById(db *sql.DB, id int64) (*User, error) {
	stmt := `select * from "users" where id=$1`
	row := db.QueryRow(stmt, id)

	var u User
	err := row.Scan(&u.Id, &u.Name, &u.Login, &u.Password, &u.CreatedAt, &u.ModifiedAt, &u.Deleted, &u.LastLogin)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
