package fpl_repositories

import (
	"database/sql"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/imadeddine-belkat/shared/fpl_models"
)

type PlayerRepo struct {
	db                     *sql.DB
	PlayerBootstrapModel   *fpl_models.PlayerBootstrap
	PlayerHistoryModel     *fpl_models.PlayerHistory
	PlayerPastHistoryModel *fpl_models.PlayerPastHistory
}

func NewPlayerRepo(db *sql.DB, playerBootstrapModel *fpl_models.PlayerBootstrap, PlayerHistoryModel *fpl_models.PlayerHistory, PlayerPastHistoryModel *fpl_models.PlayerPastHistory) *PlayerRepo {
	return &PlayerRepo{
		db:                     db,
		PlayerBootstrapModel:   playerBootstrapModel,
		PlayerHistoryModel:     PlayerHistoryModel,
		PlayerPastHistoryModel: PlayerPastHistoryModel,
	}
}

func (r *PlayerRepo) CountPlayers() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM players"
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting players: %w", err)
	}
	return count, nil
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// InsertPlayerBootstrapComplete handles all player bootstrap data at once
func (r *PlayerRepo) InsertPlayerBootstrapComplete(players []fpl_models.PlayerBootstrapMessage) error {
	log.Printf("Attempting to insert %d players...", len(players))

	if err := r.InsertPlayers(players); err != nil {
		log.Printf("❌ Error inserting players table: %v", err)
		return fmt.Errorf("inserting players: %w", err)
	}
	log.Printf("Inserted %d players into players table", len(players))

	if err := r.InsertPlayerCosts(players); err != nil {
		log.Printf("❌ Error inserting player_costs: %v", err)
		return fmt.Errorf("inserting player costs: %w", err)
	}
	log.Printf("Inserted %d records into player_costs", len(players))

	if err := r.InsertPlayerSeasonStats(players); err != nil {
		log.Printf("❌ Error inserting player_season_stats: %v", err)
		return fmt.Errorf("inserting player season stats: %w", err)
	}
	log.Printf("Inserted %d records into player_season_stats", len(players))

	if err := r.InsertPlayerICTStats(players); err != nil {
		log.Printf("❌ Error inserting player_ict_stats: %v", err)
		return fmt.Errorf("inserting player ICT stats: %w", err)
	}
	log.Printf("Inserted %d records into player_ict_stats", len(players))

	if err := r.InsertPlayerExpectedStats(players); err != nil {
		log.Printf("❌ Error inserting player_expected_stats: %v", err)
		return fmt.Errorf("inserting player expected stats: %w", err)
	}
	log.Printf("Inserted %d records into player_expected_stats", len(players))

	if err := r.InsertPlayerRankings(players); err != nil {
		log.Printf("❌ Error inserting player_rankings: %v", err)
		return fmt.Errorf("inserting player rankings: %w", err)
	}
	log.Printf("Inserted %d records into player_rankings", len(players))

	return nil
}

