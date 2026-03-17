package sofascore_models

type TopPlayers struct {
	TopPlayers TopPlayersStats `json:"topPlayers"`
}
type PlayerStat struct {
	Player       Player     `json:"player"`
	Team         Team       `json:"team"`
	PlayedEnough bool       `json:"playedEnough"`
	Statistics   Statistics `json:"statistics"`
}

type TopPlayerStatistics struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Appearances int    `json:"appearances"`

	// Decimal Stats
	Rating           float64 `json:"rating,omitempty"`
	ExpectedGoals    float64 `json:"expectedGoals,omitempty"`
	ExpectedAssists  float64 `json:"expectedAssists,omitempty"`
	ScoringFrequency float64 `json:"scoringFrequency,omitempty"`

	// Integer Stats
	Goals              int `json:"goals,omitempty"`
	Assists            int `json:"assists,omitempty"`
	GoalsAssistsSum    int `json:"goalsAssistsSum,omitempty"`
	PenaltyGoals       int `json:"penaltyGoals,omitempty"`
	FreeKickGoal       int `json:"freeKickGoal,omitempty"`
	TotalShots         int `json:"totalShots,omitempty"`
	ShotsOnTarget      int `json:"shotsOnTarget,omitempty"`
	BigChancesCreated  int `json:"bigChancesCreated,omitempty"`
	BigChancesMissed   int `json:"bigChancesMissed,omitempty"`
	AccuratePasses     int `json:"accuratePasses,omitempty"`
	KeyPasses          int `json:"keyPasses,omitempty"`
	AccurateLongBalls  int `json:"accurateLongBalls,omitempty"`
	SuccessfulDribbles int `json:"successfulDribbles,omitempty"`
	PenaltiesWon       int `json:"penaltiesWon,omitempty"`
	Tackles            int `json:"tackles,omitempty"`
	Interceptions      int `json:"interceptions,omitempty"`
	Clearances         int `json:"clearances,omitempty"`
	PossessionsLost    int `json:"possessionsLost,omitempty"`
	YellowCards        int `json:"yellowCards,omitempty"`
	RedCards           int `json:"redCards,omitempty"`
}

// TopPlayers represents the root container for all leaderboard categories.
type TopPlayersStats struct {
	Rating             []PlayerStat `json:"rating"`
	Goals              []PlayerStat `json:"goals"`
	Assists            []PlayerStat `json:"assists"`
	ExpectedGoals      []PlayerStat `json:"expectedGoals"`
	ExpectedAssists    []PlayerStat `json:"expectedAssists"`
	GoalsAssistsSum    []PlayerStat `json:"goalsAssistsSum"`
	PenaltyGoals       []PlayerStat `json:"penaltyGoals"`
	FreeKickGoal       []PlayerStat `json:"freeKickGoal"`
	ScoringFrequency   []PlayerStat `json:"scoringFrequency"`
	TotalShots         []PlayerStat `json:"totalShots"`
	ShotsOnTarget      []PlayerStat `json:"shotsOnTarget"`
	BigChancesCreated  []PlayerStat `json:"bigChancesCreated"`
	BigChancesMissed   []PlayerStat `json:"bigChancesMissed"`
	AccuratePasses     []PlayerStat `json:"accuratePasses"`
	KeyPasses          []PlayerStat `json:"keyPasses"`
	AccurateLongBalls  []PlayerStat `json:"accurateLongBalls"`
	SuccessfulDribbles []PlayerStat `json:"successfulDribbles"`
	PenaltiesWon       []PlayerStat `json:"penaltiesWon"`
	Tackles            []PlayerStat `json:"tackles"`
	Interceptions      []PlayerStat `json:"interceptions"`
	Clearances         []PlayerStat `json:"clearances"`
	PossessionsLost    []PlayerStat `json:"possessionsLost"`
	YellowCards        []PlayerStat `json:"yellowCards"`
	RedCards           []PlayerStat `json:"redCards"`
}
