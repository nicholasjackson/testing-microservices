package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*httptest.Server, *sqlx.DB) {
	db, _ := sqlx.Connect("postgres", "user=root dbname=root sslmode=disable password=password")
	bh := NewBranches(hclog.Default(), db)

	r := mux.NewRouter()
	r.HandleFunc("/", bh.Get).Methods(http.MethodGet)
	r.HandleFunc("/", bh.Insert).Methods(http.MethodPost)

	ts := httptest.NewServer(r)

	t.Cleanup(func() {
		ts.Close()
	})

	return ts, db
}

func Test_ReturnsBranchesAsJSON(t *testing.T) {
	ts, _ := setup(t)

	resp, err := http.Get(ts.URL)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)

	branches := []Branch{}

	err = json.NewDecoder(resp.Body).Decode(&branches)
	require.NoError(t, err)

	require.Len(t, branches, 7)
}

func Test_InsertBranchReturnsOK(t *testing.T) {
	ts, db := setup(t)

	branch := Branch{
		Name:   "unit test name",
		Street: "test street",
		City:   "test city",
		Zip:    "12345",
	}

	d, _ := json.Marshal(&branch)

	resp, err := http.Post(ts.URL, "application/json", bytes.NewReader(d))
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)

	// check the record has been written
	checkBranch := Branch{}

	err = db.Get(&checkBranch, `SELECT * FROM branches WHERE name=$1`, branch.Name)
	require.NoError(t, err)
	require.Equal(t, branch.Name, checkBranch.Name)
	require.Equal(t, branch.Street, checkBranch.Street)
	require.Equal(t, branch.City, checkBranch.City)
	require.Equal(t, branch.Zip, checkBranch.Zip)
}
