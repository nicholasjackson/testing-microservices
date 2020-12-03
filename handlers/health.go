package handlers

import (
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
)

type Health struct {
	log hclog.Logger
	db  *sqlx.DB
}

func NewHealth(l hclog.Logger, db *sqlx.DB) *Health {
	return &Health{l, db}
}

func (b *Health) Get(rw http.ResponseWriter, r *http.Request) {
	// check the DB is ok
	err := b.db.Ping()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		return
	}
}
