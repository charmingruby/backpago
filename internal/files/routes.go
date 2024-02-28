package files

import (
	"database/sql"

	"github.com/charmingruby/backpago/internal/auth"
	"github.com/charmingruby/backpago/internal/bucket"
	"github.com/charmingruby/backpago/internal/queue"
	"github.com/go-chi/chi"
)

type handler struct {
	db     *sql.DB
	bucket *bucket.Bucket
	queue  *queue.Queue
}

func SetRoutes(r chi.Router, db *sql.DB, b *bucket.Bucket, q *queue.Queue) {
	h := handler{db, b, q}

	r.Route("/files", func(r chi.Router) {
		r.Use(auth.Validate)

		r.Post("/", h.Create)
		r.Put("/{id}", h.Modify)
		r.Delete("/{id}", h.Delete)
	})
}
