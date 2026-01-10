package sofascore_models

// Player statistics - match and season level

type PlayerStatistics struct {
	MinutesPlayed                int     `json:"minutesPlayed,omitempty"`
	Rating                       float64 `json:"rating,omitempty"`
	Touches                      int     `json:"touches,omitempty"`
	TotalPass                    int     `json:"totalPass,omitempty"`
	AccuratePass                 int     `json:"accuratePass,omitempty"`
	TotalLongBalls               int     `json:"totalLongBalls,omitempty"`
	AccurateLongBalls            int     `json:"accurateLongBalls,omitempty"`
	AccurateOwnHalfPasses        int     `json:"accurateOwnHalfPasses,omitempty"`
	TotalOwnHalfPasses           int     `json:"totalOwnHalfPasses,omitempty"`
	AccurateOppositionHalfPasses int     `json:"accurateOppositionHalfPasses,omitempty"`
	TotalOppositionHalfPasses    int     `json:"totalOppositionHalfPasses,omitempty"`
	TotalCross                   int     `json:"totalCross,omitempty"`
	AccurateCross                int     `json:"accurateCross,omitempty"`
	KeyPass                      int     `json:"keyPass,omitempty"`
	GoalAssist                   int     `json:"goalAssist,omitempty"`
	BigChanceCreated             int     `json:"bigChanceCreated,omitempty"`
	TotalShots                   int     `json:"totalShots,omitempty"`
	ShotOffTarget                int     `json:"shotOffTarget,omitempty"`
	OnTargetScoringAttempt       int     `json:"onTargetScoringAttempt,omitempty"`
	BlockedScoringAttempt        int     `json:"blockedScoringAttempt,omitempty"`
	HitWoodwork                  int     `json:"hitWoodwork,omitempty"`
	Goals                        int     `json:"goals,omitempty"`
	BigChanceMissed              int     `json:"bigChanceMissed,omitempty"`
	ExpectedGoals                float64 `json:"expectedGoals,omitempty"`
	ExpectedGoalsOnTarget        float64 `json:"expectedGoalsOnTarget,omitempty"`
	ExpectedAssists              float64 `json:"expectedAssists,omitempty"`
	TotalTackle                  int     `json:"totalTackle,omitempty"`
	WonTackle                    int     `json:"wonTackle,omitempty"`
	TotalClearance               int     `json:"totalClearance,omitempty"`
	InterceptionWon              int     `json:"interceptionWon,omitempty"`
	OutfielderBlock              int     `json:"outfielderBlock,omitempty"`
	BallRecovery                 int     `json:"ballRecovery,omitempty"`
	DuelLost                     int     `json:"duelLost,omitempty"`
	DuelWon                      int     `json:"duelWon,omitempty"`
	AerialLost                   int     `json:"aerialLost,omitempty"`
	AerialWon                    int     `json:"aerialWon,omitempty"`
	TotalContest                 int     `json:"totalContest,omitempty"`
	WonContest                   int     `json:"wonContest,omitempty"`
	ChallengeLost                int     `json:"challengeLost,omitempty"`
	PossessionLostCtrl           int     `json:"possessionLostCtrl,omitempty"`
	UnsuccessfulTouch            int     `json:"unsuccessfulTouch,omitempty"`
	Dispossessed                 int     `json:"dispossessed,omitempty"`
	Fouls                        int     `json:"fouls,omitempty"`
	WasFouled                    int     `json:"wasFouled,omitempty"`
	TotalOffside                 int     `json:"totalOffside,omitempty"`
	Saves                        int     `json:"saves,omitempty"`
	TotalKeeperSweeper           int     `json:"totalKeeperSweeper,omitempty"`
	AccurateKeeperSweeper        int     `json:"accurateKeeperSweeper,omitempty"`
	KeeperSaveValue              float64 `json:"keeperSaveValue,omitempty"`
	GoalsPrevented               float64 `json:"goalsPrevented,omitempty"`
	ErrorLeadToAShot             int     `json:"errorLeadToAShot,omitempty"`
	ErrorLeadToAGoal             int     `json:"errorLeadToAGoal,omitempty"`
	ShotValueNormalized          float64 `json:"shotValueNormalized,omitempty"`
	PassValueNormalized          float64 `json:"passValueNormalized,omitempty"`
	DribbleValueNormalized       float64 `json:"dribbleValueNormalized,omitempty"`
	DefensiveValueNormalized     float64 `json:"defensiveValueNormalized,omitempty"`
	GoalkeeperValueNormalized    float64 `json:"goalkeeperValueNormalized,omitempty"`
	RatingVersions               *struct {
		Original    float64 `json:"original"`
		Alternative float64 `json:"alternative"`
	} `json:"ratingVersions,omitempty"`
}

type PlayerSeasonsStats struct {
	All []PlayerSeasonStats `json:"seasons"`
}

type PlayerSeasonStats struct {
	PlayerID  int    `json:"playerId"`
	Year      string `json:"year"`
	StartYear int    `json:"startYear"`
	EndYear   string `json:"endYear"`
	Team      Team   `json:"team"`
	Season    Season `json:"season"`
}

type Statistics struct {
	AccurateLongBalls           int     `json:"accurateLongBalls"`
	AccurateLongBallsPercentage float64 `json:"accurateLongBallsPercentage"`
	AccuratePasses              int     `json:"accuratePasses"`
	AccuratePassesPercentage    float64 `json:"accuratePassesPercentage"`
	AerialDuelsWon              int     `json:"aerialDuelsWon"`
	Assists                     int     `json:"assists"`
	BigChancesCreated           int     `json:"bigChancesCreated"`
	BigChancesMissed            int     `json:"bigChancesMissed"`
	BlockedShots                int     `json:"blockedShots"`
	OutfielderBlocks            int     `json:"outfielderBlocks"`
	CleanSheet                  int     `json:"cleanSheet"`
	DribbledPast                int     `json:"dribbledPast"`
	ErrorLeadToGoal             int     `json:"errorLeadToGoal"`
	ExpectedAssists             float64 `json:"expectedAssists"`
	ExpectedGoals               float64 `json:"expectedGoals"`
	Goals                       int     `json:"goals"`
	GoalsAssistsSum             int     `json:"goalsAssistsSum"`
	GoalsConceded               int     `json:"goalsConceded"`
	Interceptions               int     `json:"interceptions"`
	KeyPasses                   int     `json:"keyPasses"`
	MinutesPlayed               int     `json:"minutesPlayed"`
	PassToAssist                int     `json:"passToAssist"`
	Rating                      float64 `json:"rating"`
	RedCards                    int     `json:"redCards"`
	Saves                       int     `json:"saves"`
	ShotsOnTarget               int     `json:"shotsOnTarget"`
	SuccessfulDribbles          int     `json:"successfulDribbles"`
	Tackles                     int     `json:"tackles"`
	TotalShots                  int     `json:"totalShots"`
	YellowCards                 int     `json:"yellowCards"`
	TotalRating                 int     `json:"totalRating"`
	CountRating                 int     `json:"countRating"`
	TotalLongBalls              int     `json:"totalLongBalls"`
	TotalPasses                 int     `json:"totalPasses"`
	ShotsFromInsideTheBox       int     `json:"shotsFromInsideTheBox"`
	Appearances                 int     `json:"appearances"`
	Type                        string  `json:"type"`
	ID                          int     `json:"id"`
}
