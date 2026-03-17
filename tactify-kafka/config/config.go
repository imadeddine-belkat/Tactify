package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type KafkaConfig struct {
	KafkaBroker            string `envconfig:"KAFKA_BROKER"`
	KafkaAcks              string `envconfig:"KAFKA_ACKS"`
	KafkaRetries           int    `envconfig:"KAFKA_RETRIES"`
	KafkaRetryBackoffMs    int    `envconfig:"KAFKA_RETRY_BACKOFF_MS"`
	KafkaDeliveryTimeoutMs int    `envconfig:"KAFKA_DELIVERY_TIMEOUT_MS"`
	KafkaBatchSize         int    `envconfig:"KAFKA_BATCH_SIZE"`
	KafkaLingerMs          int    `envconfig:"KAFKA_LINGER_MS"`
	KafkaCompressionType   string `envconfig:"KAFKA_COMPRESSION_TYPE"`
	KafkaBufferMemory      int    `envconfig:"KAFKA_BUFFER_MEMORY"`
	KafkaPartitions        int    `envconfig:"KAFKA_PARTITIONS"`
	KafkaReplication       int    `envconfig:"KAFKA_REPLICATION"`
	TopicsName             Topics
	TopicsRetention        TopicsRetention
	ConsumersGroupID       ConsumersGroupID
}

type TopicsConfig struct {
	Kafka struct {
		Topics Topics `yaml:"topics"`
	} `yaml:"tactify-kafka"`
}
type Topic struct {
	Name       string `yaml:"name"`
	Partitions int    `yaml:"partitions"`
}