// InsertPlayers inserts/updates main player information
func (r *PlayerRepo) InsertPlayers(players []fpl_models.PlayerBootstrapMessage) error {
	if len(players) == 0 {
		return nil
	}

	query := sq.Insert("players").Columns(
		"season_id", "player_id", "player_code", "first_name", "second_name", "web_name",
		"team_id", "team_code", "element_type_id", "status", "photo",
		"squad_number", "birth_date", "team_join_date", "region", "opta_code",
		"can_transact", "can_select", "in_dreamteam", "dreamteam_count", "special", "removed", "unavailable",
	).Suffix("ON CONFLICT (player_id, season_id) DO UPDATE SET " +
		"season_id = EXCLUDED.season_id, " +
		"player_code = EXCLUDED.player_code, " +
		"first_name = EXCLUDED.first_name, " +
		"second_name = EXCLUDED.second_name, " +
		"web_name = EXCLUDED.web_name, " +
		"team_id = EXCLUDED.team_id, " +
		"team_code = EXCLUDED.team_code, " +
		"element_type_id = EXCLUDED.element_type_id, " +
		"status = EXCLUDED.status, " +
		"photo = EXCLUDED.photo, " +
		"squad_number = EXCLUDED.squad_number, " +
		"birth_date = EXCLUDED.birth_date, " +
		"team_join_date = EXCLUDED.team_join_date, " +
		"region = EXCLUDED.region, " +
		"opta_code = EXCLUDED.opta_code, " +
		"can_transact = EXCLUDED.can_transact, " +
		"can_select = EXCLUDED.can_select, " +
		"in_dreamteam = EXCLUDED.in_dreamteam, " +
		"dreamteam_count = EXCLUDED.dreamteam_count, " +
		"special = EXCLUDED.special, " +
		"removed = EXCLUDED.removed, " +
		"unavailable = EXCLUDED.unavailable, " +
		"updated_at = CURRENT_TIMESTAMP").
		PlaceholderFormat(sq.Dollar)

	successCount := 0
	for _, p := range players {
		query = query.Values(
			p.SeasonID, p.Player.ID, p.Player.Code, p.Player.FirstName, p.Player.SecondName, p.Player.WebName,
			p.Player.TeamID, p.Player.TeamCode, p.Player.ElementType, p.Player.Status, p.Player.Photo,
			p.Player.SquadNumber, nullIfEmpty(p.Player.BirthDate), nullIfEmpty(p.Player.TeamJoinDate), p.Player.Region, p.Player.OptaCode,
			p.Player.CanTransact, p.Player.CanSelect, p.Player.InDreamteam, p.Player.DreamteamCount, p.Player.Special, p.Player.Removed, p.Player.Unavailable,
		)
		successCount++
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building players insert query: %w", err)
	}

	result, err := r.db.Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("❌ SQL Error inserting players: %v", err)
		log.Printf("   Query length: %d characters", len(sqlQuery))
		log.Printf("   Args count: %d", len(args))
		return fmt.Errorf("executing players insert: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Players insert completed: %d rows affected", rowsAffected)

	return nil
}

// InsertPlayerCosts inserts/updates player cost information
func (r *PlayerRepo) InsertPlayerCosts(players []fpl_models.PlayerBootstrapMessage) error {
	query := sq.Insert("player_costs").Columns(
		"player_id", "now_cost", "cost_change_event", "cost_change_event_fall",
		"cost_change_start", "cost_change_start_fall",
	).Suffix("ON CONFLICT (player_id) DO UPDATE SET " +
		"now_cost = EXCLUDED.now_cost, " +
		"cost_change_event = EXCLUDED.cost_change_event, " +
		"cost_change_event_fall = EXCLUDED.cost_change_event_fall, " +
		"cost_change_start = EXCLUDED.cost_change_start, " +
		"cost_change_start_fall = EXCLUDED.cost_change_start_fall, " +
		"updated_at = CURRENT_TIMESTAMP").
		PlaceholderFormat(sq.Dollar)

	for _, p := range players {
		query = query.Values(
			p.Player.ID, p.Player.NowCost, p.Player.CostChangeEvent, p.Player.CostChangeEventFall,
			p.Player.CostChangeStart, p.Player.CostChangeStartFall,
		)
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building player_costs insert query: %w", err)
	}

	_, err = r.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("executing player_costs insert: %w", err)
	}
	return nil
}

