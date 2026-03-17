package sofascore_models

// Top teams ranking types with generic support

// GetTeamID returns the team ID - this makes it easy to use with generics
func (t *TeamStat[T]) GetTeamID() int {
	return t.Team.ID
}

func (t *TeamStat[T]) SetSeasonID(seasonId int) {
	t.SeasonID = seasonId
}

func (t *TeamStat[T]) SetLeagueID(leagueId int) {
	t.LeagueID = leagueId
}

type TeamStatistics[T int | float64] struct {
	ID             int `json:"id"`
	Matches        int `json:"matches"`
	AwardedMatches int `json:"awardedMatches"`

	// Stat-specific fields (only one will be used per type)
	AverageRating         float64 `json:"avgRating,omitempty"`
	GoalsScored           int     `json:"goalsScored,omitempty"`
	GoalsConceded         int     `json:"goalsConceded,omitempty"`
	BigChances            int     `json:"bigChances,omitempty"`
	BigChancesMissed      int     `json:"bigChancesMissed,omitempty"`
	HitWoodWork           int     `json:"hitWoodWork,omitempty"`
	YellowCards           int     `json:"yellowCards,omitempty"`
	RedCards              int     `json:"redCards,omitempty"`
	AverageBallPossession float64 `json:"avgBallPossession,omitempty"`
	AccuratePasses        int     `json:"accuratePasses,omitempty"`
	AccurateLongBalls     int     `json:"accurateLongBalls,omitempty"`
	AccurateCrosses       int     `json:"accurateCrosses,omitempty"`
	Shots                 int     `json:"shots,omitempty"`
	ShotsOnTarget         int     `json:"shotsOnTarget,omitempty"`
	SuccessfulDribbles    int     `json:"successfulDribbles,omitempty"`
	Tackles               int     `json:"tackles,omitempty"`
	Interceptions         int     `json:"interceptions,omitempty"`
	Clearances            int     `json:"clearances,omitempty"`
	Corners               int     `json:"corners,omitempty"`
	Fouls                 int     `json:"fouls,omitempty"`
	PenaltyGoals          int     `json:"penaltyGoals,omitempty"`
	PenaltyGoalsConceded  int     `json:"penaltyGoalsConceded,omitempty"`
	CleanSheets           int     `json:"cleanSheets,omitempty"`
}

type TeamStat[T int | float64] struct {
	SeasonID   int               `json:"season"`
	LeagueID   int               `json:"league"`
	Team       Team              `json:"team"`
	Statistics TeamStatistics[T] `json:"statistics"`
}

type TopTeams struct {
	AverageRating         []TopTeamsAverageRating         `json:"avgRating"`
	GoalsScored           []TopTeamsGoalsScored           `json:"goalsScored"`
	GoalsConceded         []TopTeamsGoalsConceded         `json:"goalsConceded"`
	BigChances            []TopTeamsBigChances            `json:"bigChances"`
	BigChancesMissed      []TopTeamsBigChancesMissed      `json:"bigChancesMissed"`
	HitWoodWork           []TopTeamsHitWoodWork           `json:"hitWoodWork"`
	YellowCards           []TopTeamsYellowCards           `json:"yellowCards"`
	RedCards              []TopTeamsRedCards              `json:"redCards"`
	AverageBallPossession []TopTeamsAverageBallPossession `json:"averageBallPossession"`
	AccuratePasses        []TopTeamsAccuratePasses        `json:"accuratePasses"`
	AccurateLongBalls     []TopTeamsAccurateLongBalls     `json:"accurateLongBalls"`
	AccurateCrosses       []TopTeamsAccurateCrosses       `json:"accurateCrosses"`
	Shots                 []TopTeamsShots                 `json:"shots"`
	ShotsOnTarget         []TopTeamsShotsOnTarget         `json:"shotsOnTarget"`
	SuccessfulDribbles    []TopTeamsSuccessfulDribbles    `json:"successfulDribbles"`
	Tackles               []TopTeamsTackles               `json:"tackles"`
	Interceptions         []TopTeamsInterceptions         `json:"interceptions"`
	Clearances            []TopTeamsClearances            `json:"clearances"`
	Corners               []TopTeamsCorners               `json:"corners"`
	Fouls                 []TopTeamsFouls                 `json:"fouls"`
	PenaltyGoals          []TopTeamsPenaltyGoals          `json:"penaltyGoals"`
	PenaltyGoalsConceded  []TopTeamsPenaltyGoalsConceded  `json:"penaltyGoalsConceded"`
	CleanSheets           []TopTeamsCleanSheets           `json:"cleanSheets"`
}

// Type aliases for specific statistics
type TopTeamsAverageRating = *TeamStat[float64]
type TopTeamsGoalsScored = *TeamStat[int]
type TopTeamsGoalsConceded = *TeamStat[int]
type TopTeamsBigChances = *TeamStat[int]
type TopTeamsBigChancesMissed = *TeamStat[int]
type TopTeamsHitWoodWork = *TeamStat[int]
type TopTeamsYellowCards = *TeamStat[int]
type TopTeamsRedCards = *TeamStat[int]
type TopTeamsAverageBallPossession = *TeamStat[float64]
type TopTeamsAccuratePasses = *TeamStat[int]
type TopTeamsAccurateLongBalls = *TeamStat[int]
type TopTeamsAccurateCrosses = *TeamStat[int]
type TopTeamsShots = *TeamStat[int]
type TopTeamsShotsOnTarget = *TeamStat[int]
type TopTeamsSuccessfulDribbles = *TeamStat[int]
type TopTeamsTackles = *TeamStat[int]
type TopTeamsInterceptions = *TeamStat[int]
type TopTeamsClearances = *TeamStat[int]
type TopTeamsCorners = *TeamStat[int]
type TopTeamsFouls = *TeamStat[int]
type TopTeamsPenaltyGoals = *TeamStat[int]
type TopTeamsPenaltyGoalsConceded = *TeamStat[int]
type TopTeamsCleanSheets = *TeamStat[int]
