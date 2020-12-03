package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nicholasjackson/testing-microservices/handlers"
)

func main() {

	l := hclog.Default()

	db, err := sqlx.Connect("postgres", "user=root dbname=root sslmode=disable password=password")
	if err != nil {
		l.Error("Unable to connect to the DB", "error", err)
		os.Exit(1)
	}

	bh := handlers.NewBranches(l, db)
	uh := handlers.NewUsers(l, db)

	r := mux.NewRouter()
	r.HandleFunc("/branches", bh.Get).Methods(http.MethodGet)
	r.HandleFunc("/branches", bh.Insert).Methods(http.MethodPost)

	r.HandleFunc("/users", uh.Insert).Methods(http.MethodPost)

	http.Handle("/", r)
	http.ListenAndServe(":9090", nil)

	l.Info("Server listening on :9090")
}
