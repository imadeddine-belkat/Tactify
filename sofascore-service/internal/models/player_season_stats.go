package models

type PlayerSeasonsStats struct {
	All []PlayerSeasonStats `json:"playerSeasonStats"`
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
