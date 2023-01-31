package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/http"
)

// dim tables' names.
type AppConfig struct {
	KafkaConfig          KafkaConfig
	PostgresConfig       PostgresConfig
	SchemaRegistryConfig SchemaRegistryConfig
	LogLevel             string
}

type SchemaRegistryConfig struct {
	SchemaRegistryEndpoint           string `split_words:"true"`
	SchemaRegistryRefreshTimeSeconds int    `split_words:"true"`
	W4IngressSchemaKey               string `split_words:"true"`
	SolarIngressSchemaKey            string `split_words:"true"`
	W4EgressSchemaKey                string `split_words:"true"`
	SolarEgressSchemaKey             string `split_words:"true"`
	AcntContractSchemaKey            string `split_words:"true"`
	BinTableSchemaKey                string `split_words:"true"`
	ClientAddressSchemaKey           string `split_words:"true"`
	TransCondSchemaKey               string `split_words:"true"`
	TransTypeSchemaKey               string `split_words:"true"`
}

type KafkaConfig struct {
	KafkaEndpoint                  []string `split_words:"true"`
	KafkaUsername                  string   `split_words:"true"`
	KafkaPassword                  string   `split_words:"true"`
	KafkaEgressTopic               string   `split_words:"true"`
	KafkaW4IngressTopic            string   `split_words:"true"`
	KafkaSolarIngressTopic         string   `split_words:"true"`
	KafkaAcntContractIngressTopic  string   `split_words:"true"`
	KafkaClientAddressIngressTopic string   `split_words:"true"`
	KafkaTransTypeIngressTopic     string   `split_words:"true"`
	KafkaTransCondIngressTopic     string   `split_words:"true"`
	KafkaBinTableIngressTopic      string   `split_words:"true"`
	MockedData                     bool     `split_words:"true"`
	ListenDimTables                bool     `split_words:"true"`
	ListenTransactions             bool     `split_words:"true"`
}

type PostgresConfig struct {
	PostgresConnection string `split_words:"true"`
}

func loadAppConfig() (AppConfig, error) {
	kafkaConfig, err := loadKafkaConfig()
	if err != nil {
		return AppConfig{}, err
	}
	postgresConfig, err := loadPostgresConfig()
	if err != nil {
		return AppConfig{}, err
	}
	schemaRegistryConfig, err := loadSchemaRegistryConfig()
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		KafkaConfig:          kafkaConfig,
		PostgresConfig:       postgresConfig,
		SchemaRegistryConfig: schemaRegistryConfig,
	}, nil
}

// LoadKafkaConfig loads the app config from environment variables.
func loadKafkaConfig() (KafkaConfig, error) {
	var config KafkaConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return KafkaConfig{}, err
	}
	// brokers := make([]string, len(config.KafkaEndpoint))
	// for key, broker := range config.KafkaEndpoint {
	// 	u, err := url.Parse(broker)
	// 	if err != nil {
	// 		config.KafkaEndpoint[key] = u.Host
	// 	}
	// }

	return config, nil
}

func loadPostgresConfig() (PostgresConfig, error) {
	var config PostgresConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return PostgresConfig{}, err
	}

	return config, nil
}

func loadSchemaRegistryConfig() (SchemaRegistryConfig, error) {
	var config SchemaRegistryConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return SchemaRegistryConfig{}, err
	}

	return config, nil
}

func newDefaultConfig() http.ServerConfig {
	return http.ServerConfig{
		Port:             "8080",
		HTTPReadTimeout:  2 * time.Second,
		HTTPWriteTimeout: 2 * time.Second,
	}
}
