package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nicholasjackson/testing-microservices/handlers"
)

func main() {

	dbPort := "5432"
	if p := os.Getenv("DB_PORT"); p != "" {
		dbPort = p
	}

	l := hclog.Default()

	var db *sqlx.DB
	var err error

	for n := 0; n < 30; n++ {
		l.Info("Attempting to connect to the DB", "attempt", n+1)
		db, err = sqlx.Connect(
			"postgres",
			fmt.Sprintf("user=root dbname=root sslmode=disable password=password port=%s", dbPort),
		)
		if err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if err != nil {
		l.Error("Unable to connect to the DB", "error", err)
		os.Exit(1)
	}

	bh := handlers.NewBranches(l, db)
	uh := handlers.NewUsers(l, db)
	hh := handlers.NewHealth(l, db)

	r := mux.NewRouter()
	r.HandleFunc("/branches", bh.Get).Methods(http.MethodGet)
	r.HandleFunc("/branches", bh.Insert).Methods(http.MethodPost)

	r.HandleFunc("/users", uh.Insert).Methods(http.MethodPost)

	r.HandleFunc("/health", hh.Get).Methods(http.MethodGet)

	http.Handle("/", r)

	l.Info("Starting server on :9090")

	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		l.Error("Unable to start server", "error", err)
		os.Exit(1)
	}
}
