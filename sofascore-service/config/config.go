package config

import (
	"encoding/json"
	"fmt"
	"log"

	kafka "github.com/imadbelkat1/kafka/config"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type SofascoreConfig struct {
	SofascoreApi SofascoreApi
	KafkaConfig  kafka.KafkaConfig

	PublishWorkerCount int `envconfig:"WORKER_PUBLISH_POOL_SIZE" default:"100"`
	FetchWorkerCount   int `envconfig:"WORKER_FETCH_POOL_SIZE" default:"50"`
}

type SofascoreApi struct {
	BaseURL         string `envconfig:"SOFASCOREAPI_BASE_URL"`
	SeasonsIDs      SeasonsIDs
	LeaguesIDs      LeaguesIDs
	LeagueEndpoints LeagueEndpoints
	MatchEndpoints  MatchEndpoints
	TeamEndpoints   TeamEndpoints
	PlayerEndpoints PlayerEndpoints
}

type SeasonsIDs struct {
	LaLiga2526 int `envconfig:"SOFASCOREAPI_LALIGA_2526_SEASON_ID"`
	Laliga2425 int `envconfig:"SOFASCOREAPI_LALIGA_2425_SEASON_ID"`
	LaLiga2324 int `envconfig:"SOFASCOREAPI_LALIGA_2324_SEASON_ID"`
	Laliga2223 int `envconfig:"SOFASCOREAPI_LALIGA_2223_SEASON_ID"`
	Laliga2122 int `envconfig:"SOFASCOREAPI_LALIGA_2122_SEASON_ID"`
	Laliga2021 int `envconfig:"SOFASCOREAPI_LALIGA_2021_SEASON_ID"`
	Laliga1920 int `envconfig:"SOFASCOREAPI_LALIGA_1920_SEASON_ID"`
	Laliga1819 int `envconfig:"SOFASCOREAPI_LALIGA_1819_SEASON_ID"`
	Laliga1718 int `envconfig:"SOFASCOREAPI_LALIGA_1718_SEASON_ID"`
	Laliga1617 int `envconfig:"SOFASCOREAPI_LALIGA_1617_SEASON_ID"`
	Laliga1516 int `envconfig:"SOFASCOREAPI_LALIGA_1516_SEASON_ID"`
	Laliga1415 int `envconfig:"SOFASCOREAPI_LALIGA_1415_SEASON_ID"`
	Laliga1314 int `envconfig:"SOFASCOREAPI_LALIGA_1314_SEASON_ID"`
	Laliga1213 int `envconfig:"SOFASCOREAPI_LALIGA_1213_SEASON_ID"`
	Laliga1112 int `envconfig:"SOFASCOREAPI_LALIGA_1112_SEASON_ID"`
	Laliga1011 int `envconfig:"SOFASCOREAPI_LALIGA_1011_SEASON_ID"`

	PremierLeague2526 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_2526_SEASON_ID"`
	PremierLeague2425 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_2425_SEASON_ID"`
	PremierLeague2324 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_2324_SEASON_ID"`
	PremierLeague2223 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_2223_SEASON_ID"`
	PremierLeague2122 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_2122_SEASON_ID"`
	PremierLeague2021 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_2021_SEASON_ID"`
	PremierLeague1920 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_1920_SEASON_ID"`
	PremierLeague1819 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_1819_SEASON_ID"`
	PremierLeague1718 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_1718_SEASON_ID"`
	PremierLeague1617 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_1617_SEASON_ID"`
	PremierLeague1516 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_1516_SEASON_ID"`
	PremierLeague1415 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_1415_SEASON_ID"`
	PremierLeague1314 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_1314_SEASON_ID"`
	PremierLeague1213 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_1213_SEASON_ID"`
	PremierLeague1112 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_1112_SEASON_ID"`
	PremierLeague1011 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_1011_SEASON_ID"`
	PremierLeague0910 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_0910_SEASON_ID"`
	PremierLeague0809 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_0809_SEASON_ID"`
	PremierLeague0708 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_0708_SEASON_ID"`
	PremierLeague0607 int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_0607_SEASON_ID"`
}

type LeaguesIDs struct {
	LaLiga        int `envconfig:"SOFASCOREAPI_LALIGA_ID"`
	PremierLeague int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_ID"`
}
type LeagueEndpoints struct {
	LeagueRoundMatches string `envconfig:"SOFASCOREAPI_LEAGUE_ROUND_MATCHES_ENDPOINT"`
	LaLigaId           string `envconfig:"SOFASCOREAPI_LALIGA_ID"`
	PremierLeagueId    string `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_ID"`
}
type MatchEndpoints struct {
	MatchLineups    string `envconfig:"SOFASCOREAPI_MATCH_LINEUPS_ENDPOINT"`
	MatchH2hHistory string `envconfig:"SOFASCOREAPI_MATCH_H2H_HISTORY_ENDPOINT"`
}

type TeamEndpoints struct {
	TeamOverallStats string `envconfig:"SOFASCOREAPI_TEAM_OVERALL_STATS_ENDPOINT"`
	TeamMatchStats   string `envconfig:"SOFASCOREAPI_TEAM_MATCH_STATS_ENDPOINT"`
}
type PlayerEndpoints struct {
	PlayersStats       string `envconfig:"SOFASCOREAPI_TEAM_PLAYERS_STATS_ENDPOINT"`
	PlayerSeasonsStats string `envconfig:"SOFASCOREAPI_PLAYER_SEASONS_STATS_ENDPOINT"`
	PlayerAttributes   string `envconfig:"SOFASCOREAPI_PLAYER_ATTRIBUTES_ENDPOINT"`
}

func LoadConfig() *SofascoreConfig {
	// Load .env file (tries multiple paths)
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../../../.env")

	config := &SofascoreConfig{}

	// Parse FplApi config with validation
	if err := envconfig.Process("", config); err != nil {
		log.Fatalf("sofascore-service: Unable to load FPL API config: %s", err)
	}

	config.KafkaConfig = *kafka.LoadConfig()

	return config
}

func (c *SofascoreConfig) DeleteKey(T any, key []string) (map[string]any, error) {
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

func (c *SofascoreConfig) ProcessDelete(model any, toBeDeleted []string) ([]byte, error) {
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
