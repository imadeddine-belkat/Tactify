package sofascore_helper

import "github.com/imadeddine-belkat/shared/sofascore_models"

type Helper struct {
}

func (h *Helper) MapOverviewStat(stat sofascore_models.StatsMessage, m map[string]interface{}) {
	switch stat.Key {
	case "ballPossession":
		m["ball_possession_home"] = stat.HomeValue
		m["ball_possession_away"] = stat.AwayValue
	case "expectedGoals":
		m["expected_goals_home"] = stat.HomeValue
		m["expected_goals_away"] = stat.AwayValue
	case "bigChanceCreated":
		m["big_chances_home"] = stat.HomeValue
		m["big_chances_away"] = stat.AwayValue
	case "totalShotsOnGoal":
		m["total_shots_home"] = stat.HomeValue
		m["total_shots_away"] = stat.AwayValue
	case "goalkeeperSaves":
		m["goalkeeper_saves_home"] = stat.HomeValue
		m["goalkeeper_saves_away"] = stat.AwayValue
	case "cornerKicks":
		m["corner_kicks_home"] = stat.HomeValue
		m["corner_kicks_away"] = stat.AwayValue
	case "fouls":
		m["fouls_home"] = stat.HomeValue
		m["fouls_away"] = stat.AwayValue
	case "passes":
		m["passes_home"] = stat.HomeValue
		m["passes_away"] = stat.AwayValue
	case "totalTackle":
		m["tackles_home"] = stat.HomeValue
		m["tackles_away"] = stat.AwayValue
	case "freeKicks":
		m["free_kicks_home"] = stat.HomeValue
		m["free_kicks_away"] = stat.AwayValue
	case "yellowCards":
		m["yellow_cards_home"] = stat.HomeValue
		m["yellow_cards_away"] = stat.AwayValue
	}
}

func (h *Helper) MapShotsStat(stat sofascore_models.StatsMessage, m map[string]interface{}) {
	switch stat.Key {
	case "totalShotsOnGoal":
		m["total_shots_home"] = stat.HomeValue
		m["total_shots_away"] = stat.AwayValue
	case "shotsOnGoal":
		m["shots_on_target_home"] = stat.HomeValue
		m["shots_on_target_away"] = stat.AwayValue
	case "hitWoodwork":
		m["hit_woodwork_home"] = stat.HomeValue
		m["hit_woodwork_away"] = stat.AwayValue
	case "shotsOffGoal":
		m["shots_off_target_home"] = stat.HomeValue
		m["shots_off_target_away"] = stat.AwayValue
	case "blockedScoringAttempt":
		m["blocked_shots_home"] = stat.HomeValue
		m["blocked_shots_away"] = stat.AwayValue
	case "totalShotsInsideBox":
		m["shots_inside_box_home"] = stat.HomeValue
		m["shots_inside_box_away"] = stat.AwayValue
	case "totalShotsOutsideBox":
		m["shots_outside_box_home"] = stat.HomeValue
		m["shots_outside_box_away"] = stat.AwayValue
	}
}

func (h *Helper) MapAttackStat(stat sofascore_models.StatsMessage, m map[string]interface{}) {
	switch stat.Key {
	case "bigChanceScored":
		m["big_chances_scored_home"] = stat.HomeValue
		m["big_chances_scored_away"] = stat.AwayValue
	case "bigChanceMissed":
		m["big_chances_missed_home"] = stat.HomeValue
		m["big_chances_missed_away"] = stat.AwayValue
	case "accurateThroughBall":
		m["through_balls_home"] = stat.HomeValue
		m["through_balls_away"] = stat.AwayValue
	case "touchesInOppBox":
		m["touches_in_penalty_area_home"] = stat.HomeValue
		m["touches_in_penalty_area_away"] = stat.AwayValue
	case "fouledFinalThird":
		m["fouled_in_final_third_home"] = stat.HomeValue
		m["fouled_in_final_third_away"] = stat.AwayValue
	case "offsides":
		m["offsides_home"] = stat.HomeValue
		m["offsides_away"] = stat.AwayValue
	}
}

