package data

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

// DB implements the methods used from sqlx.DB to remove hard dependency on the package for tests
type DB interface {
	Select(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

// DBMock is used in tests to validate the DB interface is called correctly
type DBMock struct {
	mock.Mock
}

func (m *DBMock) Select(dest interface{}, query string, args ...interface{}) error {
	ra := m.Called(dest, query, args)

	return ra.Error(0)
}

func (m *DBMock) NamedExec(query string, arg interface{}) (sql.Result, error) {
	ra := m.Called(query, arg)

	if sr, ok := ra.Get(0).(sql.Result); ok {
		return sr, ra.Error(1)
	}

	return nil, ra.Error(1)
}
