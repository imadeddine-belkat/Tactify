package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	kafka "github.com/imadeddine-belkat/kafka/config"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type SofascoreConfig struct {
	SofascoreApi SofascoreApi
	KafkaConfig  kafka.KafkaConfig
	Tor          TorConfig

	// Dynamic season storage - indexed by league name
	Seasons map[string]map[string]int // map[league]map[year]seasonID

	PublishWorkerCount int `envconfig:"WORKER_PUBLISH_POOL_SIZE" default:"100"`
	FetchWorkerCount   int `envconfig:"WORKER_FETCH_POOL_SIZE" default:"50"`
}

type TorConfig struct {
	Enabled     bool   `envconfig:"TOR_ENABLED" default:"false"`
	SocksAddr   string `envconfig:"TOR_SOCKS_ADDR" default:"127.0.0.1:9050"`
	ControlAddr string `envconfig:"TOR_CONTROL_ADDR" default:"127.0.0.1:9051"`
	Password    string `envconfig:"TOR_PASSWORD" default:""`
}

type SofascoreApi struct {
	BaseURL     string `envconfig:"SOFASCOREAPI_BASE_URL"`
	CurrentYear string `envconfig:"CURRENT_YEAR"`
	LeaguesID   LeaguesIDs
	LeagueEndpoints
	MatchEndpoints
	TeamEndpoints
	PlayerEndpoints
}

type LeaguesIDs struct {
	LaLiga        int `envconfig:"SOFASCOREAPI_LALIGA_ID"`
	PremierLeague int `envconfig:"SOFASCOREAPI_PREMIERLEAGUE_ID"`
}

type LeagueEndpoints struct {
	LeagueSeasonsIDs       string `envconfig:"SOFASCOREAPI_LEAGUE_SEASONS_ENDPOINT"`
	LeagueCountriesIDs     string `envconfig:"SOFASCOREAPI_COUNTRIES_IDS_ENDPOINT"`
	LeagueCountryLeagueIDs string `envconfig:"SOFASCOREAPI_COUNTRY_LEAGUES_ENDPOINT"`
	LeagueSeasonStandings  string `envconfig:"SOFASCOREAPI_LEAGUE_SEASON_STANDINGS_ENDPOINT"`
	LeagueRoundMatches     string `envconfig:"SOFASCOREAPI_LEAGUE_ROUND_MATCHES_ENDPOINT"`
}

type MatchEndpoints struct {
	MatchLineups     string `envconfig:"SOFASCOREAPI_MATCH_LINEUPS_ENDPOINT"`
	MatchH2hHistory  string `envconfig:"SOFASCOREAPI_MATCH_H2H_HISTORY_ENDPOINT"`
	MatchBestPlayers string `envconfig:"SOFASCOREAPI_MATCH_BEST_PLAYERS"`
}

type TeamEndpoints struct {
	TopTeamsStats    string `envconfig:"SOFASCOREAPI_TOP_TEAMS_OVERALL_STATS_ENDPOINT"`
	TeamOverallStats string `envconfig:"SOFASCOREAPI_TEAM_OVERALL_STATS_ENDPOINT"`
	TeamMatchStats   string `envconfig:"SOFASCOREAPI_TEAM_MATCH_STATS_ENDPOINT"`
	TeamPlayerStats  string `envconfig:"SOFASCOREAPI_TEAM_PLAYERS_STATS_ENDPOINT"`
}

type PlayerEndpoints struct {
	PlayersStats       string `envconfig:"SOFASCOREAPI_PLAYER_INFO_ENDPOINT"`
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

	// Parse config with validation
	if err := envconfig.Process("", config); err != nil {
		log.Fatalf("sofascore-service: Unable to load API config: %s", err)
	}

	config.KafkaConfig = *kafka.LoadConfig()

	// Parse all season IDs dynamically
	config.Seasons = parseSeasonIDs()

	return config
}

// parseSeasonIDs extracts season IDs from environment variables
// Pattern: SOFASCOREAPI_{LEAGUE}_{YEAR}_SEASON_ID
func parseSeasonIDs() map[string]map[string]int {
	seasons := make(map[string]map[string]int)

	// Regex to match: SOFASCOREAPI_LALIGA_2425_SEASON_ID=61643
	pattern := regexp.MustCompile(`^SOFASCOREAPI_([A-Z]+)_(\d{4})_SEASON_ID$`)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]
		matches := pattern.FindStringSubmatch(key)

		if len(matches) == 3 {
			league := strings.ToUpper(matches[1]) // LALIGA, PREMIERLEAGUE
			year := matches[2]                    // 2425, 2324, etc.

			seasonID, err := strconv.Atoi(value)
			if err != nil {
				log.Printf("Warning: Invalid season ID for %s: %v", key, err)
				continue
			}

			if seasons[league] == nil {
				seasons[league] = make(map[string]int)
			}
			seasons[league][year] = seasonID
		}
	}

	return seasons
}

// GetSeasonID retrieves a season ID for a given league and year
// Example: config.GetSeasonID("LALIGA", "2425") -> 61643
func (c *SofascoreConfig) GetSeasonID(league, year string) (int, error) {
	league = strings.ToUpper(league)

	if leagueSeasons, ok := c.Seasons[league]; ok {
		if seasonID, ok := leagueSeasons[year]; ok {
			return seasonID, nil
		}
		return 0, fmt.Errorf("season %s not found for league %s", year, league)
	}
	return 0, fmt.Errorf("league %s not found", league)
}

// AllSeasons returns all season IDs for a given league, sorted by year (descending)
func (c *SofascoreConfig) AllSeasons(league string) []int {
	league = strings.ToUpper(league)

	seasons, ok := c.Seasons[league]
	if !ok {
		return []int{}
	}

	// Extract and sort years
	years := make([]string, 0, len(seasons))
	for year := range seasons {
		years = append(years, year)
	}

	// Sort descending (most recent first)
	for i := 0; i < len(years)-1; i++ {
		for j := i + 1; j < len(years); j++ {
			if years[i] < years[j] {
				years[i], years[j] = years[j], years[i]
			}
		}
	}

	// Build result
	result := make([]int, len(years))
	for i, year := range years {
		result[i] = seasons[year]
	}

	return result
}

// MustGetSeasonID panics if season not found (use for critical paths)
func (c *SofascoreConfig) MustGetSeasonID(league, year string) int {
	id, err := c.GetSeasonID(league, year)
	if err != nil {
		log.Fatalf("Critical: %v", err)
	}
	return id
}

// Legacy helper methods for backward compatibility (if needed)
func (c *SofascoreConfig) LaLigaSeasonID(year string) int {
	id, _ := c.GetSeasonID("LALIGA", year)
	return id
}

func (c *SofascoreConfig) PremierLeagueSeasonID(year string) int {
	id, _ := c.GetSeasonID("PREMIERLEAGUE", year)
	return id
}
