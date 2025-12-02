package sofascore_repositories

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/imadeddine-belkat/shared/sofascore_models"
)

type MatchRepo struct {
	db    *sql.DB
	Event *sofascore_models.Event
}

func NewMatchRepo(
	db *sql.DB,
	event *sofascore_models.Event,
) *MatchRepo {
	return &MatchRepo{
		db:    db,
		Event: event,
	}
}

func (m *MatchRepo) InsertRoundMatches(match sofascore_models.Event) error {
	query := sq.Insert("matches").Columns(
		"match_id", "season_id", "league_id", "home_team_id", "away_team_id",
		"home_team_name", "away_team_name", "start_time", "round", "status", "status_description").
		Suffix("ON CONFLICT (match_id, season_id, league_id) DO UPDATE SET " +
			"home_team_id = EXCLUDED.home_team_id, " +
			"away_team_id = EXCLUDED.away_team_id, " +
			"home_team_name = EXCLUDED.home_team_name, " +
			"away_team_name = EXCLUDED.away_team_name, " +
			"start_time = EXCLUDED.start_time," +
			"round = EXCLUDED.round," +
			"status = EXCLUDED.status," +
			"status_description = EXCLUDED.status_description," +
			"updated_at = CURRENT_TIMESTAMP").
		PlaceholderFormat(sq.Dollar)

	query = query.Values(match.ID, match.Season.ID, match.Tournament.UniqueTournament.ID,
		match.HomeTeam.ID, match.AwayTeam.ID, match.HomeTeam.Name, match.AwayTeam.Name,
		time.Unix(match.StartTimestamp, 0), match.RoundInfo.Round,
		match.Status.Code, match.Status.Description,
	)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = m.db.Exec(sqlQuery, args...)
	if err != nil {
		return err
	}

	return nil
}
