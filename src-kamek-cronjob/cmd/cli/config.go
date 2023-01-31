package main

import (
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	PostgresConfig PostgresConfig
}

type PostgresConfig struct {
	PostgresURL string `split_words:"true"`
}

func loadAppConfig() (AppConfig, error) {
	postgresConfig, err := loadPostgresConfig()
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		PostgresConfig: postgresConfig,
	}, nil
}

// LoadPostgresConfig loads the app config from environment variables.
func loadPostgresConfig() (PostgresConfig, error) {
	var config PostgresConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return PostgresConfig{}, err
	}

	return config, nil
}
