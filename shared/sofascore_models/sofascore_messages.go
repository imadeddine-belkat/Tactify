package sofascore_models

// API message wrappers for event-driven communication

type StandingMessage struct {
	SeasonID int `json:"seasonId"`
	LeagueID int `json:"leagueId"`
	Row      Row `json:"row"`
}

type PlayerMessage struct {
	SeasonID int    `json:"seasonId"`
	LeagueID int    `json:"leagueId"`
	TeamID   int    `json:"teamId"`
	Player   Player `json:"player"`
}

type MatchLineupMessage struct {
	SeasonID    int         `json:"season"`
	LeagueID    int         `json:"league"`
	MatchID     int         `json:"match"`
	Round       int         `json:"round"`
	MatchLineup MatchLineup `json:"lineup"`
}

type PlayerMatchStatsMessage struct {
	PlayerName  string      `json:"player_name"`
	SeasonID    int         `json:"season"`
	LeagueID    int         `json:"league"`
	MatchID     int         `json:"match"`
	Round       int         `json:"round"`
	MatchPlayer MatchPlayer `json:"player"`
}

type MatchStatsMessage struct {
	SeasonId        int          `json:"season_id"`
	LeagueId        int          `json:"league_id"`
	Event           int          `json:"event"`
	MatchID         int          `json:"match_id"`
	HomeTeamID      int          `json:"homeTeamId"`
	AwayTeamID      int          `json:"awayTeamId"`
	GroupName       string       `json:"groupName"`
	MatchStatistics StatsMessage `json:"statistics"`
}

type StatsMessage struct {
	Period         string  `json:"period"`
	Key            string  `json:"key"`
	Name           string  `json:"name"`
	HomeValue      float64 `json:"homeValue"`
	AwayValue      float64 `json:"awayValue"`
	CompareCode    int     `json:"compareCode"` // Which team performed better: 1-home, 2-away, 3-equal, 0-unknown
	HomeTotal      *int    `json:"homeTotal,omitempty"`
	AwayTotal      *int    `json:"awayTotal,omitempty"`
	StatisticsType string  `json:"statisticsType"` // "positive" higher is better, "negative" higher is worse
	RenderType     int     `json:"renderType"`     // 1=raw numbers, 2=percentage bar, 3=fraction with %, 4=% only from fraction
}

type TeamOverallStatsMessage struct {
	TeamID   int              `json:"team"`
	LeagueID int              `json:"league"`
	SeasonID int              `json:"season"`
	Stats    TeamOverallStats `json:"statistics"`
}

type TopTeamsMessage struct {
	TopTeams TopTeams `json:"topTeams"`
}