func (h *Helper) MapPassesStat(stat sofascore_models.StatsMessage, m map[string]interface{}) {
	switch stat.Key {
	case "accuratePasses":
		m["accurate_passes_home"] = stat.HomeValue
		m["accurate_passes_away"] = stat.AwayValue
	case "throwIns":
		m["throw_ins_home"] = stat.HomeValue
		m["throw_ins_away"] = stat.AwayValue
	case "finalThirdEntries":
		m["final_third_entries_home"] = stat.HomeValue
		m["final_third_entries_away"] = stat.AwayValue
	case "finalThirdPhaseStatistic":
		m["final_third_phase_home"] = stat.HomeValue
		m["final_third_phase_away"] = stat.AwayValue
		if stat.HomeTotal != nil && stat.AwayTotal != nil {
			m["final_third_phase_total_home"] = *stat.HomeTotal
			m["final_third_phase_total_away"] = *stat.AwayTotal
		}
	case "accurateLongBalls":
		m["long_balls_home"] = stat.HomeValue
		m["long_balls_away"] = stat.AwayValue
		if stat.HomeTotal != nil && stat.AwayTotal != nil {
			m["long_balls_total_home"] = *stat.HomeTotal
			m["long_balls_total_away"] = *stat.AwayTotal
		}
	case "accurateCross":
		m["crosses_home"] = stat.HomeValue
		m["crosses_away"] = stat.AwayValue
		if stat.HomeTotal != nil && stat.AwayTotal != nil {
			m["crosses_total_home"] = *stat.HomeTotal
			m["crosses_total_away"] = *stat.AwayTotal
		}
	}
}

func (h *Helper) MapDuelsStat(stat sofascore_models.StatsMessage, m map[string]interface{}) {
	switch stat.Key {
	case "duelWonPercent":
		m["duels_won_percent_home"] = stat.HomeValue
		m["duels_won_percent_away"] = stat.AwayValue
	case "dispossessed":
		m["dispossessed_home"] = stat.HomeValue
		m["dispossessed_away"] = stat.AwayValue
	case "groundDuelsPercentage":
		m["ground_duels_home"] = stat.HomeValue
		m["ground_duels_away"] = stat.AwayValue
		if stat.HomeTotal != nil && stat.AwayTotal != nil {
			m["ground_duels_total_home"] = *stat.HomeTotal
			m["ground_duels_total_away"] = *stat.AwayTotal
		}
	case "aerialDuelsPercentage":
		m["aerial_duels_home"] = stat.HomeValue
		m["aerial_duels_away"] = stat.AwayValue
		if stat.HomeTotal != nil && stat.AwayTotal != nil {
			m["aerial_duels_total_home"] = *stat.HomeTotal
			m["aerial_duels_total_away"] = *stat.AwayTotal
		}
	case "dribblesPercentage":
		m["dribbles_home"] = stat.HomeValue
		m["dribbles_away"] = stat.AwayValue
		if stat.HomeTotal != nil && stat.AwayTotal != nil {
			m["dribbles_total_home"] = *stat.HomeTotal
			m["dribbles_total_away"] = *stat.AwayTotal
		}
	}
}

func (h *Helper) MapDefendingStat(stat sofascore_models.StatsMessage, m map[string]interface{}) {
	switch stat.Key {
	case "wonTacklePercent":
		m["tackles_won_home"] = stat.HomeValue
		m["tackles_won_away"] = stat.AwayValue
		if stat.HomeTotal != nil && stat.AwayTotal != nil {
			m["tackles_won_total_home"] = *stat.HomeTotal
			m["tackles_won_total_away"] = *stat.AwayTotal
		}
	case "totalTackle":
		m["total_tackles_home"] = stat.HomeValue
		m["total_tackles_away"] = stat.AwayValue
	case "interceptionWon":
		m["interceptions_home"] = stat.HomeValue
		m["interceptions_away"] = stat.AwayValue
	case "ballRecovery":
		m["recoveries_home"] = stat.HomeValue
		m["recoveries_away"] = stat.AwayValue
	case "totalClearance":
		m["clearances_home"] = stat.HomeValue
		m["clearances_away"] = stat.AwayValue
	case "errorsLeadToShot":
		m["errors_lead_to_shot_home"] = stat.HomeValue
		m["errors_lead_to_shot_away"] = stat.AwayValue
	}
}

func (h *Helper) MapGoalkeepingStat(stat sofascore_models.StatsMessage, m map[string]interface{}) {
	switch stat.Key {
	case "goalkeeperSaves":
		m["total_saves_home"] = stat.HomeValue
		m["total_saves_away"] = stat.AwayValue
	case "goalsPrevented":
		m["goals_prevented_home"] = stat.HomeValue
		m["goals_prevented_away"] = stat.AwayValue
	case "diveSaves":
		m["big_saves_home"] = stat.HomeValue
		m["big_saves_away"] = stat.AwayValue
	case "highClaims":
		m["high_claims_home"] = stat.HomeValue
		m["high_claims_away"] = stat.AwayValue
	case "punches":
		m["punches_home"] = stat.HomeValue
		m["punches_away"] = stat.AwayValue
	case "goalKicks":
		m["goal_kicks_home"] = stat.HomeValue
		m["goal_kicks_away"] = stat.AwayValue
	}
}