type Topics struct {
	// FPL Core Data TopicsRetention
	FplPlayersBootstrap      Topic `yaml:"fpl_players_bootstrap"`
	FplPlayersStats          Topic `yaml:"fpl_players_stats"`
	FplPlayerMatchStats      Topic `yaml:"fpl_player_match_history_stats"`
	FplPlayerHistoryStats    Topic `yaml:"fpl_player_past_history_stats"`
	FplTeams                 Topic `yaml:"fpl_teams"`
	FplFixtures              Topic `yaml:"fpl_fixtures"`
	FplFixtureDetails        Topic `yaml:"fpl_fixture_details"`
	FplLiveEvent             Topic `yaml:"fpl_live_event"`
	FplEntry                 Topic `yaml:"fpl_entry"`
	FplEntryHistory          Topic `yaml:"fpl_entry_history"`
	FplEntryTransfers        Topic `yaml:"fpl_entry_transfers"`
	FplEntryPicks            Topic `yaml:"fpl_entry_picks"`
	FplLeagueClassicStanding Topic `yaml:"fpl_league_classic_standing"`
	FplLeagueH2hStanding     Topic `yaml:"fpl_league_h2h_standing"`

	// Sofascore
	SofascoreLeagueIDs          Topic `yaml:"sofascore_league_ids"`
	SofascoreLeagueSeasons      Topic `yaml:"sofascore_league_seasons"`
	SofascoreLeagueStandings    Topic `yaml:"sofascore_league_standings"`
	SofascoreLeagueRoundMatches Topic `yaml:"sofascore_league_round_matches"`
	SofascoreMatchLineups       Topic `yaml:"sofascore_match_lineups"`
	SofascoreMatchH2hHistory    Topic `yaml:"sofascore_match_h2h_history"`
	SofascoreTopTeamsStats      Topic `yaml:"sofascore_top_teams_stats"`
	SofascoreTeamOverallStats   Topic `yaml:"sofascore_team_overall_stats"`
	SofascoreTeamMatchStats     Topic `yaml:"sofascore_team_match_stats"`
	SofascorePlayerInfo         Topic `yaml:"sofascore_player_info"`
	SofascorePlayerTeamStats    Topic `yaml:"sofascore_player_team_stats"`
	SofascorePlayerSeasonsStats Topic `yaml:"sofascore_player_seasons_stats"`
	SofascorePlayerAttributes   Topic `yaml:"sofascore_player_attributes"`
	SofascorePlayerMatchStats   Topic `yaml:"sofascore_player_match_stats"`
}
type TopicsRetention struct {
	FplPlayers               string `envconfig:"TOPICSRETENTION_FPL_PLAYERS"`
	FplTeams                 string `envconfig:"TOPICSRETENTION_FPL_TEAMS"`
	FplFixtures              string `envconfig:"TOPICSRETENTION_FPL_FIXTURES"`
	FplPlayerMatchStats      string `envconfig:"TOPICSRETENTION_FPL_PLAYER_MATCH_STATS"`
	FplEntry                 string `envconfig:"TOPICSRETENTION_FPL_ENTRY"`
	FplEntryEvent            string `envconfig:"TOPICSRETENTION_FPL_ENTRY_EVENT"`
	FplEntryHistory          string `envconfig:"TOPICSRETENTION_FPL_ENTRY_HISTORY"`
	FplEntryTransfers        string `envconfig:"TOPICSRETENTION_FPL_ENTRY_TRANSFERS"`
	FplEntryPicks            string `envconfig:"TOPICSRETENTION_FPL_ENTRY_PICKS"`
	FplLeagueClassicStanding string `envconfig:"TOPICSRETENTION_FPL_LEAGUE_CLASSIC_STANDING"`
	FplLeagueH2hStanding     string `envconfig:"TOPICSRETENTION_FPL_LEAGUE_H2H_STANDING"`
}
type ConsumersGroupID struct {
	// Consumer Group IDs
	FplTeams                  string `envconfig:"CONSUMERSGROUPID_FPL_TEAMS"`
	FplFixtures               string `envconfig:"CONSUMERSGROUPID_FPL_FIXTURES"`
	FplPlayers                string `envconfig:"CONSUMERSGROUPID_FPL_PLAYERS"`
	FplPlayersStats           string `envconfig:"CONSUMERSGROUPID_FPL_PLAYERS_STATS"`
	FplLive                   string `envconfig:"CONSUMERSGROUPID_FPL_LIVE_EVENT"`
	FplEntries                string `envconfig:"CONSUMERSGROUPID_FPL_ENTRY"`
	FplEntriesHistory         string `envconfig:"CONSUMERSGROUPID_FPL_ENTRY_HISTORY"`
	FplEntriesTransfers       string `envconfig:"CONSUMERSGROUPID_FPL_ENTRY_TRANSFERS"`
	FplEntriesPicks           string `envconfig:"CONSUMERSGROUPID_FPL_ENTRY_PICKS"`
	FplLeaguesClassicStanding string `envconfig:"CONSUMERSGROUPID_FPL_LEAGUES_CLASSIC_STANDING"`
	FplLeaguesH2hStanding     string `envconfig:"CONSUMERSGROUPID_FPL_LEAGUES_H2H_STANDING"`

	SofascoreLeagueStanding     string `envconfig:"CONSUMERSGROUPID_SOFASCORE_LEAGUE_STANDINGS"`
	SofascoreLeagueRoundMatches string `envconfig:"CONSUMERSGROUPID_SOFASCORE_LEAGUE_ROUND_MATCHES"`
	SofascoreMatchLineups       string `envconfig:"CONSUMERSGROUPID_SOFASCORE_MATCH_LINEUPS"`
	SofascoreMatchH2hHistory    string `envconfig:"CONSUMERSGROUPID_SOFASCORE_MATCH_H2H_HISTORY"`
	SofascoreTeamOverallStats   string `envconfig:"CONSUMERSGROUPID_SOFASCORE_TEAM_OVERALL_STATS"`
	SofascoreTeamMatchStats     string `envconfig:"CONSUMERSGROUPID_SOFASCORE_TEAM_MATCH_STATS"`
	SofascorePlayerTeamStats    string `envconfig:"CONSUMERSGROUPID_SOFASCORE_PLAYER_TEAM_STATS"`
	SofascorePlayerInfo         string `envconfig:"CONSUMERSGROUPID_SOFASCORE_PLAYER_INFO"`
	SofascorePlayerSeasonsStats string `envconfig:"CONSUMERSGROUPID_SOFASCORE_PLAYER_SEASONS_STATS"`
	SofascorePlayerAttributes   string `envconfig:"CONSUMERSGROUPID_SOFASCORE_PLAYER_ATTRIBUTES"`

	Test string `envconfig:"CONSUMERSGROUPID_FPL_TEST"`
}

func LoadConfig() *KafkaConfig {
	_, filename, _, _ := runtime.Caller(0)

	// tactify-kafka/config/
	configDir := filepath.Dir(filename)
	// project root (tactify-kafka/)
	rootDir := filepath.Dir(configDir)

	// Load .env
	_ = godotenv.Load(filepath.Join(rootDir, ".env"))
	_ = godotenv.Load(filepath.Join(rootDir, "..", ".env"))

	cfg := &KafkaConfig{}

	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("tactify-kafka: unable to load env config: %v", err)
	}

	// Load topics.yaml
	topicsPath := filepath.Join(rootDir, "topics.yaml")

	data, err := os.ReadFile(topicsPath)
	if err != nil {
		log.Fatalf("tactify-kafka: unable to read topics.yaml at %s: %v", topicsPath, err)
	}

	var topicsCfg TopicsConfig
	if err := yaml.Unmarshal(data, &topicsCfg); err != nil {
		log.Fatalf("tactify-kafka: unable to parse topics.yaml: %v", err)
	}

	cfg.TopicsName = topicsCfg.Kafka.Topics
	return cfg
}
