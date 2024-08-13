package models

import (
	"time"

	"github.com/gocql/gocql"
)

type User struct {
	Id        gocql.UUID `json:"id" cql:"id uuid PRIMARY KEY"`
	Name      string     `json:"name" cql:"name text"`
	Surname   string     `json:"surname" cql:"surname text"`
	Username  string     `json:"username" cql:"username text"`
	Password  string     `json:"password" cql:"password text"`
	Email     string     `json:"email" cql:"email text"`
	CreatedAt time.Time  `json:"-" cql:"created_at timestamp"`
	UpdatedAt time.Time  `json:"-" cql:"updated_at timestamp"`
}
