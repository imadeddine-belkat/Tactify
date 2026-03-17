package fpl_models

// Player data types

type Player struct {
	PlayerHistory []PlayerHistory     `json:"history"`
	PlayerPast    []PlayerPastHistory `json:"history_past"`
}

type PlayerBootstrap struct {
	ID         int    `json:"id" db:"player_id"`
	Code       int    `json:"code" db:"player_code"`
	FirstName  string `json:"first_name" db:"first_name"`
	SecondName string `json:"second_name" db:"second_name"`
	WebName    string `json:"web_name" db:"web_name"`

	TeamID      int    `json:"team" db:"team_id"`
	TeamCode    int    `json:"team_code" db:"team_code"`
	ElementType int    `json:"element_type" db:"element_type_id"`
	Status      string `json:"status" db:"status"`

	CanTransact bool `json:"can_transact" db:"can_transact"`
	CanSelect   bool `json:"can_select" db:"can_select"`
	InDreamteam bool `json:"in_dreamteam" db:"in_dreamteam"`
	Special     bool `json:"special" db:"special"`
	Removed     bool `json:"removed" db:"removed"`
	Unavailable bool `json:"unavailable" db:"unavailable"`

	NowCost             int `json:"now_cost" db:"now_cost"`
	CostChangeEvent     int `json:"cost_change_event" db:"cost_change_event"`
	CostChangeEventFall int `json:"cost_change_event_fall" db:"cost_change_event_fall"`
	CostChangeStart     int `json:"cost_change_start" db:"cost_change_start"`
	CostChangeStartFall int `json:"cost_change_start_fall" db:"cost_change_start_fall"`

	DreamteamCount int `json:"dreamteam_count" db:"dreamteam_count"`

	TotalPoints       int    `json:"total_points" db:"total_points"`
	EventPoints       int    `json:"event_points" db:"event_points"`
	PointsPerGame     string `json:"points_per_game" db:"points_per_game"`
	Form              string `json:"form" db:"form"`
	SelectedByPercent string `json:"selected_by_percent" db:"selected_by_percent"`
	ValueForm         string `json:"value_form" db:"value_form"`
	ValueSeason       string `json:"value_season" db:"value_season"`

	Photo        string `json:"photo" db:"photo"`
	SquadNumber  *int   `json:"squad_number" db:"squad_number"`
	BirthDate    string `json:"birth_date" db:"birth_date"`
	TeamJoinDate string `json:"team_join_date" db:"team_join_date"`
	Region       int    `json:"region" db:"region"`
	OptaCode     string `json:"opta_code" db:"opta_code"`

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

	ExpectedGoalsPer90            float64 `json:"expected_goals_per_90" db:"expected_goals_per_90"`
	ExpectedAssistsPer90          float64 `json:"expected_assists_per_90" db:"expected_assists_per_90"`
	ExpectedGoalInvolvementsPer90 float64 `json:"expected_goal_involvements_per_90" db:"expected_goal_involvements_per_90"`
	ExpectedGoalsConcededPer90    float64 `json:"expected_goals_conceded_per_90" db:"expected_goals_conceded_per_90"`
	SavesPer90                    float64 `json:"saves_per_90" db:"saves_per_90"`
	GoalsConcededPer90            float64 `json:"goals_conceded_per_90" db:"goals_conceded_per_90"`
	StartsPer90                   float64 `json:"starts_per_90" db:"starts_per_90"`
	CleanSheetsPer90              float64 `json:"clean_sheets_per_90" db:"clean_sheets_per_90"`
	DefensiveContributionPer90    float64 `json:"defensive_contribution_per_90" db:"defensive_contribution_per_90"`

	// Rankings
	InfluenceRank      int `json:"influence_rank" db:"influence_rank"`
	InfluenceRankType  int `json:"influence_rank_type" db:"influence_rank_type"`
	CreativityRank     int `json:"creativity_rank" db:"creativity_rank"`
	CreativityRankType int `json:"creativity_rank_type" db:"creativity_rank_type"`
	ThreatRank         int `json:"threat_rank" db:"threat_rank"`
	ThreatRankType     int `json:"threat_rank_type" db:"threat_rank_type"`
	ICTIndexRank       int `json:"ict_index_rank" db:"ict_index_rank"`
	ICTIndexRankType   int `json:"ict_index_rank_type" db:"ict_index_rank_type"`

	NowCostRank           int `json:"now_cost_rank" db:"now_cost_rank"`
	NowCostRankType       int `json:"now_cost_rank_type" db:"now_cost_rank_type"`
	FormRank              int `json:"form_rank" db:"form_rank"`
	FormRankType          int `json:"form_rank_type" db:"form_rank_type"`
	PointsPerGameRank     int `json:"points_per_game_rank" db:"points_per_game_rank"`
	PointsPerGameRankType int `json:"points_per_game_rank_type" db:"points_per_game_rank_type"`
	SelectedRank          int `json:"selected_rank" db:"selected_rank"`
	SelectedRankType      int `json:"selected_rank_type" db:"selected_rank_type"`
}

type PlayersBootstrap struct {
	PlayerBootstrap []PlayerBootstrap `json:"elements"`
}
