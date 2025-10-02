package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
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
	TopicsName             TopicsName
	TopicsRetention        TopicsRetention
	ConsumersGroupID       ConsumersGroupID
}

type TopicsName struct {
	// FPL Core Data TopicsRetention
	FplPlayersStats          string `envconfig:"TOPICSNAME_FPL_PLAYERS_STATS"`
	FplPlayerMatchStats      string `envconfig:"TOPICSNAME_FPL_PLAYER_MATCH_HISTORY_STATS"`
	FplPlayerHistoryStats    string `envconfig:"TOPICSNAME_FPL_PLAYER_PAST_HISTORY_STATS"`
	FplTeams                 string `envconfig:"TOPICSNAME_FPL_TEAMS"`
	FplFixtures              string `envconfig:"TOPICSNAME_FPL_FIXTURES"`
	FplFixtureDetails        string `envconfig:"TOPICSNAME_FPL_FIXTURE_DETAILS"`
	FplLiveEvent             string `envconfig:"TOPICSNAME_FPL_LIVE_EVENT"`
	FplEntry                 string `envconfig:"TOPICSNAME_FPL_ENTRY"`
	FplEntryEvent            string `envconfig:"TOPICSNAME_FPL_ENTRY_EVENT"`
	FplEntryHistory          string `envconfig:"TOPICSNAME_FPL_ENTRY_HISTORY"`
	FplEntryTransfers        string `envconfig:"TOPICSNAME_FPL_ENTRY_TRANSFERS"`
	FplEntryPicks            string `envconfig:"TOPICSNAME_FPL_ENTRY_PICKS"`
	FplLeagueClassicStanding string `envconfig:"TOPICSNAME_FPL_LEAGUE_CLASSIC_STANDING"`
	FplLeagueH2hStanding     string `envconfig:"TOPICSNAME_FPL_LEAGUE_H2H_STANDING"`

	// Sofascore Data Topics
	SofascoreLeagueRoundMatches string `envconfig:"TOPICSNAME_SOFASCORE_LEAGUE_ROUND_MATCHES"`
	SofascoreMatchLineups       string `envconfig:"TOPICSNAME_SOFASCORE_MATCH_LINEUPS"`
	SofascoreMatchH2hHistory    string `envconfig:"TOPICSNAME_SOFASCORE_MATCH_H2H_HISTORY"`
	SofascoreTeamOverallStats   string `envconfig:"TOPICSNAME_SOFASCORE_TEAM_OVERALL_STATS"`
	SofascoreTeamMatchStats     string `envconfig:"TOPICSNAME_SOFASCORE_TEAM_MATCH_STATS"`
	SofascorePlayerTeamStats    string `envconfig:"TOPICSNAME_SOFASCORE_PLAYER_TEAM_STATS"`
	SofascorePlayerSeasonsStats string `envconfig:"TOPICSNAME_SOFASCORE_PLAYER_SEASONS_STATS"`
	SofascorePlayerAttributes   string `envconfig:"TOPICSNAME_SOFASCORE_PLAYER_ATTRIBUTES"`
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
	FplEntriesEvent           string `envconfig:"CONSUMERSGROUPID_FPL_ENTRY_EVENT"`
	FplEntriesHistory         string `envconfig:"CONSUMERSGROUPID_FPL_ENTRY_HISTORY"`
	FplEntriesTransfers       string `envconfig:"CONSUMERSGROUPID_FPL_ENTRY_TRANSFERS"`
	FplEntriesPicks           string `envconfig:"CONSUMERSGROUPID_FPL_ENTRY_PICKS"`
	FplLeaguesClassicStanding string `envconfig:"CONSUMERSGROUPID_FPL_LEAGUES_CLASSIC_STANDING"`
	FplLeaguesH2hStanding     string `envconfig:"CONSUMERSGROUPID_FPL_LEAGUES_H2H_STANDING"`

	SofascoreLeagueRoundMatches string `envconfig:"CONSUMERSGROUPID_SOFASCORE_LEAGUE_ROUND_MATCHES"`
	SofascoreMatchLineups       string `envconfig:"CONSUMERSGROUPID_SOFASCORE_MATCH_LINEUPS"`
	SofascoreMatchH2hHistory    string `envconfig:"CONSUMERSGROUPID_SOFASCORE_MATCH_H2H_HISTORY"`
	SofascoreTeamOverallStats   string `envconfig:"CONSUMERSGROUPID_SOFASCORE_TEAM_OVERALL_STATS"`
	SofascoreTeamMatchStats     string `envconfig:"CONSUMERSGROUPID_SOFASCORE_TEAM_MATCH_STATS"`
	SofascorePlayerTeamStats    string `envconfig:"CONSUMERSGROUPID_SOFASCORE_PLAYER_TEAM_STATS"`
	SofascorePlayerSeasonsStats string `envconfig:"CONSUMERSGROUPID_SOFASCORE_PLAYER_SEASONS_STATS"`
	SofascorePlayerAttributes   string `envconfig:"CONSUMERSGROUPID_SOFASCORE_PLAYER_ATTRIBUTES"`

	Test string `envconfig:"CONSUMERSGROUPID_FPL_TEST"`
}

func LoadConfig() *KafkaConfig {
	_, filename, _, _ := runtime.Caller(0)
	ConfigDir := filepath.Dir(filename)
	RootDir := filepath.Dir(ConfigDir)

	// Load .env file
	_ = godotenv.Load(filepath.Join(RootDir, ".env"))
	_ = godotenv.Load(filepath.Join(RootDir, "..", ".env"))
	_ = godotenv.Load(filepath.Join(RootDir, "..", "..", ".env"))

	config := &KafkaConfig{}

	// Parse main Kafka config
	if err := envconfig.Process("", config); err != nil {
		log.Fatalf("kafka: Unable to load Kafka config: %s", err)
	}

	return config
}
