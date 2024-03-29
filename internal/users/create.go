package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) Create(rw http.ResponseWriter, r *http.Request) {
	u := new(User)

	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = u.SetPassword(u.Password)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = u.Validate()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := Insert(h.db, u)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	u.ID = id

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(u)
}

func Insert(db *sql.DB, u *User) (id int64, err error) {
	stmt := `insert into "users" ("name", "login", "password", "modified_at") VALUES ($1, $2, $3, $4) RETURNING id`
	err = db.QueryRow(stmt, u.Name, u.Login, u.Password, u.ModifiedAt).Scan(&id)
	if err != nil {
		return -1, err
	}

	return
}
