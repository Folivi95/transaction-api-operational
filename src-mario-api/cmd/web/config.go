package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/http"
)

type AppConfig struct {
	PostgresConfig PostgresConfig
	ServerConfig   http.ServerConfig
	LogLevel       LogLevel
}

type LogLevel struct {
	LogLevel string `split_words:"true"`
}

type PostgresConfig struct {
	PostgresConnection string `split_words:"true"`
	MockedData         bool   `split_words:"true"`
}

type ServerConfig struct {
	ServerPort   string        `split_words:"true" default:"8080"`
	ReadTimeOut  time.Duration `split_words:"true" default:"2s"`
	WriteTimeOut time.Duration `split_words:"true" default:"2s"`
}

// LoadAppConfig loads the app config from environment variables.
func LoadAppConfig() (AppConfig, error) {
	postgresConfig, err := loadPostgresConfig()
	serverConfig, err := loadServerConfig()
	if err != nil {
		return AppConfig{}, err
	}
	logLevel := loadLogLevel()
	return AppConfig{
		PostgresConfig: postgresConfig,
		ServerConfig:   serverConfig,
		LogLevel:       logLevel,
	}, nil
}

func loadLogLevel() LogLevel {
	var config LogLevel
	err := envconfig.Process("", &config)
	if err != nil {
		return LogLevel{LogLevel: "warn"}
	}
	return config
}

// LoadPostgresConfig loads the postgres config from environment variables.
func loadPostgresConfig() (PostgresConfig, error) {
	var config PostgresConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return PostgresConfig{}, err
	}

	return config, nil
}

// LoadServerConfig loads the server config from environment variables.
func loadServerConfig() (http.ServerConfig, error) {
	var config ServerConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return newDefaultServerConfig(), err
	}

	return http.ServerConfig{
		Port:             config.ServerPort,
		HTTPReadTimeout:  config.ReadTimeOut,
		HTTPWriteTimeout: config.WriteTimeOut,
	}, nil
}

func newDefaultServerConfig() http.ServerConfig {
	return http.ServerConfig{
		Port:             "8080",
		HTTPReadTimeout:  2 * time.Second,
		HTTPWriteTimeout: 2 * time.Second,
	}
}
