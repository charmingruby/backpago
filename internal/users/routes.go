package users

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	db *sql.DB
}

func SetRoutes(r chi.Router, db *sql.DB) {
	h := handler{db}

	r.Get("/", h.List)
	r.Get("/{id}", h.GetById)
	r.Post("/sign-up", h.Create)
	r.Put("/{id}", h.Modify)
	r.Put("/{id}", h.Delete)
}
