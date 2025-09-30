package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	SofascoreApi     SofascoreApi
	TopicsName       TopicsName
	ConsumersGroupId ConsumersGroupId
}

type SofascoreApi struct {
	BaseURL         string `mapstructure:"SOFASCOREAPI_BASE_URL"`
	SeasonsIDs      SeasonsIDs
	LeagueEndpoints LeagueEndpoints
	MatchEndpoints  MatchEndpoints
	TeamEndpoints   TeamEndpoints
	PlayerEndpoints PlayerEndpoints
}

type SeasonsIDs struct {
	LaLiga2526 string `mapstructure:"SOFASCOREAPI_LALIGA_2526_SEASON_ID"`
	Laliga2425 string `mapstructure:"SOFASCOREAPI_LALIGA_2425_SEASON_ID"`
	LaLiga2324 string `mapstructure:"SOFASCOREAPI_LALIGA_2324_SEASON_ID"`
	Laliga2223 string `mapstructure:"SOFASCOREAPI_LALIGA_2223_SEASON_ID"`
	Laliga2122 string `mapstructure:"SOFASCOREAPI_LALIGA_2122_SEASON_ID"`
	Laliga2021 string `mapstructure:"SOFASCOREAPI_LALIGA_2021_SEASON_ID"`
	Laliga1920 string `mapstructure:"SOFASCOREAPI_LALIGA_1920_SEASON_ID"`
	Laliga1819 string `mapstructure:"SOFASCOREAPI_LALIGA_1819_SEASON_ID"`
	Laliga1718 string `mapstructure:"SOFASCOREAPI_LALIGA_1718_SEASON_ID"`
	Laliga1617 string `mapstructure:"SOFASCOREAPI_LALIGA_1617_SEASON_ID"`
	Laliga1516 string `mapstructure:"SOFASCOREAPI_LALIGA_1516_SEASON_ID"`
	Laliga1415 string `mapstructure:"SOFASCOREAPI_LALIGA_1415_SEASON_ID"`
	Laliga1314 string `mapstructure:"SOFASCOREAPI_LALIGA_1314_SEASON_ID"`
	Laliga1213 string `mapstructure:"SOFASCOREAPI_LALIGA_1213_SEASON_ID"`
	Laliga1112 string `mapstructure:"SOFASCOREAPI_LALIGA_1112_SEASON_ID"`
	Laliga1011 string `mapstructure:"SOFASCOREAPI_LALIGA_1011_SEASON_ID"`

	PremierLeague2526 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_2526_SEASON_ID"`
	PremierLeague2425 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_2425_SEASON_ID"`
	PremierLeague2324 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_2324_SEASON_ID"`
	PremierLeague2223 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_2223_SEASON_ID"`
	PremierLeague2122 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_2122_SEASON_ID"`
	PremierLeague2021 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_2021_SEASON_ID"`
	PremierLeague1920 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_1920_SEASON_ID"`
	PremierLeague1819 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_1819_SEASON_ID"`
	PremierLeague1718 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_1718_SEASON_ID"`
	PremierLeague1617 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_1617_SEASON_ID"`
	PremierLeague1516 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_1516_SEASON_ID"`
	PremierLeague1415 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_1415_SEASON_ID"`
	PremierLeague1314 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_1314_SEASON_ID"`
	PremierLeague1213 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_1213_SEASON_ID"`
	PremierLeague1112 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_1112_SEASON_ID"`
	PremierLeague1011 string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_1011_SEASON_ID"`
}
type LeagueEndpoints struct {
	LeagueRoundMatches string `mapstructure:"SOFASCOREAPI_LEAGUE_ROUND_MATCHES_ENDPOINT"`
	LaLigaId           string `mapstructure:"SOFASCOREAPI_LALIGA_ID"`
	PremierLeagueId    string `mapstructure:"SOFASCOREAPI_PREMIERLEAGUE_ID"`
}
type MatchEndpoints struct {
	MatchLineups    string `mapstructure:"SOFASCOREAPI_MATCH_LINEUPS_ENDPOINT"`
	MatchH2hHistory string `mapstructure:"SOFASCOREAPI_MATCH_H2H_HISTORY_ENDPOINT"`
}

type TeamEndpoints struct {
	TeamOverallStats string `mapstructure:"SOFASCOREAPI_TEAM_OVERALL_STATS_ENDPOINT"`
	TeamMatchStats   string `mapstructure:"SOFASCOREAPI_TEAM_MATCH_STATS_ENDPOINT"`
}
type PlayerEndpoints struct {
	PlayersStats       string `mapstructure:"SOFASCOREAPI_TEAM_PLAYERS_STATS_ENDPOINT"`
	PlayerSeasonsStats string `mapstructure:"SOFASCOREAPI_PLAYER_SEASONS_STATS_ENDPOINT"`
	PlayerAttributes   string `mapstructure:"SOFASCOREAPI_PLAYER_ATTRIBUTES_ENDPOINT"`
}

type TopicsName struct {
	LeagueRoundMatches string `mapstructure:"TOPICSNAME_SOFASCORE_LEAGUE_ROUND_MATCHES"`
	MatchLineups       string `mapstructure:"TOPICSNAME_SOFASCORE_MATCH_LINEUPS"`
	MatchH2hHistory    string `mapstructure:"TOPICSNAME_SOFASCORE_MATCH_H2H_HISTORY"`
	TeamOverallStats   string `mapstructure:"TOPICSNAME_SOFASCORE_TEAM_OVERALL_STATS"`
	TeamMatchStats     string `mapstructure:"TOPICSNAME_SOFASCORE_TEAM_MATCH_STATS"`
	PlayerTeamStats    string `mapstructure:"TOPICSNAME_SOFASCORE_PLAYER_TEAM_STATS"`
	PlayerSeasonsStats string `mapstructure:"TOPICSNAME_SOFASCORE_PLAYER_SEASONS_STATS"`
	PlayerAttributes   string `mapstructure:"TOPICSNAME_SOFASCORE_PLAYER_ATTRIBUTES"`
}

type ConsumersGroupId struct {
	LeagueRoundMatches string `mapstructure:"CONSUMERSGROUPID_SOFASCORE_LEAGUE_ROUND_MATCHES"`
	MatchLineups       string `mapstructure:"CONSUMERSGROUPID_SOFASCORE_MATCH_LINEUPS"`
	MatchH2hHistory    string `mapstructure:"CONSUMERSGROUPID_SOFASCORE_MATCH_H2H_HISTORY"`
	TeamOverallStats   string `mapstructure:"CONSUMERSGROUPID_SOFASCORE_TEAM_OVERALL_STATS"`
	TeamMatchStats     string `mapstructure:"CONSUMERSGROUPID_SOFASCORE_TEAM_MATCH_STATS"`
	PlayerTeamStats    string `mapstructure:"CONSUMERSGROUPID_SOFASCORE_PLAYER_TEAM_STATS"`
	PlayerSeasonsStats string `mapstructure:"CONSUMERSGROUPID_SOFASCORE_PLAYER_SEASONS_STATS"`
	PlayerAttributes   string `mapstructure:"CONSUMERSGROUPID_SOFASCORE_PLAYER_ATTRIBUTES"`
}

func LoadConfig() *Config {
	_, filename, _, _ := runtime.Caller(0)
	ConfigDir := filepath.Dir(filename)
	RootDir := filepath.Dir(ConfigDir)

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(RootDir)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("sofascore-service: Error reading config file, %s", err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("sofascore-service: Unable to decode config into struct, %s", err)
	}

	return config
}
