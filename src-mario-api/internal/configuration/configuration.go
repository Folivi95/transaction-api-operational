package configuration

import (
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/http"
)

type AppConfig struct {
	PostgresConfig       PostgresConfig
	KafkaConfig          KafkaConfig
	SchemaRegistryConfig SchemaRegistryConfig
	LogLevel             LogConfig
}

type LogConfig struct {
	LogLevel string `split_words:"true"`
}

type SchemaRegistryConfig struct {
	SchemaRegistryEndpoint           string `split_words:"true"`
	SchemaRegistryRefreshTimeSeconds int    `split_words:"true"`
}

type PostgresConfig struct {
	PostgresConnection string `split_words:"true"`
	MockedData         bool   `split_words:"true"`
}

type KafkaConfig struct {
	KafkaTransactionTopic  string `split_words:"true"`
	KafkaEndpoint          string `split_words:"true"`
	KafkaUsername          string `split_words:"true"`
	KafkaPassword          string `split_words:"true"`
	KafkaGroupId           string `split_words:"true"`
	KafkaTransactionSchema string `split_words:"true"`
}

func LoadAppConfig() (AppConfig, error) {
	postgresConfig, err := loadPostgresConfig()
	if err != nil {
		return AppConfig{}, err
	}

	kafkaConfig, err := loadKafkaConfig()
	if err != nil {
		return AppConfig{}, err
	}

	logLevelConfig, err := loadLogLevelConfig()
	if err != nil {
		return AppConfig{}, err
	}

	schemaRegistryConfig, err := loadSchemaRegistryConfig()
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		PostgresConfig:       postgresConfig,
		KafkaConfig:          kafkaConfig,
		SchemaRegistryConfig: schemaRegistryConfig,
		LogLevel:             logLevelConfig,
	}, nil
}

// loadPostgresConfig loads postgres configuration from environment variables.
func loadPostgresConfig() (PostgresConfig, error) {
	var config PostgresConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return PostgresConfig{}, err
	}

	return config, nil
}

// loadKafkaConfig loads kafka configuration from the environment variables.
func loadKafkaConfig() (KafkaConfig, error) {
	var config KafkaConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return KafkaConfig{}, err
	}

	return config, nil
}

// loadSchemaRegistryConfig loads kafka configuration from the environment variables.
func loadSchemaRegistryConfig() (SchemaRegistryConfig, error) {
	var config SchemaRegistryConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return SchemaRegistryConfig{}, err
	}

	return config, nil
}

// loadLogLevelConfig loads log level configuration from the environment variables.
func loadLogLevelConfig() (LogConfig, error) {
	var config LogConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return LogConfig{}, err
	}

	return config, nil
}

func NewDefaultConfig() http.ServerConfig {
	return http.ServerConfig{
		Port:             "8080",
		HTTPReadTimeout:  2 * time.Second,
		HTTPWriteTimeout: 2 * time.Second,
	}
}