// InsertPlayerSeasonStats inserts/updates player season statistics
func (r *PlayerRepo) InsertPlayerSeasonStats(players []fpl_models.PlayerBootstrapMessage) error {
	query := sq.Insert("player_season_stats").Columns(
		"player_id", "dreamteam_count", "total_points", "event_points", "points_per_game",
		"form", "selected_by_percent", "value_form", "value_season",
		"minutes", "goals_scored", "assists", "clean_sheets", "goals_conceded",
		"own_goals", "penalties_saved", "penalties_missed", "yellow_cards", "red_cards",
		"saves", "bonus", "bps", "starts", "clearances_blocks_interceptions",
		"recoveries", "tackles", "defensive_contribution",
	).Suffix("ON CONFLICT (player_id) DO UPDATE SET " +
		"dreamteam_count = EXCLUDED.dreamteam_count, " +
		"total_points = EXCLUDED.total_points, " +
		"event_points = EXCLUDED.event_points, " +
		"points_per_game = EXCLUDED.points_per_game, " +
		"form = EXCLUDED.form, " +
		"selected_by_percent = EXCLUDED.selected_by_percent, " +
		"value_form = EXCLUDED.value_form, " +
		"value_season = EXCLUDED.value_season, " +
		"minutes = EXCLUDED.minutes, " +
		"goals_scored = EXCLUDED.goals_scored, " +
		"assists = EXCLUDED.assists, " +
		"clean_sheets = EXCLUDED.clean_sheets, " +
		"goals_conceded = EXCLUDED.goals_conceded, " +
		"own_goals = EXCLUDED.own_goals, " +
		"penalties_saved = EXCLUDED.penalties_saved, " +
		"penalties_missed = EXCLUDED.penalties_missed, " +
		"yellow_cards = EXCLUDED.yellow_cards, " +
		"red_cards = EXCLUDED.red_cards, " +
		"saves = EXCLUDED.saves, " +
		"bonus = EXCLUDED.bonus, " +
		"bps = EXCLUDED.bps, " +
		"starts = EXCLUDED.starts, " +
		"clearances_blocks_interceptions = EXCLUDED.clearances_blocks_interceptions, " +
		"recoveries = EXCLUDED.recoveries, " +
		"tackles = EXCLUDED.tackles, " +
		"defensive_contribution = EXCLUDED.defensive_contribution, " +
		"updated_at = CURRENT_TIMESTAMP").
		PlaceholderFormat(sq.Dollar)

	for _, p := range players {
		query = query.Values(
			p.Player.ID, p.Player.DreamteamCount, p.Player.TotalPoints, p.Player.EventPoints, p.Player.PointsPerGame,
			p.Player.Form, p.Player.SelectedByPercent, p.Player.ValueForm, p.Player.ValueSeason,
			p.Player.Minutes, p.Player.GoalsScored, p.Player.Assists, p.Player.CleanSheets, p.Player.GoalsConceded,
			p.Player.OwnGoals, p.Player.PenaltiesSaved, p.Player.PenaltiesMissed, p.Player.YellowCards, p.Player.RedCards,
			p.Player.Saves, p.Player.Bonus, p.Player.BPS, p.Player.Starts, p.Player.ClearancesBlocksInterceptions,
			p.Player.Recoveries, p.Player.Tackles, p.Player.DefensiveContribution,
		)
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building player_season_stats insert query: %w", err)
	}

	_, err = r.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("executing player_season_stats insert: %w", err)
	}
	return nil
}

