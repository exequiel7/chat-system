package config

import (
	"context"
	"fmt"
	"os"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Environment represents the available environments
//
//go:generate go run github.com/lindell/string-enumer --text -t Environment -o ./generated_enviroment_enum_utils.go .
type Environment string

var env *Environment

const (
	LOCAL Environment = "local"
	DEV   Environment = "dev"
	TEST  Environment = "test"
)

// Config represents the application configuration.
type Config struct {
	Environment       Environment `mapstructure:"ENVIRONMENT"`
	ServiceName       string      `mapstructure:"SERVICE_NAME"`
	CassandraPort     int         `mapstructure:"CASSANDRA_PORT"`
	CassandraHost     string      `mapstructure:"CASSANDRA_HOST"`
	CassandraKeyspace string      `mapstructure:"CASSANDRA_KEYSPACE"`
	JwtSecretKey      string      `mapstructure:"JWT_SECRET_KEY"`
}

var (
	config *Config
	once   sync.Once
)

func InitConfig(ctx context.Context, env Environment, options ...Option) (*Config, error) {
	var err error
	once.Do(func() {
		config, err = readConfig(ctx, env, options...)
	})
	return config, err
}

func DefaultTestConfig() *Config {
	config = &Config{
		Environment: TEST,
		ServiceName: "chat-system",
	}

	return config
}

func readConfig(ctx context.Context, env Environment, options ...Option) (*Config, error) {
	conf := &Config{
		Environment: env,
	}
	for _, opts := range options {
		opts(conf)
	}
	viper.AutomaticEnv()
	viper.SetDefault("SERVICE_NAME", "chat-system")
	viper.SetDefault("CASSANDRA_PORT", "9042")

	if err := viper.ReadInConfig(); err != nil {
		logger.Error(ctx, "Unable to read config. Assuming env variables or defaults will be used.")
		return conf, err
	}

	if err := viper.Unmarshal(&conf); err != nil {
		logger.Error(ctx, "Unable to decode into struct.")
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	return conf, nil
}

func GetConfig() *Config {
	return config
}

func DefineEnvironment(ctx context.Context) Environment {
	if config != nil {
		return config.Environment
	}

	env := LOCAL
	e, exists := os.LookupEnv("ENVIRONMENT")
	if !exists {
		logger.Info(ctx, fmt.Sprintf("'ENVIRONMENT' variable not defined. Environment set to '%s'", env))
		return env
	}

	environment := Environment(e)
	logger.Info(ctx, fmt.Sprintf("'ENVIRONMENT' variable defined as '%s'", environment))
	if !environment.Valid() {
		logger.Warn(ctx, fmt.Sprintf("'ENVIRONMENT' variable value (%s) is not valid. Environment set to '%s'", e, env))
		return env
	}

	return environment
}
