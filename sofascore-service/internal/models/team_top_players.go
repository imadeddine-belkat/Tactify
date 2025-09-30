package models

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

// Core types
type Player struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Position  string `json:"position"`
}

// TopPlayers struct
type TopPlayers struct {
	Ratings            Ratings            `json:"ratings"`
	Goals              Goals              `json:"goals"`
	Assists            Assists            `json:"assists"`
	ExpectedGoals      ExpectedGoals      `json:"expectedGoals"`
	ExpectedAssists    ExpectedAssists    `json:"expectedAssists"`
	GoalsAssistsSum    GoalsAssistsSum    `json:"goalsAssistsSum"`
	PenaltyGoals       PenaltyGoals       `json:"penaltyGoals"`
	FreeKickGoals      FreeKickGoals      `json:"freeKickGoals"`
	ScoringFrequencies ScoringFrequencies `json:"scoringFrequencies"`
	TotalShots         TotalShots         `json:"totalShots"`
	ShotsOnTargets     ShotsOnTargets     `json:"shotsOnTargets"`
	BigChancesCreated  BigChancesCreated  `json:"bigChancesCreated"`
	BigChancesMissed   BigChancesMissed   `json:"bigChancesMissed"`
	AccuratePasses     AccuratePasses     `json:"accuratePasses"`
	KeyPasses          KeyPasses          `json:"keyPasses"`
	AccurateLongBalls  AccurateLongBalls  `json:"accurateLongBalls"`
	SuccessfulDribbles SuccessfulDribbles `json:"successfulDribbles"`
	PenaltiesWon       PenaltiesWon       `json:"penaltiesWon"`
	Tackles            Tackles            `json:"tackles"`
	Interceptions      Interceptions      `json:"interceptions"`
	Clearances         Clearances         `json:"clearances"`
	PossessionsLost    PossessionsLost    `json:"possessionsLost"`
	YellowCards        YellowCards        `json:"yellowCards"`
	RedCards           RedCards           `json:"redCards"`
}

// Ratings
type Ratings struct {
	All []Rating `json:"all"`
}

type Rating struct {
	PlayerInfo
	RatingStats RatingStats `json:"statistics"`
}

type RatingStats struct {
	BaseStats
	Rating int `json:"rating"`
}

// Goals
type Goals struct {
	All []Goal `json:"all"`
}

type Goal struct {
	PlayerInfo
	GoalStats GoalStats `json:"statistics"`
}

type GoalStats struct {
	BaseStats
	Goals int `json:"goals"`
}

// Assists
type Assists struct {
	All []Assist `json:"all"`
}

type Assist struct {
	PlayerInfo
	AssistStats AssistStats `json:"statistics"`
}

type AssistStats struct {
	BaseStats
	Assists int `json:"assists"`
}

// ExpectedGoals
type ExpectedGoals struct {
	All []ExpectedGoal `json:"all"`
}

type ExpectedGoal struct {
	PlayerInfo
	ExpectedGoalStats ExpectedGoalStats `json:"statistics"`
}

type ExpectedGoalStats struct {
	BaseStats
	Goals         int     `json:"goals"`
	ExpectedGoals float64 `json:"expectedGoals"`
}

// ExpectedAssists
type ExpectedAssists struct {
	All []ExpectedAssist `json:"all"`
}

type ExpectedAssist struct {
	PlayerInfo
	ExpectedAssistStats ExpectedAssistStats `json:"statistics"`
}

type ExpectedAssistStats struct {
	BaseStats
	Assists         int     `json:"assists"`
	ExpectedAssists float64 `json:"expectedAssists"`
}

// GoalsAssistsSum
type GoalsAssistsSum struct {
	All []GoalAssistSum `json:"all"`
}

type GoalAssistSum struct {
	PlayerInfo
	GoalAssistSumStats GoalAssistSumStats `json:"statistics"`
}

type GoalAssistSumStats struct {
	BaseStats
	GoalsAssistsSum int `json:"goalsAssistsSum"`
}

// PenaltyGoals
type PenaltyGoals struct {
	All []PenaltyGoal `json:"all"`
}

type PenaltyGoal struct {
	PlayerInfo
	PenaltyGoalStats PenaltyGoalStats `json:"statistics"`
}

type PenaltyGoalStats struct {
	BaseStats
	PenaltiesTaken int `json:"penaltiesTaken"`
	PenaltyGoals   int `json:"penaltyGoals"`
}

// FreeKickGoals
type FreeKickGoals struct {
	All []FreeKickGoal `json:"all"`
}

type FreeKickGoal struct {
	PlayerInfo
	FreeKickGoalStats FreeKickGoalStats `json:"statistics"`
}

type FreeKickGoalStats struct {
	BaseStats
	ShotFromSetPiece int `json:"shotFromSetPiece"`
	FreeKickGoal     int `json:"freeKickGoal"`
}

// ScoringFrequencies
type ScoringFrequencies struct {
	All []ScoringFrequency `json:"all"`
}

type ScoringFrequency struct {
	PlayerInfo
	ScoringFrequencyStats ScoringFrequencyStats `json:"statistics"`
}

type ScoringFrequencyStats struct {
	BaseStats
	ScoringFrequency float64 `json:"scoringFrequency"`
}