// InsertPlayerICTStats inserts/updates player ICT statistics
func (r *PlayerRepo) InsertPlayerICTStats(players []fpl_models.PlayerBootstrapMessage) error {
	query := sq.Insert("player_ict_stats").Columns(
		"player_id", "influence", "creativity", "threat", "ict_index",
		"influence_rank", "influence_rank_type", "creativity_rank", "creativity_rank_type",
		"threat_rank", "threat_rank_type", "ict_index_rank", "ict_index_rank_type",
	).Suffix("ON CONFLICT (player_id) DO UPDATE SET " +
		"influence = EXCLUDED.influence, " +
		"creativity = EXCLUDED.creativity, " +
		"threat = EXCLUDED.threat, " +
		"ict_index = EXCLUDED.ict_index, " +
		"influence_rank = EXCLUDED.influence_rank, " +
		"influence_rank_type = EXCLUDED.influence_rank_type, " +
		"creativity_rank = EXCLUDED.creativity_rank, " +
		"creativity_rank_type = EXCLUDED.creativity_rank_type, " +
		"threat_rank = EXCLUDED.threat_rank, " +
		"threat_rank_type = EXCLUDED.threat_rank_type, " +
		"ict_index_rank = EXCLUDED.ict_index_rank, " +
		"ict_index_rank_type = EXCLUDED.ict_index_rank_type, " +
		"updated_at = CURRENT_TIMESTAMP").
		PlaceholderFormat(sq.Dollar)

	for _, p := range players {
		query = query.Values(
			p.Player.ID, p.Player.Influence, p.Player.Creativity, p.Player.Threat, p.Player.ICTIndex,
			p.Player.InfluenceRank, p.Player.InfluenceRankType, p.Player.CreativityRank, p.Player.CreativityRankType,
			p.Player.ThreatRank, p.Player.ThreatRankType, p.Player.ICTIndexRank, p.Player.ICTIndexRankType,
		)
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building player_ict_stats insert query: %w", err)
	}

	_, err = r.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("executing player_ict_stats insert: %w", err)
	}
	return nil
}

// InsertPlayerExpectedStats inserts/updates player expected statistics
func (r *PlayerRepo) InsertPlayerExpectedStats(players []fpl_models.PlayerBootstrapMessage) error {
	query := sq.Insert("player_expected_stats").Columns(
		"player_id", "expected_goals", "expected_assists", "expected_goal_involvements", "expected_goals_conceded",
		"expected_goals_per_90", "expected_assists_per_90", "expected_goal_involvements_per_90",
		"expected_goals_conceded_per_90", "saves_per_90", "goals_conceded_per_90",
		"starts_per_90", "clean_sheets_per_90", "defensive_contribution_per_90",
	).Suffix("ON CONFLICT (player_id) DO UPDATE SET " +
		"expected_goals = EXCLUDED.expected_goals, " +
		"expected_assists = EXCLUDED.expected_assists, " +
		"expected_goal_involvements = EXCLUDED.expected_goal_involvements, " +
		"expected_goals_conceded = EXCLUDED.expected_goals_conceded, " +
		"expected_goals_per_90 = EXCLUDED.expected_goals_per_90, " +
		"expected_assists_per_90 = EXCLUDED.expected_assists_per_90, " +
		"expected_goal_involvements_per_90 = EXCLUDED.expected_goal_involvements_per_90, " +
		"expected_goals_conceded_per_90 = EXCLUDED.expected_goals_conceded_per_90, " +
		"saves_per_90 = EXCLUDED.saves_per_90, " +
		"goals_conceded_per_90 = EXCLUDED.goals_conceded_per_90, " +
		"starts_per_90 = EXCLUDED.starts_per_90, " +
		"clean_sheets_per_90 = EXCLUDED.clean_sheets_per_90, " +
		"defensive_contribution_per_90 = EXCLUDED.defensive_contribution_per_90, " +
		"updated_at = CURRENT_TIMESTAMP").
		PlaceholderFormat(sq.Dollar)

	for _, p := range players {
		query = query.Values(
			p.Player.ID, p.Player.ExpectedGoals, p.Player.ExpectedAssists, p.Player.ExpectedGoalInvolvements, p.Player.ExpectedGoalsConceded,
			p.Player.ExpectedGoalsPer90, p.Player.ExpectedAssistsPer90, p.Player.ExpectedGoalInvolvementsPer90,
			p.Player.ExpectedGoalsConcededPer90, p.Player.SavesPer90, p.Player.GoalsConcededPer90,
			p.Player.StartsPer90, p.Player.CleanSheetsPer90, p.Player.DefensiveContributionPer90,
		)
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building player_expected_stats insert query: %w", err)
	}

	_, err = r.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("executing player_expected_stats insert: %w", err)
	}
	return nil
}

