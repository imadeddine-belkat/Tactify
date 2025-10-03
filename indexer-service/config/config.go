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
	BatchSize     int
	FlushInterval time.Duration
}

type PostgresConfig struct {
	Host           string
	Port           int
	Database       string
	User           string
	Password       string
	MaxConnections int
	MinConnections int
	MaxIdleTime    time.Duration
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
