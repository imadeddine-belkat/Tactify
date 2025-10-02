package models

type LiveEvent struct {
	Elements []LiveElement `json:"elements"`
}

type LiveElementDTO struct {
	Event int       `json:"event"`
	Stats LiveStats `json:"stats"`
	// add/remove fields as required
}

type LiveElement struct {
	ID       int           `json:"id" db:"player_id"`
	Stats    LiveStats     `json:"stats"`
	Explain  []ExplainItem `json:"explain"`
	Modified bool          `json:"modified" db:"modified"`
}

type LiveStats struct {
	Minutes         int `json:"minutes" db:"minutes"`
	GoalsScored     int `json:"goals_scored" db:"goals_scored"`
	Assists         int `json:"assists" db:"assists"`
	CleanSheets     int `json:"clean_sheets" db:"clean_sheets"`
	GoalsConceded   int `json:"goals_conceded" db:"goals_conceded"`
	OwnGoals        int `json:"own_goals" db:"own_goals"`
	PenaltiesSaved  int `json:"penalties_saved" db:"penalties_saved"`
	PenaltiesMissed int `json:"penalties_missed" db:"penalties_missed"`
	YellowCards     int `json:"yellow_cards" db:"yellow_cards"`
	RedCards        int `json:"red_cards" db:"red_cards"`
	Saves           int `json:"saves" db:"saves"`
	Bonus           int `json:"bonus" db:"bonus"`
	BPS             int `json:"bps" db:"bps"`

	Influence  string `json:"influence" db:"influence"`
	Creativity string `json:"creativity" db:"creativity"`
	Threat     string `json:"threat" db:"threat"`
	ICTIndex   string `json:"ict_index" db:"ict_index"`

	ClearancesBlocksInterceptions int `json:"clearances_blocks_interceptions" db:"clearances_blocks_interceptions"`
	Recoveries                    int `json:"recoveries" db:"recoveries"`
	Tackles                       int `json:"tackles" db:"tackles"`
	DefensiveContribution         int `json:"defensive_contribution" db:"defensive_contribution"`
	Starts                        int `json:"starts" db:"starts"`

	ExpectedGoals            string `json:"expected_goals" db:"expected_goals"`
	ExpectedAssists          string `json:"expected_assists" db:"expected_assists"`
	ExpectedGoalInvolvements string `json:"expected_goal_involvements" db:"expected_goal_involvements"`
	ExpectedGoalsConceded    string `json:"expected_goals_conceded" db:"expected_goals_conceded"`

	TotalPoints int  `json:"total_points" db:"total_points"`
	InDreamteam bool `json:"in_dreamteam" db:"in_dreamteam"`
}

type ExplainItem struct {
	Fixture int               `json:"fixture" db:"fixture_id"`
	Stats   []ExplainStatItem `json:"stats"`
}

type ExplainStatItem struct {
	Identifier         string `json:"identifier" db:"identifier"`
	Points             int    `json:"points" db:"points"`
	Value              int    `json:"value" db:"value"`
	PointsModification int    `json:"points_modification" db:"points_modification"`
}