// TotalShots
type TotalShots struct {
	All []TotalShot `json:"all"`
}

type TotalShot struct {
	PlayerInfo
	TotalShotStats TotalShotStats `json:"statistics"`
}

type TotalShotStats struct {
	BaseStats
	TotalShots int `json:"totalShots"`
}

// ShotsOnTargets
type ShotsOnTargets struct {
	All []ShotsOnTarget `json:"all"`
}

type ShotsOnTarget struct {
	PlayerInfo
	ShotsOnTargetStats ShotsOnTargetStats `json:"statistics"`
}

type ShotsOnTargetStats struct {
	BaseStats
	ShotsOnTarget int `json:"shotsOnTarget"`
}

// BigChancesCreated
type BigChancesCreated struct {
	All []BigChanceCreated `json:"all"`
}

type BigChanceCreated struct {
	PlayerInfo
	BigChanceCreatedStats BigChanceCreatedStats `json:"statistics"`
}

type BigChanceCreatedStats struct {
	BaseStats
	BigChancesCreated int `json:"bigChancesCreated"`
}

// BigChancesMissed
type BigChancesMissed struct {
	All []BigChanceMissed `json:"all"`
}

type BigChanceMissed struct {
	PlayerInfo
	BigChanceMissedStats BigChanceMissedStats `json:"statistics"`
}

type BigChanceMissedStats struct {
	BaseStats
	BigChancesMissed int `json:"bigChancesMissed"`
}

// AccuratePasses
type AccuratePasses struct {
	All []AccuratePass `json:"all"`
}

type AccuratePass struct {
	PlayerInfo
	AccuratePassStats AccuratePassStats `json:"statistics"`
}

type AccuratePassStats struct {
	BaseStats
	AccuratePasses           int     `json:"accuratePasses"`
	AccuratePassesPercentage float64 `json:"accuratePassesPercentage"`
}

// KeyPasses
type KeyPasses struct {
	All []KeyPass `json:"all"`
}

type KeyPass struct {
	PlayerInfo
	KeyPassStats KeyPassStats `json:"statistics"`
}

type KeyPassStats struct {
	BaseStats
	KeyPasses int `json:"keyPasses"`
}

// AccurateLongBalls
type AccurateLongBalls struct {
	All []AccurateLongBall `json:"all"`
}

type AccurateLongBall struct {
	PlayerInfo
	AccurateLongBallStats AccurateLongBallStats `json:"statistics"`
}

type AccurateLongBallStats struct {
	BaseStats
	AccurateLongBalls int `json:"accurateLongBalls"`
}

// SuccessfulDribbles
type SuccessfulDribbles struct {
	All []SuccessfulDribble `json:"all"`
}

type SuccessfulDribble struct {
	PlayerInfo
	SuccessfulDribbleStats SuccessfulDribbleStats `json:"statistics"`
}

type SuccessfulDribbleStats struct {
	BaseStats
	SuccessfulDribbles           int     `json:"successfulDribbles"`
	SuccessfulDribblesPercentage float64 `json:"successfulDribblesPercentage"`
}

// PenaltiesWon
type PenaltiesWon struct {
	All []PenaltyWon `json:"all"`
}

type PenaltyWon struct {
	PlayerInfo
	PenaltyWonStats PenaltyWonStats `json:"statistics"`
}

type PenaltyWonStats struct {
	BaseStats
	PenaltyWon int `json:"penaltyWon"`
}

// Tackles
type Tackles struct {
	All []Tackle `json:"all"`
}

type Tackle struct {
	PlayerInfo
	TackleStats TackleStats `json:"statistics"`
}

type TackleStats struct {
	BaseStats
	Tackles int `json:"tackles"`
}

// Interceptions
type Interceptions struct {
	All []Interception `json:"all"`
}

type Interception struct {
	PlayerInfo
	InterceptionStats InterceptionStats `json:"statistics"`
}

type InterceptionStats struct {
	BaseStats
	Interceptions int `json:"interceptions"`
}

// Clearances
type Clearances struct {
	All []Clearance `json:"all"`
}

type Clearance struct {
	PlayerInfo
	ClearanceStats ClearanceStats `json:"statistics"`
}

type ClearanceStats struct {
	BaseStats
	Clearances int `json:"clearances"`
}

// PossessionsLost
type PossessionsLost struct {
	All []PossessionLost `json:"all"`
}

type PossessionLost struct {
	PlayerInfo
	PossessionLostStats PossessionLostStats `json:"statistics"`
}

type PossessionLostStats struct {
	BaseStats
	PossessionLost int `json:"possessionLost"`
}

// YellowCards
type YellowCards struct {
	All []YellowCard `json:"all"`
}

type YellowCard struct {
	PlayerInfo
	YellowCardStats YellowCardStats `json:"statistics"`
}

type YellowCardStats struct {
	BaseStats
	YellowCards int `json:"yellowCards"`
}

// RedCards
type RedCards struct {
	All []RedCard `json:"all"`
}

type RedCard struct {
	PlayerInfo
	RedCardStats RedCardStats `json:"statistics"`
}

type RedCardStats struct {
	BaseStats
	RedCards int `json:"redCards"`
}
