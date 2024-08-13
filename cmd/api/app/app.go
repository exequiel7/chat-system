package app

import (
	"chat-system/cmd/api/router"

	"github.com/gocql/gocql"
)

func Start(dbSession *gocql.Session) {
	router := router.NewRouter(":8080")
	router.Setup(dbSession)
}
