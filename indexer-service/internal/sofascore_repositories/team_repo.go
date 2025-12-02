package sofascore_repositories

import (
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/imadeddine-belkat/indexer-service/internal/sofascore_helper"
	"github.com/imadeddine-belkat/shared/sofascore_models"
)

type TeamRepo struct {
	db               *sql.DB
	Helper           *sofascore_helper.Helper
	LeagueStanding   *sofascore_models.StandingMessage
	TeamOverallStats *sofascore_models.TeamOverallStatsMessage
	MatchStats       *sofascore_models.MatchStatsMessage
}

func NewTeamRepo(
	db *sql.DB,
	leagueStanding *sofascore_models.StandingMessage,
	ovrStats *sofascore_models.TeamOverallStatsMessage,
	MatchStats *sofascore_models.MatchStatsMessage) *TeamRepo {
	return &TeamRepo{
		db:               db,
		LeagueStanding:   leagueStanding,
		TeamOverallStats: ovrStats,
		MatchStats:       MatchStats,
	}
}

func (t *TeamRepo) InsertTeamInfo(standing sofascore_models.StandingMessage) error {
	query := sq.Insert("teams").
		Columns(
			"team_id", "name", "league_id", "primary_color", "secondary_color").
		Suffix(
			"ON CONFLICT (team_id, league_id) DO NOTHING ").
		PlaceholderFormat(sq.Dollar)

	query = query.Values(
		standing.Row.Team.ID,
		standing.Row.Team.Name,
		standing.LeagueID,
		standing.Row.Team.Colors.PrimaryColor,
		standing.Row.Team.Colors.SecondaryColor,
	)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = t.db.Exec(sqlQuery, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *TeamRepo) InsertTeamOverallStats(ovrStats sofascore_models.TeamOverallStatsMessage) error {
	stats := ovrStats.Stats

	query := sq.Insert("team_overall_stats").
		Columns(
			// Primary Keys
			"team_id",
			"league_id",
			"season_id",

			// Offensive Stats
			"goals_scored",
			"goals_conceded",
			"own_goals",
			"assists",
			"shots",
			"shots_on_target",
			"shots_off_target",
			"penalty_goals",
			"penalties_taken",
			"free_kick_goals",
			"free_kick_shots",

			// Positional Goals
			"goals_from_inside_box",
			"goals_from_outside_box",
			"headed_goals",
			"left_foot_goals",
			"right_foot_goals",

			// Chances
			"big_chances",
			"big_chances_created",
			"big_chances_missed",

			// Shooting Details
			"shots_from_inside_box",
			"shots_from_outside_box",
			"blocked_scoring_attempt",
			"hit_woodwork",

			// Dribbling
			"successful_dribbles",
			"dribble_attempts",

			// Set Pieces
			"corners",
			"free_kicks",
			"throw_ins",
			"goal_kicks",

			// Fast Breaks
			"fast_breaks",
			"fast_break_goals",
			"fast_break_shots",

			// Possession & Passing
			"average_ball_possession",
			"total_passes",
			"accurate_passes",
			"accurate_passes_percentage",

			// Passing by Zone
			"total_own_half_passes",
			"accurate_own_half_passes",
			"accurate_own_half_passes_percentage",
			"total_opposition_half_passes",
			"accurate_opposition_half_passes",
			"accurate_opposition_half_passes_percentage",

			// Long Balls & Crosses
			"total_long_balls",
			"accurate_long_balls",
			"accurate_long_balls_percentage",
			"total_crosses",
			"accurate_crosses",
			"accurate_crosses_percentage",

			// Defensive Stats
			"clean_sheets",
			"tackles",
			"interceptions",
			"saves",
			"clearances",
			"clearances_off_line",
			"last_man_tackles",
			"ball_recovery",

			// Errors
			"errors_leading_to_goal",
			"errors_leading_to_shot",

			// Penalties
			"penalties_commited",
			"penalty_goals_conceded",

			// Duels
			"total_duels",
			"duels_won",
			"duels_won_percentage",
			"total_ground_duels",
			"ground_duels_won",
			"ground_duels_won_percentage",
			"total_aerial_duels",
			"aerial_duels_won",
			"aerial_duels_won_percentage",

			// Discipline
			"possession_lost",
			"offsides",
			"fouls",
			"yellow_cards",
			"yellow_red_cards",
			"red_cards",

			// Performance
			"avg_rating",
			"matches",
			"awarded_matches",

			// Stats Against (Defensive Perspective)
			"accurate_final_third_passes_against",
			"accurate_opposition_half_passes_against",
			"accurate_own_half_passes_against",
			"accurate_passes_against",
			"big_chances_against",
			"big_chances_created_against",
			"big_chances_missed_against",
			"clearances_against",
			"corners_against",
			"crosses_successful_against",
			"crosses_total_against",
			"dribble_attempts_total_against",
			"dribble_attempts_won_against",
			"errors_leading_to_goal_against",
			"errors_leading_to_shot_against",
			"hit_woodwork_against",
			"interceptions_against",
			"key_passes_against",
			"long_balls_successful_against",
			"long_balls_total_against",
			"offsides_against",
			"red_cards_against",
			"shots_against",
			"shots_blocked_against",
			"shots_from_inside_box_against",
			"shots_from_outside_box_against",
			"shots_off_target_against",
			"shots_on_target_against",
			"blocked_scoring_attempt_against",
			"tackles_against",
			"total_final_third_passes_against",
			"opposition_half_passes_total_against",
			"own_half_passes_total_against",
			"total_passes_against",
			"yellow_cards_against",
		).
		Values(
			// Primary Keys
			ovrStats.TeamID,
			ovrStats.LeagueID,
			ovrStats.SeasonID,

			// Offensive Stats
			stats.GoalsScored,
			stats.GoalsConceded,
			stats.OwnGoals,
			stats.Assists,
			stats.Shots,
			stats.ShotsOnTarget,
			stats.ShotsOffTarget,
			stats.PenaltyGoals,
			stats.PenaltiesTaken,
			stats.FreeKickGoals,
			stats.FreeKickShots,

			// Positional Goals
			stats.GoalsFromInsideTheBox,
			stats.GoalsFromOutsideTheBox,
			stats.HeadedGoals,
			stats.LeftFootGoals,
			stats.RightFootGoals,

			// Chances
			stats.BigChances,
			stats.BigChancesCreated,
			stats.BigChancesMissed,

			// Shooting Details
			stats.ShotsFromInsideTheBox,
			stats.ShotsFromOutsideTheBox,
			stats.BlockedScoringAttempt,
			stats.HitWoodwork,

			// Dribbling
			stats.SuccessfulDribbles,
			stats.DribbleAttempts,

			// Set Pieces
			stats.Corners,
			stats.FreeKicks,
			stats.ThrowIns,
			stats.GoalKicks,

			// Fast Breaks
			stats.FastBreaks,
			stats.FastBreakGoals,
			stats.FastBreakShots,

			// Possession & Passing
			stats.AverageBallPossession,
			stats.TotalPasses,
			stats.AccuratePasses,
			stats.AccuratePassesPercentage,

			// Passing by Zone
			stats.TotalOwnHalfPasses,
			stats.AccurateOwnHalfPasses,
			stats.AccurateOwnHalfPassesPercentage,
			stats.TotalOppositionHalfPasses,
			stats.AccurateOppositionHalfPasses,
			stats.AccurateOppositionHalfPassesPercentage,

			// Long Balls & Crosses
			stats.TotalLongBalls,
			stats.AccurateLongBalls,
			stats.AccurateLongBallsPercentage,
			stats.TotalCrosses,
			stats.AccurateCrosses,
			stats.AccurateCrossesPercentage,

			// Defensive Stats
			stats.CleanSheets,
			stats.Tackles,
			stats.Interceptions,
			stats.Saves,
			stats.Clearances,
			stats.ClearancesOffLine,
			stats.LastManTackles,
			stats.BallRecovery,

			// Errors
			stats.ErrorsLeadingToGoal,
			stats.ErrorsLeadingToShot,

			// Penalties
			stats.PenaltiesCommited,
			stats.PenaltyGoalsConceded,

			// Duels
			stats.TotalDuels,
			stats.DuelsWon,
			stats.DuelsWonPercentage,
			stats.TotalGroundDuels,
			stats.GroundDuelsWon,
			stats.GroundDuelsWonPercentage,
			stats.TotalAerialDuels,
			stats.AerialDuelsWon,
			stats.AerialDuelsWonPercentage,

			// Discipline
			stats.PossessionLost,
			stats.Offsides,
			stats.Fouls,
			stats.YellowCards,
			stats.YellowRedCards,
			stats.RedCards,

			// Performance
			stats.AvgRating,
			stats.Matches,
			stats.AwardedMatches,

			// Stats Against (Defensive Perspective)
			stats.AccurateFinalThirdPassesAgainst,
			stats.AccurateOppositionHalfPassesAgainst,
			stats.AccurateOwnHalfPassesAgainst,
			stats.AccuratePassesAgainst,
			stats.BigChancesAgainst,
			stats.BigChancesCreatedAgainst,
			stats.BigChancesMissedAgainst,
			stats.ClearancesAgainst,
			stats.CornersAgainst,
			stats.CrossesSuccessfulAgainst,
			stats.CrossesTotalAgainst,
			stats.DribbleAttemptsTotalAgainst,
			stats.DribbleAttemptsWonAgainst,
			stats.ErrorsLeadingToGoalAgainst,
			stats.ErrorsLeadingToShotAgainst,
			stats.HitWoodworkAgainst,
			stats.InterceptionsAgainst,
			stats.KeyPassesAgainst,
			stats.LongBallsSuccessfulAgainst,
			stats.LongBallsTotalAgainst,
			stats.OffsidesAgainst,
			stats.RedCardsAgainst,
			stats.ShotsAgainst,
			stats.ShotsBlockedAgainst,
			stats.ShotsFromInsideTheBoxAgainst,
			stats.ShotsFromOutsideTheBoxAgainst,
			stats.ShotsOffTargetAgainst,
			stats.ShotsOnTargetAgainst,
			stats.BlockedScoringAttemptAgainst,
			stats.TacklesAgainst,
			stats.TotalFinalThirdPassesAgainst,
			stats.OppositionHalfPassesTotalAgainst,
			stats.OwnHalfPassesTotalAgainst,
			stats.TotalPassesAgainst,
			stats.YellowCardsAgainst,
		).
		Suffix(`
            ON CONFLICT (team_id, league_id, season_id) 
            DO UPDATE SET
                updated_at = CURRENT_TIMESTAMP,
                goals_scored = EXCLUDED.goals_scored,
                goals_conceded = EXCLUDED.goals_conceded,
                own_goals = EXCLUDED.own_goals,
                assists = EXCLUDED.assists,
                shots = EXCLUDED.shots,
                shots_on_target = EXCLUDED.shots_on_target,
                shots_off_target = EXCLUDED.shots_off_target,
                penalty_goals = EXCLUDED.penalty_goals,
                penalties_taken = EXCLUDED.penalties_taken,
                free_kick_goals = EXCLUDED.free_kick_goals,
                free_kick_shots = EXCLUDED.free_kick_shots,
                goals_from_inside_box = EXCLUDED.goals_from_inside_box,
                goals_from_outside_box = EXCLUDED.goals_from_outside_box,
                headed_goals = EXCLUDED.headed_goals,
                left_foot_goals = EXCLUDED.left_foot_goals,
                right_foot_goals = EXCLUDED.right_foot_goals,
                big_chances = EXCLUDED.big_chances,
                big_chances_created = EXCLUDED.big_chances_created,
                big_chances_missed = EXCLUDED.big_chances_missed,
                shots_from_inside_box = EXCLUDED.shots_from_inside_box,
                shots_from_outside_box = EXCLUDED.shots_from_outside_box,
                blocked_scoring_attempt = EXCLUDED.blocked_scoring_attempt,
                hit_woodwork = EXCLUDED.hit_woodwork,
                successful_dribbles = EXCLUDED.successful_dribbles,
                dribble_attempts = EXCLUDED.dribble_attempts,
                corners = EXCLUDED.corners,
                free_kicks = EXCLUDED.free_kicks,
                throw_ins = EXCLUDED.throw_ins,
                goal_kicks = EXCLUDED.goal_kicks,
                fast_breaks = EXCLUDED.fast_breaks,
                fast_break_goals = EXCLUDED.fast_break_goals,
                fast_break_shots = EXCLUDED.fast_break_shots,
                average_ball_possession = EXCLUDED.average_ball_possession,
                total_passes = EXCLUDED.total_passes,
                accurate_passes = EXCLUDED.accurate_passes,
                accurate_passes_percentage = EXCLUDED.accurate_passes_percentage,
                total_own_half_passes = EXCLUDED.total_own_half_passes,
                accurate_own_half_passes = EXCLUDED.accurate_own_half_passes,
                accurate_own_half_passes_percentage = EXCLUDED.accurate_own_half_passes_percentage,
                total_opposition_half_passes = EXCLUDED.total_opposition_half_passes,
                accurate_opposition_half_passes = EXCLUDED.accurate_opposition_half_passes,
                accurate_opposition_half_passes_percentage = EXCLUDED.accurate_opposition_half_passes_percentage,
                total_long_balls = EXCLUDED.total_long_balls,
                accurate_long_balls = EXCLUDED.accurate_long_balls,
                accurate_long_balls_percentage = EXCLUDED.accurate_long_balls_percentage,
                total_crosses = EXCLUDED.total_crosses,
                accurate_crosses = EXCLUDED.accurate_crosses,
                accurate_crosses_percentage = EXCLUDED.accurate_crosses_percentage,
                clean_sheets = EXCLUDED.clean_sheets,
                tackles = EXCLUDED.tackles,
                interceptions = EXCLUDED.interceptions,
                saves = EXCLUDED.saves,
                clearances = EXCLUDED.clearances,
                clearances_off_line = EXCLUDED.clearances_off_line,
                last_man_tackles = EXCLUDED.last_man_tackles,
                ball_recovery = EXCLUDED.ball_recovery,
                errors_leading_to_goal = EXCLUDED.errors_leading_to_goal,
                errors_leading_to_shot = EXCLUDED.errors_leading_to_shot,
                penalties_commited = EXCLUDED.penalties_commited,
                penalty_goals_conceded = EXCLUDED.penalty_goals_conceded,
                total_duels = EXCLUDED.total_duels,
                duels_won = EXCLUDED.duels_won,
                duels_won_percentage = EXCLUDED.duels_won_percentage,
                total_ground_duels = EXCLUDED.total_ground_duels,
                ground_duels_won = EXCLUDED.ground_duels_won,
                ground_duels_won_percentage = EXCLUDED.ground_duels_won_percentage,
                total_aerial_duels = EXCLUDED.total_aerial_duels,
                aerial_duels_won = EXCLUDED.aerial_duels_won,
                aerial_duels_won_percentage = EXCLUDED.aerial_duels_won_percentage,
                possession_lost = EXCLUDED.possession_lost,
                offsides = EXCLUDED.offsides,
                fouls = EXCLUDED.fouls,
                yellow_cards = EXCLUDED.yellow_cards,
                yellow_red_cards = EXCLUDED.yellow_red_cards,
                red_cards = EXCLUDED.red_cards,
                avg_rating = EXCLUDED.avg_rating,
                matches = EXCLUDED.matches,
                awarded_matches = EXCLUDED.awarded_matches,
                accurate_final_third_passes_against = EXCLUDED.accurate_final_third_passes_against,
                accurate_opposition_half_passes_against = EXCLUDED.accurate_opposition_half_passes_against,
                accurate_own_half_passes_against = EXCLUDED.accurate_own_half_passes_against,
                accurate_passes_against = EXCLUDED.accurate_passes_against,
                big_chances_against = EXCLUDED.big_chances_against,
                big_chances_created_against = EXCLUDED.big_chances_created_against,
                big_chances_missed_against = EXCLUDED.big_chances_missed_against,
                clearances_against = EXCLUDED.clearances_against,
                corners_against = EXCLUDED.corners_against,
                crosses_successful_against = EXCLUDED.crosses_successful_against,
                crosses_total_against = EXCLUDED.crosses_total_against,
                dribble_attempts_total_against = EXCLUDED.dribble_attempts_total_against,
                dribble_attempts_won_against = EXCLUDED.dribble_attempts_won_against,
                errors_leading_to_goal_against = EXCLUDED.errors_leading_to_goal_against,
                errors_leading_to_shot_against = EXCLUDED.errors_leading_to_shot_against,
                hit_woodwork_against = EXCLUDED.hit_woodwork_against,
                interceptions_against = EXCLUDED.interceptions_against,
                key_passes_against = EXCLUDED.key_passes_against,
                long_balls_successful_against = EXCLUDED.long_balls_successful_against,
                long_balls_total_against = EXCLUDED.long_balls_total_against,
                offsides_against = EXCLUDED.offsides_against,
                red_cards_against = EXCLUDED.red_cards_against,
                shots_against = EXCLUDED.shots_against,
                shots_blocked_against = EXCLUDED.shots_blocked_against,
                shots_from_inside_box_against = EXCLUDED.shots_from_inside_box_against,
                shots_from_outside_box_against = EXCLUDED.shots_from_outside_box_against,
                shots_off_target_against = EXCLUDED.shots_off_target_against,
                shots_on_target_against = EXCLUDED.shots_on_target_against,
                blocked_scoring_attempt_against = EXCLUDED.blocked_scoring_attempt_against,
                tackles_against = EXCLUDED.tackles_against,
                total_final_third_passes_against = EXCLUDED.total_final_third_passes_against,
                opposition_half_passes_total_against = EXCLUDED.opposition_half_passes_total_against,
                own_half_passes_total_against = EXCLUDED.own_half_passes_total_against,
                total_passes_against = EXCLUDED.total_passes_against,
                yellow_cards_against = EXCLUDED.yellow_cards_against
        `).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(sqlQuery, args...)
	return err
}

func (r *TeamRepo) InsertTeamMatchStats(matchStat sofascore_models.MatchStatsMessage) error {
	var tableName string
	columnMap := make(map[string]interface{})

	// Always add base columns
	columnMap["match_id"] = matchStat.MatchID
	columnMap["home_team_id"] = matchStat.HomeTeamID
	columnMap["away_team_id"] = matchStat.AwayTeamID
	columnMap["period"] = matchStat.MatchStatistics.Period

	stat := matchStat.MatchStatistics

	switch matchStat.GroupName {
	case "Match overview":
		tableName = "match_overview"
		r.Helper.MapOverviewStat(stat, columnMap)
	case "Shots":
		tableName = "match_shots"
		r.Helper.MapShotsStat(stat, columnMap)
	case "Attack":
		tableName = "match_attack"
		r.Helper.MapAttackStat(stat, columnMap)
	case "Passes":
		tableName = "match_passes"
		r.Helper.MapPassesStat(stat, columnMap)
	case "Duels":
		tableName = "match_duels"
		r.Helper.MapDuelsStat(stat, columnMap)
	case "Defending":
		tableName = "match_defending"
		r.Helper.MapDefendingStat(stat, columnMap)
	case "Goalkeeping":
		tableName = "match_goalkeeping"
		r.Helper.MapGoalkeepingStat(stat, columnMap)
	default:
		return fmt.Errorf("unknown group name: %s", matchStat.GroupName)
	}

	if len(columnMap) == 4 { // Only base columns, no stat columns added
		return nil
	}

	// Build query from map
	columns := make([]string, 0, len(columnMap))
	values := make([]interface{}, 0, len(columnMap))

	for col, val := range columnMap {
		columns = append(columns, col)
		values = append(values, val)
	}

	updateClauses := make([]string, 0, len(columns)-4)
	for _, col := range columns {
		if col != "match_id" && col != "home_team_id" && col != "away_team_id" && col != "period" {
			updateClauses = append(updateClauses, fmt.Sprintf("%s = EXCLUDED.%s", col, col))
		}
	}

	query := sq.Insert(tableName).
		Columns(columns...).
		Values(values...).
		Suffix("ON CONFLICT (match_id, period) DO UPDATE SET " + strings.Join(updateClauses, ", ")).
		PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(sqlQuery, args...)
	return err
}
