package sofascore_models

// Top players ranking types

// PlayerInfo Common embedded structs
type PlayerInfo struct {
	Player       Player `json:"player"`
	Team         Team   `json:"team"`
	PlayedEnough bool   `json:"playedEnough"`
}

type BaseStats struct {
	Type        string `json:"type"`
	Appearances int    `json:"appearances"`
}

type PlayerOverallStats[T int | float64] struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Appearances int    `json:"appearances"`

	Rating             float64 `json:"rating,omitempty"`
	Goals              int     `json:"goals,omitempty"`
	Assists            int     `json:"assists,omitempty"`
	ExpectedGoals      float64 `json:"expectedGoals,omitempty"`
	ExpectedAssists    float64 `json:"expectedAssists,omitempty"`
	GoalsAssistsSum    int     `json:"goalsAssistsSum,omitempty"`
	PenaltyGoals       int     `json:"penaltyGoals,omitempty"`
	FreeKickGoal       int     `json:"freeKickGoal,omitempty"`
	ScoringFrequency   float64 `json:"scoringFrequency,omitempty"`
	TotalShots         int     `json:"totalShots,omitempty"`
	ShotsOnTarget      int     `json:"shotsOnTarget,omitempty"`
	BigChancesCreated  int     `json:"bigChancesCreated,omitempty"`
	BigChancesMissed   int     `json:"bigChancesMissed,omitempty"`
	AccuratePasses     int     `json:"accuratePasses,omitempty"`
	KeyPasses          int     `json:"keyPasses,omitempty"`
	AccurateLongBalls  int     `json:"accurateLongBalls,omitempty"`
	SuccessfulDribbles int     `json:"successfulDribbles,omitempty"`
	PenaltiesWon       int     `json:"penaltiesWon,omitempty"`
	Tackles            int     `json:"tackles,omitempty"`
	Interceptions      int     `json:"interceptions,omitempty"`
	Clearances         int     `json:"clearances,omitempty"`
	PossessionsLost    int     `json:"possessionsLost,omitempty"`
	YellowCards        int     `json:"yellowCards,omitempty"`
	RedCards           int     `json:"redCards,omitempty"`
}

type PlayerStat[T int | float64] struct {
	SeasonID   int                   `json:"season"`
	LeagueID   int                   `json:"league"`
	PlayerInfo PlayerInfo            `json:"player"`
	Statistics PlayerOverallStats[T] `json:"statistics"`
}

// TopPlayers struct
type TopPlayers struct {
	Ratings            []TopPlayersRatings            `json:"rating"`
	Goals              []TopPlayersGoals              `json:"goals"`
	Assists            []TopPlayersAssists            `json:"assists"`
	ExpectedGoals      []TopPlayersExpectedGoals      `json:"expectedGoals"`
	ExpectedAssists    []TopPlayersExpectedAssists    `json:"expectedAssists"`
	GoalsAssistsSum    []TopPlayersGoalsAssistsSum    `json:"goalsAssistsSum"`
	PenaltyGoals       []TopPlayersPenaltyGoals       `json:"penaltyGoals"`
	FreeKickGoal       []TopPlayersFreeKickGoal       `json:"freeKickGoal"`
	ScoringFrequency   []TopPlayersScoringFrequency   `json:"scoringFrequency"`
	TotalShots         []TopPlayersTotalShots         `json:"totalShots"`
	ShotsOnTarget      []TopPlayersShotsOnTarget      `json:"shotsOnTarget"`
	BigChancesCreated  []TopPlayersBigChancesCreated  `json:"bigChancesCreated"`
	BigChancesMissed   []TopPlayersBigChancesMissed   `json:"bigChancesMissed"`
	AccuratePasses     []TopPlayersAccuratePasses     `json:"accuratePasses"`
	KeyPasses          []TopPlayersKeyPasses          `json:"keyPasses"`
	AccurateLongBalls  []TopPlayersAccurateLongBalls  `json:"accurateLongBalls"`
	SuccessfulDribbles []TopPlayersSuccessfulDribbles `json:"successfulDribbles"`
	PenaltiesWon       []TopPlayersPenaltiesWon       `json:"penaltiesWon"`
	Tackles            []TopPlayersTackles            `json:"tackles"`
	Interceptions      []TopPlayersInterceptions      `json:"interceptions"`
	Clearances         []TopPlayersClearances         `json:"clearances"`
	PossessionsLost    []TopPlayersPossessionsLost    `json:"possessionsLost"`
	YellowCards        []TopPlayersYellowCards        `json:"yellowCards"`
	RedCards           []TopPlayersRedCards           `json:"redCards"`
}

type TopPlayersRatings = []*PlayerStat[float64]
type TopPlayersGoals = []*PlayerStat[int]
type TopPlayersAssists = []*PlayerStat[int]
type TopPlayersExpectedGoals = []*PlayerStat[float64]
type TopPlayersExpectedAssists = []*PlayerStat[float64]
type TopPlayersGoalsAssistsSum = []*PlayerStat[int]
type TopPlayersPenaltyGoals = []*PlayerStat[int]
type TopPlayersFreeKickGoal = []*PlayerStat[int]
type TopPlayersScoringFrequency = []*PlayerStat[float64]
type TopPlayersTotalShots = []*PlayerStat[int]
type TopPlayersShotsOnTarget = []*PlayerStat[int]
type TopPlayersBigChancesCreated = []*PlayerStat[int]
type TopPlayersBigChancesMissed = []*PlayerStat[int]
type TopPlayersAccuratePasses = []*PlayerStat[int]
type TopPlayersKeyPasses = []*PlayerStat[int]
type TopPlayersAccurateLongBalls = []*PlayerStat[int]
type TopPlayersSuccessfulDribbles = []*PlayerStat[int]
type TopPlayersPenaltiesWon = []*PlayerStat[int]
type TopPlayersTackles = []*PlayerStat[int]
type TopPlayersInterceptions = []*PlayerStat[int]
type TopPlayersClearances = []*PlayerStat[int]
type TopPlayersPossessionsLost = []*PlayerStat[int]
type TopPlayersYellowCards = []*PlayerStat[int]
type TopPlayersRedCards = []*PlayerStat[int]
