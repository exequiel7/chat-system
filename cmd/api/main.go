package main

import (
	"chat-system/cmd/api/app"
	"chat-system/cmd/api/config"
	"chat-system/cmd/api/databases"
	"context"
	"fmt"

	logger "github.com/sirupsen/logrus"
)

func main() {
	logger.Info("starting application...")
	ctx := context.Background()
	env := config.DefineEnvironment(ctx)
	_, err := config.InitConfig(ctx,
		config.Environment(env),
		config.WithConfigFile("./cmd/api/env", fmt.Sprintf("%s.config", env), "env"),
	)

	if err != nil {
		logger.Panic(ctx, "fatal reading the config")
	}

	dbProvider := databases.NewCassandraDB()
	defer dbProvider.GetDbClient().Close()

	dbProvider.CreateSchema()
	app.Start(dbProvider.GetDbClient())
}
