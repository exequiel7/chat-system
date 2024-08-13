package databases

import "github.com/gocql/gocql"

type DBProvider interface {
	GetDbClient() *gocql.Session
}
