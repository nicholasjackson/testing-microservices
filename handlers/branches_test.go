package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-hclog"
	_ "github.com/lib/pq"
	"github.com/nicholasjackson/testing-microservices/data"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*Branches, *data.DBMock) {
	//db, _ := sqlx.Connect("postgres", "user=root dbname=root sslmode=disable password=password")
	db := &data.DBMock{}
	bh := NewBranches(hclog.Default(), db)

	return bh, db
}

func Test_ReturnsBranchesAsJSON(t *testing.T) {
	bh, db := setup(t)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	// setup expectations for the mock
	db.On("Select", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			// sqlx mutates the object passed in
			a1 := args.Get(0).(*[]Branch)

			// setup return obj
			branches := []Branch{
				Branch{},
				Branch{},
				Branch{},
				Branch{},
				Branch{},
				Branch{},
				Branch{},
			}

			// copy to the original arg
			*a1 = branches
		}).
		Return(nil, nil)

	bh.Get(rr, r)
	require.Equal(t, rr.Code, http.StatusOK)

	branches := []Branch{}

	err := json.NewDecoder(rr.Body).Decode(&branches)
	require.NoError(t, err)

	require.Len(t, branches, 7)
}

func Test_ReturnsErrorWhenUnableToRetrieveBranches(t *testing.T) {
	bh, db := setup(t)

	db.On("Select", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("Boom"))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	bh.Get(rr, r)
	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Contains(t, rr.Body.String(), "Boom")
}

func Test_InsertBranchReturnsOK(t *testing.T) {
	bh, db := setup(t)

	branch := Branch{
		Name:   "unit test name",
		Street: "test street",
		City:   "test city",
		Zip:    "12345",
	}

	// check the db has been called with the correct params
	// this check is also testing the deserialization of the code
	db.On("NamedExec", mock.Anything, branch).Return(nil, nil)

	d, _ := json.Marshal(&branch)

	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(d))
	rr := httptest.NewRecorder()

	bh.Insert(rr, r)
	require.Equal(t, rr.Code, http.StatusOK)
}
