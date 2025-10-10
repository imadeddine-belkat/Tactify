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
	EntryHistory          string `envconfig:"FPLAPI_ENTRY_HISTORY" required:"true"`
	EntryTransfers        string `envconfig:"FPLAPI_ENTRY_TRANSFERS" required:"true"`
	EntryPicks            string `envconfig:"FPLAPI_ENTRY_PICKS" required:"true"`
	LiveEvent             string `envconfig:"FPLAPI_LIVE_EVENT" required:"true"`
	LeagueClassicStanding string `envconfig:"FPLAPI_LEAGUE_CLASSIC_STANDING" required:"true"`
	LeagueH2hStanding     string `envconfig:"FPLAPI_LEAGUE_H2H_STANDING" required:"true"`

	CurrentSeasonID int `envconfig:"FPL_CURRENT_SEASON_ID" required:"true"`
	Season2526      int `envconfig:"FPL_2526_SEASON_ID" required:"true"`
	Season2425      int `envconfig:"FPL_2425_SEASON_ID" required:"true"`
	Season2324      int `envconfig:"FPL_2324_SEASON_ID" required:"true"`
	Season2223      int `envconfig:"FPL_2223_SEASON_ID" required:"true"`
	Season2122      int `envconfig:"FPL_2122_SEASON_ID" required:"true"`
	Season2021      int `envconfig:"FPL_2021_SEASON_ID" required:"true"`
	Season1920      int `envconfig:"FPL_1920_SEASON_ID" required:"true"`
	Season1819      int `envconfig:"FPL_1819_SEASON_ID" required:"true"`
	Season1718      int `envconfig:"FPL_1718_SEASON_ID" required:"true"`
	Season1617      int `envconfig:"FPL_1617_SEASON_ID" required:"true"`
	Season1516      int `envconfig:"FPL_1516_SEASON_ID" required:"true"`
	Season1415      int `envconfig:"FPL_1415_SEASON_ID" required:"true"`
	Season1314      int `envconfig:"FPL_1314_SEASON_ID" required:"true"`
	Season1213      int `envconfig:"FPL_1213_SEASON_ID" required:"true"`
	Season1112      int `envconfig:"FPL_1112_SEASON_ID" required:"true"`
	Season1011      int `envconfig:"FPL_1011_SEASON_ID" required:"true"`
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
