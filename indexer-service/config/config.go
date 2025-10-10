package config

import (
	"log"
	"time"

	kafkaConfig "github.com/imadbelkat1/kafka/config"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type IndexerConfig struct {
	Kafka         kafkaConfig.KafkaConfig
	Postgres      PostgresConfig
	WorkerCount   int
	BatchSize     int           `envconfig:"DB_BATCH_SIZE" default:"100"`
	FlushInterval time.Duration `envconfig:"DB_FLUSH_INTERVAL" default:"30s"`
}

type PostgresConfig struct {
	Host              string `envconfig:"DB_HOST" default:"localhost"`
	Port              int    `envconfig:"DB_PORT" default:"5432"`
	FplDatabase       string `envconfig:"DB_NAME" default:"fpl"`
	SofascoreDatabase string `envconfig:"DB_SOFA_DB_NAME" default:"sofascore"`
	User              string `envconfig:"DB_USER" default:"tactify"`
	Password          string `envconfig:"DB_PASSWORD" default:"admin"`
}

func LoadConfig() *IndexerConfig {
	// Load .env file (tries multiple paths)
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../../../.env")

	config := &IndexerConfig{}

	// Parse FplApi config with validation
	if err := envconfig.Process("", config); err != nil {
		log.Fatalf("fpl-service: Unable to load FPL API config: %s", err)
	}

	config.Kafka = *kafkaConfig.LoadConfig()

	return config
}