// InsertPlayerRankings inserts/updates player rankings
func (r *PlayerRepo) InsertPlayerRankings(players []fpl_models.PlayerBootstrapMessage) error {
	query := sq.Insert("player_rankings").Columns(
		"player_id", "now_cost_rank", "now_cost_rank_type",
		"form_rank", "form_rank_type", "points_per_game_rank", "points_per_game_rank_type",
		"selected_rank", "selected_rank_type",
	).Suffix("ON CONFLICT (player_id) DO UPDATE SET " +
		"now_cost_rank = EXCLUDED.now_cost_rank, " +
		"now_cost_rank_type = EXCLUDED.now_cost_rank_type, " +
		"form_rank = EXCLUDED.form_rank, " +
		"form_rank_type = EXCLUDED.form_rank_type, " +
		"points_per_game_rank = EXCLUDED.points_per_game_rank, " +
		"points_per_game_rank_type = EXCLUDED.points_per_game_rank_type, " +
		"selected_rank = EXCLUDED.selected_rank, " +
		"selected_rank_type = EXCLUDED.selected_rank_type, " +
		"updated_at = CURRENT_TIMESTAMP").
		PlaceholderFormat(sq.Dollar)

	for _, p := range players {
		query = query.Values(
			p.Player.ID, p.Player.NowCostRank, p.Player.NowCostRankType,
			p.Player.FormRank, p.Player.FormRankType, p.Player.PointsPerGameRank, p.Player.PointsPerGameRankType,
			p.Player.SelectedRank, p.Player.SelectedRankType,
		)
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building player_rankings insert query: %w", err)
	}

	_, err = r.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("executing player_rankings insert: %w", err)
	}
	return nil
}

