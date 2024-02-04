package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{db}

	u, _ := New("JohnDoe", "john@email.com", "password")

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(u)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", &b)

	mock.ExpectExec(`insert into "users" ("name", "login", "password", "modified_at")*`).
		WithArgs(u.Name, u.Login, u.Password, u.ModifiedAt.Truncate(0)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h.Create(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Error: %v", rr)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	u, err := New("John Doe", "john@doe.com", "password")
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(`insert into "users" ("name", "login", "password", "modified_at")*`).
		WithArgs("John Doe", "john@doe.com", u.Password, u.ModifiedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = Insert(db, u)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
