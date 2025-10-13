package sofascore_repositories

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/imadbelkat1/shared/sofascore_models"
)

type TeamRepo struct {
	db        *sql.DB
	TeamModel *sofascore_models.TopTeamsMessage
}

func NewTeamRepo(db *sql.DB, teamModel *sofascore_models.TopTeamsMessage) *TeamRepo {
	return &TeamRepo{
		db:        db,
		TeamModel: teamModel,
	}
}

func (r *TeamRepo) InsertTeamOverallStats(stats sofascore_models.TopTeamsMessage) error {
	query := sq.Insert("teams").Columns(
		"team_id", "name", "primary_color", "secondary_color",
		"league_id", "season_id",
	).Suffix(
		"ON CONFLICT (team_id, league_id, season_id) DO UPDATE SET " +
			"name = EXCLUDED.name, " +
			"primary_color = EXCLUDED.primary_color, " +
			"secondary_color = EXCLUDED.secondary_color",
	).PlaceholderFormat(sq.Dollar)

	for _, team := range stats.TopTeams.AccuratePasses {
		query = query.Values(
			team.Team.ID,
			team.Team.Name,
			team.Team.Colors.PrimaryColor,
			team.Team.Colors.SecondaryColor,
			team.LeagueID,
			team.SeasonID,
		)
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(sqlQuery, args...)
	if err != nil {
		return err
	}

	return nil
}
