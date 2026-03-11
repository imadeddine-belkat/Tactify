package config

import (
	"log"
	"strconv"

	kafkaConfig "github.com/imadeddine-belkat/tactify-kafka/config"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type FplConfig struct {
	FplApi      FplApi
	KafkaConfig kafkaConfig.KafkaConfig

	DeleteWorkerCount  int32 `envconfig:"WORKER_DELETE_POOL_SIZE" default:"10"`
	PublishWorkerCount int32 `envconfig:"WORKER_PUBLISH_POOL_SIZE" default:"100"`
}

type FplApi struct {
	BaseUrl               string `envconfig:"FPLAPI_BASE_URL" required:"true"`
	Bootstrap             string `envconfig:"FPLAPI_BOOTSTRAP" required:"true"`
	Fixtures              string `envconfig:"FPLAPI_FIXTURES" required:"true"`
	PlayerSummary         string `envconfig:"FPLAPI_PLAYER_SUMMARY" required:"true"`
	Entry                 string `envconfig:"FPLAPI_ENTRY" required:"true"`
	EntryHistory          string `envconfig:"FPLAPI_ENTRY_HISTORY" required:"true"`
	EntryTransfers        string `envconfig:"FPLAPI_ENTRY_TRANSFERS" required:"true"`
	EntryPicks            string `envconfig:"FPLAPI_ENTRY_PICKS" required:"true"`
	LiveEvent             string `envconfig:"FPLAPI_LIVE_EVENT" required:"true"`
	LeagueClassicStanding string `envconfig:"FPLAPI_LEAGUE_CLASSIC_STANDING" required:"true"`
	LeagueH2hStanding     string `envconfig:"FPLAPI_LEAGUE_H2H_STANDING" required:"true"`

	CurrentSeasonID int32 `envconfig:"FPL_CURRENT_SEASON_ID" required:"true"`
}

type ProcessedModel struct {
	ID   int32
	Data []byte
}

func LoadConfig() *FplConfig {
	// Load .env file (tries multiple paths)
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../../../.env")

	config := &FplConfig{}

	// Parse FplApi config with validation
	if err := envconfig.Process("", config); err != nil {
		log.Fatalf("fpl-service: Unable to load FPL API config: %s", err)
	}

	config.KafkaConfig = *kafkaConfig.LoadConfig()

	return config
}

func (c *FplConfig) MapSeasonNameToID(season string) int32 {
	if len(season) < 4 {
		return 0
	}

	year, err := strconv.Atoi(season[:4])
	if err != nil {
		return 0
	}

	return int32(year)
}
