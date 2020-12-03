package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
)

type Branches struct {
	log hclog.Logger
	db  *sqlx.DB
}

func NewBranches(l hclog.Logger, db *sqlx.DB) *Branches {
	return &Branches{l, db}
}

func (b *Branches) Get(rw http.ResponseWriter, r *http.Request) {
	b.log.Info("Get handler called")

	branches := []Branch{}
	err := b.db.Select(&branches, `SELECT * FROM branches`)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	je := json.NewEncoder(rw)

	je.Encode(branches)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (b *Branches) Insert(rw http.ResponseWriter, r *http.Request) {
	b.log.Info("Insert handler caller")

	branch := Branch{}

	err := json.NewDecoder(r.Body).Decode(&branch)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = b.db.NamedExec("INSERT INTO branches (name, street, city, zip) VALUES (:name, :street, :city, :zip)", branch)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Branch struct {
	ID     string `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Street string `db:"street" json:"street"`
	City   string `db:"city" json:"city"`
	Zip    string `db:"zip" json:"zip"`
}
