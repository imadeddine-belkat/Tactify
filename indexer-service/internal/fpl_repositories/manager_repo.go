package fpl_repositories

import (
	"database/sql"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/imadeddine-belkat/shared/fpl_models"
)

type ManagerRepo struct {
	db             *sql.DB
	Entry          *fpl_models.EntryMessage
	EntryPicks     *fpl_models.EntryEventPicksMessage
	EntryTransfers *fpl_models.EntryTransfersMessage
	EntryHistory   *fpl_models.EntryHistoryMessage
}

func NewManagerRepo(db *sql.DB, entry *fpl_models.EntryMessage, entryPicks *fpl_models.EntryEventPicksMessage, entryTransfers *fpl_models.EntryTransfersMessage, entryHistory *fpl_models.EntryHistoryMessage) *ManagerRepo {
	return &ManagerRepo{
		db:             db,
		Entry:          entry,
		EntryPicks:     entryPicks,
		EntryTransfers: entryTransfers,
		EntryHistory:   entryHistory,
	}
}

func (r *ManagerRepo) InsertManagerInfo(entry *fpl_models.EntryMessage) error {
	if entry == nil {
		return nil
	}

	query := sq.Insert("managers").Columns(
		"manager_id", "season_id", "manager_name", "player_first_name", "player_last_name", "player_region_id",
		"player_region_name", "player_region_iso_code_short", "player_region_iso_code_long", "favourite_team_id", "joined_time",
		"started_event", "years_active", "summary_overall_points", "summary_overall_rank", "summary_points", "summary_rank",
		"current_event", "name_change_blocked", "last_deadline_bank", "last_deadline_value", "last_deadline_total_transfers",
		"club_badge_src",
	).Values(
		entry.Entry.ID, entry.SeasonId, entry.Entry.Name, entry.Entry.PlayerFirstName, entry.Entry.PlayerLastName, entry.Entry.PlayerRegionID,
		entry.Entry.PlayerRegionName, entry.Entry.PlayerRegionShort, entry.Entry.PlayerRegionLong, entry.Entry.FavouriteTeam, entry.Entry.JoinedTime,
		entry.Entry.StartedEvent, entry.Entry.YearsActive, entry.Entry.SummaryOverallPoints, entry.Entry.SummaryOverallRank, entry.Entry.SummaryEventPoints, entry.Entry.SummaryEventRank,
		entry.Entry.CurrentEvent, entry.Entry.NameChangeBlocked, entry.Entry.LastDeadlineBank, entry.Entry.LastDeadlineValue, entry.Entry.LastDeadlineTransfers,
		entry.Entry.ClubBadgeSrc,
	).Suffix(`ON CONFLICT (manager_id, season_id) DO UPDATE SET 
		manager_name = EXCLUDED.manager_name,
		player_first_name = EXCLUDED.player_first_name,
		player_last_name = EXCLUDED.player_last_name,
		player_region_id = EXCLUDED.player_region_id,
		player_region_name = EXCLUDED.player_region_name,
		player_region_iso_code_short = EXCLUDED.player_region_iso_code_short,
		player_region_iso_code_long = EXCLUDED.player_region_iso_code_long,
		favourite_team_id = EXCLUDED.favourite_team_id,
		joined_time = EXCLUDED.joined_time,
		started_event = EXCLUDED.started_event,
		years_active = EXCLUDED.years_active,
		summary_overall_points = EXCLUDED.summary_overall_points,
		summary_overall_rank = EXCLUDED.summary_overall_rank,
		summary_points = EXCLUDED.summary_points,
		summary_rank = EXCLUDED.summary_rank,
		current_event = EXCLUDED.current_event,
		name_change_blocked = EXCLUDED.name_change_blocked,
		last_deadline_bank = EXCLUDED.last_deadline_bank,
		last_deadline_value = EXCLUDED.last_deadline_value,
		last_deadline_total_transfers = EXCLUDED.last_deadline_total_transfers,
		club_badge_src = EXCLUDED.club_badge_src,
		updated_at = CURRENT_TIMESTAMP`,
	).PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building managers insert query: %w", err)
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
	return err
}

