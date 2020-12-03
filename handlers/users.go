package handlers

import (
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
)

type Users struct {
	log hclog.Logger
	db  *sqlx.DB
}

func NewUsers(l hclog.Logger, db *sqlx.DB) *Users {
	return &Users{l, db}
}

func (u *Users) Insert(rw http.ResponseWriter, r *http.Request) {

}
