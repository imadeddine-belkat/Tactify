package sofascore_helper

import (
	sofascore "github.com/imadeddine-belkat/tactify-protos/go/sofascore/v1"
)

type Helper struct {
}

func (h *Helper) MapOverviewStat(stat *sofascore.MatchStatsMessage, m map[string]interface{}) {
	switch stat.Statistics.Key {
	case "ballPossession":
		m["ball_possession_home"] = stat.Statistics.HomeValue
		m["ball_possession_away"] = stat.Statistics.AwayValue
	case "expectedGoals":
		m["expected_goals_home"] = stat.Statistics.HomeValue
		m["expected_goals_away"] = stat.Statistics.AwayValue
	case "bigChanceCreated":
		m["big_chances_home"] = stat.Statistics.HomeValue
		m["big_chances_away"] = stat.Statistics.AwayValue
	case "totalShotsOnGoal":
		m["total_shots_home"] = stat.Statistics.HomeValue
		m["total_shots_away"] = stat.Statistics.AwayValue
	case "goalkeeperSaves":
		m["goalkeeper_saves_home"] = stat.Statistics.HomeValue
		m["goalkeeper_saves_away"] = stat.Statistics.AwayValue
	case "cornerKicks":
		m["corner_kicks_home"] = stat.Statistics.HomeValue
		m["corner_kicks_away"] = stat.Statistics.AwayValue
	case "fouls":
		m["fouls_home"] = stat.Statistics.HomeValue
		m["fouls_away"] = stat.Statistics.AwayValue
	case "passes":
		m["passes_home"] = stat.Statistics.HomeValue
		m["passes_away"] = stat.Statistics.AwayValue
	case "totalTackle":
		m["tackles_home"] = stat.Statistics.HomeValue
		m["tackles_away"] = stat.Statistics.AwayValue
	case "freeKicks":
		m["free_kicks_home"] = stat.Statistics.HomeValue
		m["free_kicks_away"] = stat.Statistics.AwayValue
	case "yellowCards":
		m["yellow_cards_home"] = stat.Statistics.HomeValue
		m["yellow_cards_away"] = stat.Statistics.AwayValue
	}
}

func (h *Helper) MapShotsStat(stat *sofascore.MatchStatsMessage, m map[string]interface{}) {
	switch stat.Statistics.Key {
	case "totalShotsOnGoal":
		m["total_shots_home"] = stat.Statistics.HomeValue
		m["total_shots_away"] = stat.Statistics.AwayValue
	case "shotsOnGoal":
		m["shots_on_target_home"] = stat.Statistics.HomeValue
		m["shots_on_target_away"] = stat.Statistics.AwayValue
	case "hitWoodwork":
		m["hit_woodwork_home"] = stat.Statistics.HomeValue
		m["hit_woodwork_away"] = stat.Statistics.AwayValue
	case "shotsOffGoal":
		m["shots_off_target_home"] = stat.Statistics.HomeValue
		m["shots_off_target_away"] = stat.Statistics.AwayValue
	case "blockedScoringAttempt":
		m["blocked_shots_home"] = stat.Statistics.HomeValue
		m["blocked_shots_away"] = stat.Statistics.AwayValue
	case "totalShotsInsideBox":
		m["shots_inside_box_home"] = stat.Statistics.HomeValue
		m["shots_inside_box_away"] = stat.Statistics.AwayValue
	case "totalShotsOutsideBox":
		m["shots_outside_box_home"] = stat.Statistics.HomeValue
		m["shots_outside_box_away"] = stat.Statistics.AwayValue
	}
}

func (h *Helper) MapAttackStat(stat *sofascore.MatchStatsMessage, m map[string]interface{}) {
	switch stat.Statistics.Key {
	case "bigChanceScored":
		m["big_chances_scored_home"] = stat.Statistics.HomeValue
		m["big_chances_scored_away"] = stat.Statistics.AwayValue
	case "bigChanceMissed":
		m["big_chances_missed_home"] = stat.Statistics.HomeValue
		m["big_chances_missed_away"] = stat.Statistics.AwayValue
	case "accurateThroughBall":
		m["through_balls_home"] = stat.Statistics.HomeValue
		m["through_balls_away"] = stat.Statistics.AwayValue
	case "touchesInOppBox":
		m["touches_in_penalty_area_home"] = stat.Statistics.HomeValue
		m["touches_in_penalty_area_away"] = stat.Statistics.AwayValue
	case "fouledFinalThird":
		m["fouled_in_final_third_home"] = stat.Statistics.HomeValue
		m["fouled_in_final_third_away"] = stat.Statistics.AwayValue
	case "offsides":
		m["offsides_home"] = stat.Statistics.HomeValue
		m["offsides_away"] = stat.Statistics.AwayValue
	}
}