func (r *ManagerRepo) InsertManagerPicks(entryPicks *fpl_models.EntryEventPicksMessage) error {
	if entryPicks == nil {
		return nil
	}

	queryPicks := sq.Insert("manager_picks").Columns(
		"manager_id", "season_id", "event", "player_id", "is_captain",
		"is_vice_captain", "multiplier", "position", "element_type",
	).Suffix(`ON CONFLICT (manager_id, season_id, event, player_id) DO UPDATE SET 
		is_captain = EXCLUDED.is_captain,
		is_vice_captain = EXCLUDED.is_vice_captain,
		multiplier = EXCLUDED.multiplier,
		position = EXCLUDED.position,
		element_type = EXCLUDED.element_type,
		updated_at = CURRENT_TIMESTAMP`,
	).PlaceholderFormat(sq.Dollar)

	querySubs := sq.Insert("manager_automatic_subs").Columns(
		"manager_id", "season_id", "event", "player_out_id", "player_in_id",
	).Suffix(`ON CONFLICT (manager_id, season_id, event, player_out_id, player_in_id) DO UPDATE SET 
		updated_at = CURRENT_TIMESTAMP`,
	).PlaceholderFormat(sq.Dollar)

	for _, pick := range entryPicks.Picks.Picks {
		queryPicks = queryPicks.Values(
			entryPicks.EntryId, entryPicks.SeasonId, entryPicks.EventId, pick.Element, pick.IsCaptain,
			pick.IsViceCaptain, pick.Multiplier, pick.Position, pick.ElementType,
		)
	}

	for _, sub := range entryPicks.Picks.AutomaticSubs {
		querySubs = querySubs.Values(
			entryPicks.EntryId, entryPicks.SeasonId, entryPicks.EventId, sub.ElementOut, sub.ElementIn,
		)
	}

	q1, args, err := queryPicks.ToSql()
	if err != nil {
		return fmt.Errorf("building manager_picks insert query: %w", err)
	}

	q2, args2, err := querySubs.ToSql()
	if err != nil {
		return fmt.Errorf("building manager_automatic_subs insert query: %w", err)
	}

	result, err := r.db.Exec(q1, args...)
	if err != nil {
		log.Printf("❌ SQL Error inserting manager_picks: %v", err)
		log.Printf("   Query length: %d characters", len(q1))
		log.Printf("   Args count: %d", len(args))
		return fmt.Errorf("executing manager_picks insert: %w", err)
	}

	if len(args2) > 0 {
		result2, err := r.db.Exec(q2, args2...)
		if err != nil {
			log.Printf("❌ SQL Error inserting manager_automatic_subs: %v", err)
			log.Printf("   Query length: %d characters", len(q2))
			log.Printf("   Args count: %d", len(args2))
			return fmt.Errorf("executing manager_automatic_subs insert: %w", err)
		}

		rowsAffected2, _ := result2.RowsAffected()
		log.Printf("✅ manager_automatic_subs insert completed: %d rows affected", rowsAffected2)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ manager_picks insert completed: %d rows affected", rowsAffected)
	return err
}

func (r *ManagerRepo) InsertManagerTransfers(entryTransfers *fpl_models.EntryTransfersMessage) error {
	if entryTransfers == nil {
		return nil
	}

	query := sq.Insert("manager_transfers").Columns(
		"manager_id", "season_id", "event", "player_in_id", "player_in_cost",
		"player_out_id", "player_out_cost").
		Suffix(`ON CONFLICT (manager_id, season_id, event, player_in_id, player_out_id) DO UPDATE SET 
		player_in_cost = EXCLUDED.player_in_cost,
		player_out_cost = EXCLUDED.player_out_cost,
		updated_at = CURRENT_TIMESTAMP`,
		).PlaceholderFormat(sq.Dollar)

	for _, transfer := range entryTransfers.Transfers {
		query = query.Values(
			entryTransfers.EntryId, entryTransfers.SeasonId, transfer.Event, transfer.ElementIn, transfer.ElementInCost,
			transfer.ElementOut, transfer.ElementOutCost,
		)
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building manager_transfers insert query: %w", err)
	}

	result, err := r.db.Exec(sqlQuery, args...)
	if err != nil {
		log.Printf("❌ SQL Error inserting manager_transfers: %v", err)
		log.Printf("   Query length: %d characters", len(sqlQuery))
		log.Printf("   Args count: %d", len(args))
		return fmt.Errorf("executing manager_transfers insert: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ manager_transfers insert completed: %d rows affected", rowsAffected)
	return err
}

func (r *ManagerRepo) InsertManagerFullHistory(entryHistory *fpl_models.EntryHistoryMessage) error {
	if entryHistory == nil {
		return nil
	}

	queryCurrent := sq.Insert("manager_gameweek_history").Columns(
		"manager_id", "season_id", "event", "points", "total_points", "rank",
		"rank_sort", "overall_rank", "percentile_rank", "bank", "value", "event_transfers", "event_transfers_cost",
		"points_on_bench",
	).Suffix(`ON CONFLICT (manager_id, season_id, event) DO UPDATE SET 
		points = EXCLUDED.points,
		total_points = EXCLUDED.total_points,
		rank = EXCLUDED.rank,
		rank_sort = EXCLUDED.rank_sort,
		overall_rank = EXCLUDED.overall_rank,
		percentile_rank = EXCLUDED.percentile_rank,
		bank = EXCLUDED.bank,
		value = EXCLUDED.value,
		event_transfers = EXCLUDED.event_transfers,
		event_transfers_cost = EXCLUDED.event_transfers_cost,
		points_on_bench = EXCLUDED.points_on_bench,
		updated_at = CURRENT_TIMESTAMP`,
	).PlaceholderFormat(sq.Dollar)

	queryPast := sq.Insert("manager_season_history").Columns(
		"manager_id", "season_id", "season_name", "total_points", "rank",
	).Suffix(`ON CONFLICT (manager_id, season_id) DO UPDATE SET 
		total_points = EXCLUDED.total_points,
		rank = EXCLUDED.rank,
		updated_at = CURRENT_TIMESTAMP`,
	).PlaceholderFormat(sq.Dollar)

	queryChips := sq.Insert("manager_chips").Columns(
		"manager_id", "season_id", "event", "chip_name", "time",
	).Suffix(`ON CONFLICT (manager_id, season_id, event, chip_name) DO UPDATE SET 
		time = EXCLUDED.time,
		updated_at = CURRENT_TIMESTAMP`,
	).PlaceholderFormat(sq.Dollar)

	for _, history := range entryHistory.EntryHistory.Current {
		queryCurrent = queryCurrent.Values(
			entryHistory.EntryId, entryHistory.SeasonId, history.Event, history.Points, history.TotalPoints, history.Rank,
			history.RankSort, history.OverallRank, history.PercentileRank, history.Bank, history.Value, history.EventTransfers, history.EventTransfersCost,
			history.PointsOnBench,
		)
	}

	for _, history := range entryHistory.EntryHistory.Past {
		queryPast = queryPast.Values(
			entryHistory.EntryId, history.SeasonId, history.SeasonName, history.TotalPoints, history.Rank,
		)
	}

	for _, chip := range entryHistory.EntryHistory.Chips {
		queryChips = queryChips.Values(
			entryHistory.EntryId, entryHistory.SeasonId, chip.Event, chip.Name, chip.Time,
		)
	}

	q1, args1, err := queryCurrent.ToSql()
	if err != nil {
		return fmt.Errorf("building manager_gameweek_history insert query: %w", err)
	}

	q2, args2, err := queryPast.ToSql()
	if err != nil {
		return fmt.Errorf("building manager_season_history insert query: %w", err)
	}

	q3, args3, err := queryChips.ToSql()
	if err != nil {
		return fmt.Errorf("building manager_chips insert query: %w", err)
	}

	result, err := r.db.Exec(q1, args1...)
	if err != nil {
		log.Printf("❌ SQL Error inserting manager_gameweek_history: %v", err)
		log.Printf("   Query length: %d characters", len(q1))
		log.Printf("   Args count: %d", len(args1))
		return fmt.Errorf("executing manager_gameweek_history insert: %w", err)
	}

	var rowsAffected2 int64
	if len(args2) > 0 {
		result2, err := r.db.Exec(q2, args2...)
		if err != nil {
			log.Printf("❌ SQL Error inserting manager_season_history: %v", err)
			log.Printf("   Query length: %d characters", len(q2))
			log.Printf("   Args count: %d", len(args2))
			return fmt.Errorf("executing manager_season_history insert: %w", err)
		}

		rowsAffected2, _ = result2.RowsAffected()
		log.Printf("✅ manager_season_history insert completed: %d rows affected", rowsAffected2)
	}

	var rowsAffected3 int64
	if len(args3) > 0 {
		result3, err := r.db.Exec(q3, args3...)
		if err != nil {
			log.Printf("❌ SQL Error inserting manager_chips: %v", err)
			log.Printf("   Query length: %d characters", len(q3))
			log.Printf("   Args count: %d", len(args3))
			return fmt.Errorf("executing manager_chips insert: %w", err)
		}

		rowsAffected3, _ = result3.RowsAffected()
		log.Printf("✅ manager_chips insert completed: %d rows affected", rowsAffected3)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ manager_gameweek_history insert completed: %d rows affected", rowsAffected)
	return err
}
