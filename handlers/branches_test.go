package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*Branches, *sqlx.DB) {
	db, _ := sqlx.Connect("postgres", "user=root dbname=root sslmode=disable password=password")
	bh := NewBranches(hclog.Default(), db)

	return bh, db
}

func Test_ReturnsBranchesAsJSON(t *testing.T) {
	bh, _ := setup(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	bh.Get(rr, r)
	require.Equal(t, rr.Code, http.StatusOK)

	branches := []Branch{}

	err := json.NewDecoder(rr.Body).Decode(&branches)
	require.NoError(t, err)

	require.Len(t, branches, 7)
}

func Test_InsertBranchReturnsOK(t *testing.T) {
	bh, db := setup(t)

	branch := Branch{
		Name:   "unit test name",
		Street: "test street",
		City:   "test city",
		Zip:    "12345",
	}

	d, _ := json.Marshal(&branch)

	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(d))
	rr := httptest.NewRecorder()

	bh.Insert(rr, r)
	require.Equal(t, rr.Code, http.StatusOK)

	// check the record has been written
	checkBranch := Branch{}

	err := db.Get(&checkBranch, `SELECT * FROM branches WHERE name=$1`, branch.Name)
	require.NoError(t, err)
	require.Equal(t, branch.Name, checkBranch.Name)
	require.Equal(t, branch.Street, checkBranch.Street)
	require.Equal(t, branch.City, checkBranch.City)
	require.Equal(t, branch.Zip, checkBranch.Zip)
}
