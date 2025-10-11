package config

import (
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
	Season0910      int `envconfig:"FPL_0910_SEASON_ID" required:"true"`
	Season0809      int `envconfig:"FPL_0809_SEASON_ID" required:"true"`
	Season0708      int `envconfig:"FPL_0708_SEASON_ID" required:"true"`
	Season0607      int `envconfig:"FPL_0607_SEASON_ID" required:"true"`
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

func (c *FplConfig) MapSeasonNameToID(seasons string) int {
	switch seasons {
	case "2025/26":
		return c.FplApi.Season2526
	case "2024/25":
		return c.FplApi.Season2425
	case "2023/24":
		return c.FplApi.Season2324
	case "2022/23":
		return c.FplApi.Season2223
	case "2021/22":
		return c.FplApi.Season2122
	case "2020/21":
		return c.FplApi.Season2021
	case "2019/20":
		return c.FplApi.Season1920
	case "2018/19":
		return c.FplApi.Season1819
	case "2017/18":
		return c.FplApi.Season1718
	case "2016/17":
		return c.FplApi.Season1617
	case "2015/16":
		return c.FplApi.Season1516
	case "2014/15":
		return c.FplApi.Season1415
	case "2013/14":
		return c.FplApi.Season1314
	case "2012/13":
		return c.FplApi.Season1213
	case "2011/12":
		return c.FplApi.Season1112
	case "2010/11":
		return c.FplApi.Season1011
	case "2009/10":
		return c.FplApi.Season0910
	case "2008/09":
		return c.FplApi.Season0809
	case "2007/08":
		return c.FplApi.Season0708
	case "2006/07":
		return c.FplApi.Season0607

	default:
		return 0 // Unknown season
	}
}
