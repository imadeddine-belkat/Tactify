package config

import (
	"encoding/json"
	"fmt"
	"log"

	kafkaConfig "github.com/imadbelkat1/kafka/config"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type FplConfig struct {
	FplApi      FplApi
	KafkaConfig kafkaConfig.KafkaConfig

	DeleteWorkerCount  int `envconfig:"WORKER_DELETE_POOL_SIZE" default:"10"`
	PublishWorkerCount int `envconfig:"WORKER_PUBLISH_POOL_SIZE" default:"100"`
}

type FplApi struct {
	BaseUrl               string `envconfig:"FPLAPI_BASE_URL" required:"true"`
	Bootstrap             string `envconfig:"FPLAPI_BOOTSTRAP" required:"true"`
	Fixtures              string `envconfig:"FPLAPI_FIXTURES" required:"true"`
	PlayerSummary         string `envconfig:"FPLAPI_PLAYER_SUMMARY" required:"true"`
	Entry                 string `envconfig:"FPLAPI_ENTRY" required:"true"`
	EntryEvent            string `envconfig:"FPLAPI_ENTRY_EVENT" required:"true"`
	EntryHistory          string `envconfig:"FPLAPI_ENTRY_HISTORY" required:"true"`
	EntryTransfers        string `envconfig:"FPLAPI_ENTRY_TRANSFERS" required:"true"`
	EntryPicks            string `envconfig:"FPLAPI_ENTRY_PICKS" required:"true"`
	LiveEvent             string `envconfig:"FPLAPI_LIVE_EVENT" required:"true"`
	LeagueClassicStanding string `envconfig:"FPLAPI_LEAGUE_CLASSIC_STANDING" required:"true"`
	LeagueH2hStanding     string `envconfig:"FPLAPI_LEAGUE_H2H_STANDING" required:"true"`
}

type ProcessedModel struct {
	ID   int
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

func (c *FplConfig) DeleteKey(T any, key []string) (map[string]interface{}, error) {
	Bytes, err := json.Marshal(T)
	if err != nil {
		fmt.Println("Error marshaling struct:", err)
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(Bytes, &data)
	if err != nil {
		fmt.Println("Error unmarshaling to map:", err)
		return nil, err
	}

	for _, k := range key {
		delete(data, k)
	}

	return data, nil
}

func (c *FplConfig) ProcessDelete(model interface{}, toBeDeleted []string) ([]byte, error) {
	newElement, err := c.DeleteKey(model, toBeDeleted)
	if err != nil {
		fmt.Errorf("failed to delete keys from newElement: %v", err)
		return nil, err
	}

	elementJSON, err := json.Marshal(newElement)
	if err != nil {
		fmt.Errorf("failed to marshal elementJSON: %v", err)
		return nil, err
	}

	return elementJSON, nil
}
