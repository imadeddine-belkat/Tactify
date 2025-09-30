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
	FplPlayers               string `envconfig:"TOPICSNAME_FPL_PLAYERS"`
	FplTeams                 string `envconfig:"TOPICSNAME_FPL_TEAMS"`
	FplFixtures              string `envconfig:"TOPICSNAME_FPL_FIXTURES"`
	FplFixtureDetails        string `envconfig:"TOPICSNAME_FPL_FIXTURE_DETAILS"`
	FplPlayerMatchStats      string `envconfig:"TOPICSNAME_FPL_PLAYER_MATCH_STATS"`
	FplLiveEvent             string `envconfig:"TOPICSNAME_FPL_LIVE_EVENT"`
	FplEntry                 string `envconfig:"TOPICSNAME_FPL_ENTRY"`
	FplEntryEvent            string `envconfig:"TOPICSNAME_FPL_ENTRY_EVENT"`
	FplEntryHistory          string `envconfig:"TOPICSNAME_FPL_ENTRY_HISTORY"`
	FplEntryTransfers        string `envconfig:"TOPICSNAME_FPL_ENTRY_TRANSFERS"`
	FplEntryPicks            string `envconfig:"TOPICSNAME_FPL_ENTRY_PICKS"`
	FplLeagueClassicStanding string `envconfig:"TOPICSNAME_FPL_LEAGUE_CLASSIC_STANDING"`
	FplLeagueH2hStanding     string `envconfig:"TOPICSNAME_FPL_LEAGUE_H2H_STANDING"`
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
	Teams                  string `envconfig:"CONSUMERSGROUPID_KAFKA_TEAMS"`
	Fixtures               string `envconfig:"CONSUMERSGROUPID_KAFKA_FIXTURES"`
	Players                string `envconfig:"CONSUMERSGROUPID_KAFKA_PLAYERS"`
	PlayersStats           string `envconfig:"CONSUMERSGROUPID_KAFKA_PLAYERS_STATS"`
	Live                   string `envconfig:"CONSUMERSGROUPID_KAFKA_LIVE_EVENT"`
	Entries                string `envconfig:"CONSUMERSGROUPID_KAFKA_ENTRY"`
	EntriesEvent           string `envconfig:"CONSUMERSGROUPID_KAFKA_ENTRY_EVENT"`
	EntriesHistory         string `envconfig:"CONSUMERSGROUPID_KAFKA_ENTRY_HISTORY"`
	EntriesTransfers       string `envconfig:"CONSUMERSGROUPID_KAFKA_ENTRY_TRANSFERS"`
	EntriesPicks           string `envconfig:"CONSUMERSGROUPID_KAFKA_ENTRY_PICKS"`
	LeaguesClassicStanding string `envconfig:"CONSUMERSGROUPID_KAFKA_LEAGUES_CLASSIC_STANDING"`
	LeaguesH2hStanding     string `envconfig:"CONSUMERSGROUPID_KAFKA_LEAGUES_H2H_STANDING"`
	Test                   string `envconfig:"CONSUMERSGROUPID_KAFKA_TEST"`
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
