package mocks

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/mock"
)

type MockSession struct {
	mock.Mock
}

func (m *MockSession) Query(stmt string, values ...interface{}) *gocql.Query {
	args := m.Called(stmt, values)
	return args.Get(0).(*gocql.Query)
}

func (m *MockSession) Close() {
	m.Called()
}

type MockQuery struct {
	mock.Mock
}

func (m *MockQuery) WithContext(ctx context.Context) *gocql.Query {
	args := m.Called(ctx)
	return args.Get(0).(*gocql.Query)
}

func (m *MockQuery) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *MockQuery) Iter() *gocql.Iter {
	args := m.Called()
	return args.Get(0).(*gocql.Iter)
}