// InsertPlayerGameweekStats inserts/updates player gameweek performance stats
func (r *PlayerRepo) InsertPlayerGameweekStats(history fpl_models.PlayerHistoryMessage) error {
	playerHistory := history.History

	if len(playerHistory) == 0 {
		return nil
	}

	query := sq.Insert("player_gameweek_stats").Columns(
		"player_id", "fixture_id", "season_id", "event", "opponent_team_id", "kickoff_time", "was_home",
		"team_h_score", "team_a_score", "minutes", "goals_scored", "assists", "clean_sheets",
		"goals_conceded", "own_goals", "penalties_saved", "penalties_missed", "yellow_cards",
		"red_cards", "saves", "bonus", "bps", "starts", "clearances_blocks_interceptions",
		"recoveries", "tackles", "defensive_contribution", "influence", "creativity", "threat",
		"ict_index", "expected_goals", "expected_assists", "expected_goal_involvements",
		"expected_goals_conceded", "total_points", "value", "transfers_balance", "selected",
		"transfers_in", "transfers_out", "modified",
	).Suffix("ON CONFLICT (player_id, fixture_id, season_id) DO UPDATE SET " +
		"event = EXCLUDED.event, " +
		"opponent_team_id = EXCLUDED.opponent_team_id, " +
		"kickoff_time = EXCLUDED.kickoff_time, " +
		"was_home = EXCLUDED.was_home, " +
		"team_h_score = EXCLUDED.team_h_score, " +
		"team_a_score = EXCLUDED.team_a_score, " +
		"minutes = EXCLUDED.minutes, " +
		"goals_scored = EXCLUDED.goals_scored, " +
		"assists = EXCLUDED.assists, " +
		"clean_sheets = EXCLUDED.clean_sheets, " +
		"goals_conceded = EXCLUDED.goals_conceded, " +
		"own_goals = EXCLUDED.own_goals, " +
		"penalties_saved = EXCLUDED.penalties_saved, " +
		"penalties_missed = EXCLUDED.penalties_missed, " +
		"yellow_cards = EXCLUDED.yellow_cards, " +
		"red_cards = EXCLUDED.red_cards, " +
		"saves = EXCLUDED.saves, " +
		"bonus = EXCLUDED.bonus, " +
		"bps = EXCLUDED.bps, " +
		"starts = EXCLUDED.starts, " +
		"clearances_blocks_interceptions = EXCLUDED.clearances_blocks_interceptions, " +
		"recoveries = EXCLUDED.recoveries, " +
		"tackles = EXCLUDED.tackles, " +
		"defensive_contribution = EXCLUDED.defensive_contribution, " +
		"influence = EXCLUDED.influence, " +
		"creativity = EXCLUDED.creativity, " +
		"threat = EXCLUDED.threat, " +
		"ict_index = EXCLUDED.ict_index, " +
		"expected_goals = EXCLUDED.expected_goals, " +
		"expected_assists = EXCLUDED.expected_assists, " +
		"expected_goal_involvements = EXCLUDED.expected_goal_involvements, " +
		"expected_goals_conceded = EXCLUDED.expected_goals_conceded, " +
		"total_points = EXCLUDED.total_points, " +
		"value = EXCLUDED.value, " +
		"transfers_balance = EXCLUDED.transfers_balance, " +
		"selected = EXCLUDED.selected, " +
		"transfers_in = EXCLUDED.transfers_in, " +
		"transfers_out = EXCLUDED.transfers_out, " +
		"modified = EXCLUDED.modified," +
		"updated_at = CURRENT_TIMESTAMP").
		PlaceholderFormat(sq.Dollar)

	for _, h := range playerHistory {
		query = query.Values(
			h.PlayerID, h.FixtureID, history.SeasonID, h.Round, h.OpponentTeam, h.KickoffTime,
			h.WasHome, h.TeamHScore, h.TeamAScore, h.Minutes, h.GoalsScored,
			h.Assists, h.CleanSheets, h.GoalsConceded, h.OwnGoals,
			h.PenaltiesSaved, h.PenaltiesMissed, h.YellowCards, h.RedCards,
			h.Saves, h.Bonus, h.BPS, h.Starts, h.ClearancesBlocksInterceptions,
			h.Recoveries, h.Tackles, h.DefensiveContribution, h.Influence,
			h.Creativity, h.Threat, h.ICTIndex, h.ExpectedGoals,
			h.ExpectedAssists, h.ExpectedGoalInvolvements, h.ExpectedGoalsConceded,
			h.TotalPoints, h.Value, h.TransfersBalance, h.Selected,
			h.TransfersIn, h.TransfersOut, h.Modified,
		)

	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building player_gameweek_stats insert query: %w", err)
	}

	_, err = r.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("executing player_gameweek_stats insert: %w", err)
	}
	return nil
}

