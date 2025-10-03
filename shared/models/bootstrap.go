package models

// BootstrapResponse represents the main FPL api bootstrap-static response
type BootstrapResponse struct {
	Events       []Event           `json:"events"`
	GameSettings GameSettings      `json:"game_settings"`
	Phases       []Phase           `json:"phases"`
	Teams        []Team            `json:"teams"`
	TotalPlayers int               `json:"total_players"`
	Elements     []PlayerBootstrap `json:"elements"`
	ElementStats []ElementStat     `json:"element_stats"`
	ElementTypes []ElementType     `json:"element_types"`
}

// Event represents a gameweek
type Event struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	DeadlineTime      string `json:"deadline_time"`
	AverageEntryScore int    `json:"average_entry_score"`
	Finished          bool   `json:"finished"`
	DataChecked       bool   `json:"data_checked"`
	HighestScore      int    `json:"highest_score"`
	IsPrevious        bool   `json:"is_previous"`
	IsCurrent         bool   `json:"is_current"`
	IsNext            bool   `json:"is_next"`
	ChipPlays         []Chip `json:"chip_plays"`
	MostSelected      int    `json:"most_selected"`
	MostTransferredIn int    `json:"most_transferred_in"`
	TopElement        int    `json:"top_element"`
	TopElementInfo    struct {
		ID     int `json:"id"`
		Points int `json:"points"`
	} `json:"top_element_info"`
	TransfersThisGameweek int `json:"transfers_made"`
	MostCaptained         int `json:"most_captained"`
	MostViceCaptained     int `json:"most_vice_captained"`
}

// GameSettings contains general game information
type GameSettings struct {
	LeagueJoinPrivateMax         int           `json:"league_join_private_max"`
	LeagueJoinPublicMax          int           `json:"league_join_public_max"`
	LeagueMaxSizePublicClassic   int           `json:"league_max_size_public_classic"`
	LeagueMaxSizePublicH2h       int           `json:"league_max_size_public_h2h"`
	LeagueMaxSizePrivateH2h      int           `json:"league_max_size_private_h2h"`
	LeagueMaxKoRounds            int           `json:"league_max_ko_rounds"`
	LeaguePrefixPublic           string        `json:"league_prefix_public"`
	LeaguePointsH2hWin           int           `json:"league_points_h2h_win"`
	LeaguePointsH2hLose          int           `json:"league_points_h2h_lose"`
	LeaguePointsH2hDraw          int           `json:"league_points_h2h_draw"`
	LeagueKoFirstInsteadOfRandom bool          `json:"league_ko_first_instead_of_random"`
	CupStartEventId              int           `json:"cup_start_event_id"`
	CupStopEventId               int           `json:"cup_stop_event_id"`
	CupQualifyingMethod          string        `json:"cup_qualifying_method"`
	CupType                      string        `json:"cup_type"`
	SquadSquadplay               int           `json:"squad_squadplay"`
	SquadSquadsize               int           `json:"squad_squadsize"`
	SquadTeamLimit               int           `json:"squad_team_limit"`
	SquadTotalSpend              int           `json:"squad_total_spend"`
	UICurrencyMultiplier         int           `json:"ui_currency_multiplier"`
	UIUseSpecialShirts           bool          `json:"ui_use_special_shirts"`
	UISpecialShirtExclusions     []interface{} `json:"ui_special_shirt_exclusions"`
	StatsFormDays                int           `json:"stats_form_days"`
	SysViceCaptainEnabled        bool          `json:"sys_vice_captain_enabled"`
	TransfersCap                 int           `json:"transfers_cap"`
	TransfersSellOnFee           float64       `json:"transfers_sell_on_fee"`
	MaxBudget                    int           `json:"max_budget"`
}

// Phase represents a phase of the season
type Phase struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	StartEvent int    `json:"start_event"`
	StopEvent  int    `json:"stop_event"`
}

// Chip represents a chip play
type Chip struct {
	Name      string `json:"name"`
	NumPlayed int    `json:"num_played"`
}

// ElementStat represents element statistics
type ElementStat struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}
