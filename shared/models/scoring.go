package models

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
