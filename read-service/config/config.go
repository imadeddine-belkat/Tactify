package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type ReaderConfig struct {
	Postgres    PostgresConfig
	WorkerCount int
}
type PostgresConfig struct {
	Host              string `envconfig:"DB_HOST" default:"localhost"`
	Port              int    `envconfig:"DB_PORT" default:"5432"`
	FplDatabase       string `envconfig:"DB_FPL_NAME" default:"fpl"`
	SofascoreDatabase string `envconfig:"DB_SOFASCORE_NAME" default:"sofascore"`
	User              string `envconfig:"DB_USER" default:"tactify"`
	Password          string `envconfig:"DB_PASSWORD" default:"admin"`
}

func LoadConfig() *ReaderConfig {
	// Load .env file (tries multiple paths)
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../../../.env")

	config := &ReaderConfig{}

	// Parse FplApi config with validation
	if err := envconfig.Process("", config); err != nil {
		log.Fatalf("read-service: Unable to load config: %s", err)
	}

	return config
}