// InsertPlayerPastSeasons inserts/updates player past season history
func (r *PlayerRepo) InsertPlayerPastSeasons(pastHistory []fpl_models.PlayerPastHistoryMessage) error {
	if len(pastHistory) == 0 {
		return nil
	}

	query := sq.Insert("player_past_seasons").Columns(
		"player_code", "season_name", "season_id", "start_cost", "end_cost", "total_points",
		"minutes", "goals_scored", "assists", "clean_sheets", "goals_conceded",
		"own_goals", "penalties_saved", "penalties_missed", "yellow_cards", "red_cards",
		"saves", "bonus", "bps", "starts", "clearances_blocks_interceptions",
		"recoveries", "tackles", "defensive_contribution", "influence", "creativity",
		"threat", "ict_index", "expected_goals", "expected_assists",
		"expected_goal_involvements", "expected_goals_conceded",
	).Suffix("ON CONFLICT (player_code, season_id) DO UPDATE SET " +
		"start_cost = EXCLUDED.start_cost, " +
		"end_cost = EXCLUDED.end_cost, " +
		"total_points = EXCLUDED.total_points, " +
		"minutes = EXCLUDED.minutes, " +
		"goals_scored = EXCLUDED.goals_scored, " +
		"assists = EXCLUDED.assists, " +
		"clean_sheets = EXCLUDED.clean_sheets, " +
		"goals_conceded = EXCLUDED.goals_conceded, " +
		"own_goals = EXCLUDED.own_goals, " +
		"penalties_saved = EXCLUDED.penalties_saved, " +
		"penalties_missed = EXCLUDED.penalties_missed, " +
		"yellow_cards = EXCLUDED.yellow_cards, " +
		"red_cards = EXCLUDED.red_cards, " +
		"saves = EXCLUDED.saves, " +
		"bonus = EXCLUDED.bonus, " +
		"bps = EXCLUDED.bps, " +
		"starts = EXCLUDED.starts, " +
		"clearances_blocks_interceptions = EXCLUDED.clearances_blocks_interceptions, " +
		"recoveries = EXCLUDED.recoveries, " +
		"tackles = EXCLUDED.tackles, " +
		"defensive_contribution = EXCLUDED.defensive_contribution, " +
		"influence = EXCLUDED.influence, " +
		"creativity = EXCLUDED.creativity, " +
		"threat = EXCLUDED.threat, " +
		"ict_index = EXCLUDED.ict_index, " +
		"expected_goals = EXCLUDED.expected_goals, " +
		"expected_assists = EXCLUDED.expected_assists, " +
		"expected_goal_involvements = EXCLUDED.expected_goal_involvements, " +
		"expected_goals_conceded = EXCLUDED.expected_goals_conceded," +
		"updated_at = CURRENT_TIMESTAMP").
		PlaceholderFormat(sq.Dollar)

	for _, past := range pastHistory {
		for _, ph := range past.PlayerPastHistory {
			query = query.Values(
				past.PlayerCode, ph.SeasonName, ph.SeasonId, ph.StartCost, ph.EndCost, ph.TotalPoints,
				ph.Minutes, ph.GoalsScored, ph.Assists, ph.CleanSheets, ph.GoalsConceded,
				ph.OwnGoals, ph.PenaltiesSaved, ph.PenaltiesMissed, ph.YellowCards, ph.RedCards,
				ph.Saves, ph.Bonus, ph.BPS, ph.Starts, ph.ClearancesBlocksInterceptions,
				ph.Recoveries, ph.Tackles, ph.DefensiveContribution, ph.Influence, ph.Creativity,
				ph.Threat, ph.ICTIndex, ph.ExpectedGoals, ph.ExpectedAssists,
				ph.ExpectedGoalInvolvements, ph.ExpectedGoalsConceded,
			)
		}
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building player_past_seasons insert query: %w", err)
	}

	_, err = r.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("executing player_past_seasons insert: %w", err)
	}
	return nil
}

func (r *PlayerRepo) InsertPlayerGameweekExplain(explain []fpl_models.LiveEventMessage) error {
	if len(explain) == 0 {
		return nil
	}

	query := sq.Insert("player_gameweek_explain").Columns(
		"player_id", "fixture_id", "season_id", "event", "points",
		"identifier", "value", "points_modification",
	).Suffix("ON CONFLICT (player_id, season_id, fixture_id, identifier) DO UPDATE SET " +
		"event = EXCLUDED.event," +
		"points = EXCLUDED.points," +
		"value = EXCLUDED.value," +
		"points_modification = EXCLUDED.points_modification," +
		"updated_at = CURRENT_TIMESTAMP",
	).PlaceholderFormat(sq.Dollar)

	for _, e := range explain {
		for _, ev := range e.Explain {
			for _, detail := range ev.Stats {
				query = query.Values(
					e.PlayerID, ev.Fixture, e.SeasonID, e.Event, detail.Points,
					detail.Identifier, detail.Value, detail.PointsModification,
				)
			}
		}
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building player_gameweek_explain insert query: %w", err)
	}

	_, err = r.db.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("executing player_gameweek_explain insert: %w", err)
	}
	return nil
}
