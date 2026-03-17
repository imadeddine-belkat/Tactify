package fpl_models

// Game configuration and scoring rules

// GameSettings contains general game configuration
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

// Scoring represents the point system for different actions
type Scoring struct {
	LongPlay                 int            `json:"long_play"`
	ShortPlay                int            `json:"short_play"`
	GoalsConceded            map[string]int `json:"goals_conceded"`
	Saves                    int            `json:"saves"`
	GoalsScored              map[string]int `json:"goals_scored"`
	Assists                  int            `json:"assists"`
	CleanSheets              map[string]int `json:"clean_sheets"`
	PenaltiesSaved           int            `json:"penalties_saved"`
	PenaltiesMissed          int            `json:"penalties_missed"`
	YellowCards              int            `json:"yellow_cards"`
	RedCards                 int            `json:"red_cards"`
	OwnGoals                 int            `json:"own_goals"`
	Bonus                    int            `json:"bonus"`
	BPS                      int            `json:"bps"`
	Influence                int            `json:"influence"`
	Creativity               int            `json:"creativity"`
	Threat                   int            `json:"threat"`
	ICTIndex                 int            `json:"ict_index"`
	SpecialMultiplier        int            `json:"special_multiplier"`
	Tackles                  int            `json:"tackles"`
	CBI                      int            `json:"clearances_blocks_interceptions"`
	Recoveries               int            `json:"recoveries"`
	DefensiveContribution    map[string]int `json:"defensive_contribution"`
	MngGoalsScored           map[string]int `json:"mng_goals_scored"`
	MngCleanSheets           map[string]int `json:"mng_clean_sheets"`
	MngWin                   map[string]int `json:"mng_win"`
	MngDraw                  map[string]int `json:"mng_draw"`
	MngLoss                  int            `json:"mng_loss"`
	MngUnderdogWin           map[string]int `json:"mng_underdog_win"`
	MngUnderdogDraw          map[string]int `json:"mng_underdog_draw"`
	Starts                   int            `json:"starts"`
	ExpectedAssists          int            `json:"expected_assists"`
	ExpectedGoalInvolvements int            `json:"expected_goal_involvements"`
	ExpectedGoalsConceded    int            `json:"expected_goals_conceded"`
	ExpectedGoals            int            `json:"expected_goals"`
}