func (h *Helper) MapPassesStat(stat *sofascore.MatchStatsMessage, m map[string]interface{}) {
	switch stat.Statistics.Key {
	case "accuratePasses":
		m["accurate_passes_home"] = stat.Statistics.HomeValue
		m["accurate_passes_away"] = stat.Statistics.AwayValue
	case "throwIns":
		m["throw_ins_home"] = stat.Statistics.HomeValue
		m["throw_ins_away"] = stat.Statistics.AwayValue
	case "finalThirdEntries":
		m["final_third_entries_home"] = stat.Statistics.HomeValue
		m["final_third_entries_away"] = stat.Statistics.AwayValue
	case "finalThirdPhaseStatistic":
		m["final_third_phase_home"] = stat.Statistics.HomeValue
		m["final_third_phase_away"] = stat.Statistics.AwayValue
		if stat.Statistics.HomeTotal != nil && stat.Statistics.AwayTotal != nil {
			m["final_third_phase_total_home"] = *stat.Statistics.HomeTotal
			m["final_third_phase_total_away"] = *stat.Statistics.AwayTotal
		}
	case "accurateLongBalls":
		m["long_balls_home"] = stat.Statistics.HomeValue
		m["long_balls_away"] = stat.Statistics.AwayValue
		if stat.Statistics.HomeTotal != nil && stat.Statistics.AwayTotal != nil {
			m["long_balls_total_home"] = *stat.Statistics.HomeTotal
			m["long_balls_total_away"] = *stat.Statistics.AwayTotal
		}
	case "accurateCross":
		m["crosses_home"] = stat.Statistics.HomeValue
		m["crosses_away"] = stat.Statistics.AwayValue
		if stat.Statistics.HomeTotal != nil && stat.Statistics.AwayTotal != nil {
			m["crosses_total_home"] = *stat.Statistics.HomeTotal
			m["crosses_total_away"] = *stat.Statistics.AwayTotal
		}
	}
}

func (h *Helper) MapDuelsStat(stat *sofascore.MatchStatsMessage, m map[string]interface{}) {
	switch stat.Statistics.Key {
	case "duelWonPercent":
		m["duels_won_percent_home"] = stat.Statistics.HomeValue
		m["duels_won_percent_away"] = stat.Statistics.AwayValue
	case "dispossessed":
		m["dispossessed_home"] = stat.Statistics.HomeValue
		m["dispossessed_away"] = stat.Statistics.AwayValue
	case "groundDuelsPercentage":
		m["ground_duels_home"] = stat.Statistics.HomeValue
		m["ground_duels_away"] = stat.Statistics.AwayValue
		if stat.Statistics.HomeTotal != nil && stat.Statistics.AwayTotal != nil {
			m["ground_duels_total_home"] = *stat.Statistics.HomeTotal
			m["ground_duels_total_away"] = *stat.Statistics.AwayTotal
		}
	case "aerialDuelsPercentage":
		m["aerial_duels_home"] = stat.Statistics.HomeValue
		m["aerial_duels_away"] = stat.Statistics.AwayValue
		if stat.Statistics.HomeTotal != nil && stat.Statistics.AwayTotal != nil {
			m["aerial_duels_total_home"] = *stat.Statistics.HomeTotal
			m["aerial_duels_total_away"] = *stat.Statistics.AwayTotal
		}
	case "dribblesPercentage":
		m["dribbles_home"] = stat.Statistics.HomeValue
		m["dribbles_away"] = stat.Statistics.AwayValue
		if stat.Statistics.HomeTotal != nil && stat.Statistics.AwayTotal != nil {
			m["dribbles_total_home"] = *stat.Statistics.HomeTotal
			m["dribbles_total_away"] = *stat.Statistics.AwayTotal
		}
	}
}

func (h *Helper) MapDefendingStat(stat *sofascore.MatchStatsMessage, m map[string]interface{}) {
	switch stat.Statistics.Key {
	case "wonTacklePercent":
		m["tackles_won_home"] = stat.Statistics.HomeValue
		m["tackles_won_away"] = stat.Statistics.AwayValue
		if stat.Statistics.HomeTotal != nil && stat.Statistics.AwayTotal != nil {
			m["tackles_won_total_home"] = *stat.Statistics.HomeTotal
			m["tackles_won_total_away"] = *stat.Statistics.AwayTotal
		}
	case "totalTackle":
		m["total_tackles_home"] = stat.Statistics.HomeValue
		m["total_tackles_away"] = stat.Statistics.AwayValue
	case "interceptionWon":
		m["interceptions_home"] = stat.Statistics.HomeValue
		m["interceptions_away"] = stat.Statistics.AwayValue
	case "ballRecovery":
		m["recoveries_home"] = stat.Statistics.HomeValue
		m["recoveries_away"] = stat.Statistics.AwayValue
	case "totalClearance":
		m["clearances_home"] = stat.Statistics.HomeValue
		m["clearances_away"] = stat.Statistics.AwayValue
	case "errorsLeadToShot":
		m["errors_lead_to_shot_home"] = stat.Statistics.HomeValue
		m["errors_lead_to_shot_away"] = stat.Statistics.AwayValue
	}
}

func (h *Helper) MapGoalkeepingStat(stat *sofascore.MatchStatsMessage, m map[string]interface{}) {
	switch stat.Statistics.Key {
	case "goalkeeperSaves":
		m["total_saves_home"] = stat.Statistics.HomeValue
		m["total_saves_away"] = stat.Statistics.AwayValue
	case "goalsPrevented":
		m["goals_prevented_home"] = stat.Statistics.HomeValue
		m["goals_prevented_away"] = stat.Statistics.AwayValue
	case "diveSaves":
		m["big_saves_home"] = stat.Statistics.HomeValue
		m["big_saves_away"] = stat.Statistics.AwayValue
	case "highClaims":
		m["high_claims_home"] = stat.Statistics.HomeValue
		m["high_claims_away"] = stat.Statistics.AwayValue
	case "punches":
		m["punches_home"] = stat.Statistics.HomeValue
		m["punches_away"] = stat.Statistics.AwayValue
	case "goalKicks":
		m["goal_kicks_home"] = stat.Statistics.HomeValue
		m["goal_kicks_away"] = stat.Statistics.AwayValue
	}
}
