package fpl_repositories

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/imadeddine-belkat/shared/fpl_models"
)

type TeamRepo struct {
	db        *sql.DB
	TeamModel *fpl_models.Team
}

func NewTeamRepo(db *sql.DB, teamModel *fpl_models.Team) *TeamRepo {
	return &TeamRepo{
		db:        db,
		TeamModel: teamModel,
	}
}

func (r *TeamRepo) InsertTeams(teams []fpl_models.TeamMessage) error {
	query := sq.Insert("teams").Columns("season_id", "team_id", "team_code", "name", "short_name", "strength",
		"form", "position", "points", "played", "win", "draw", "loss", "team_division", "unavailable",
		"pulse_id", "strength_overall_home", "strength_overall_away", "strength_attack_home", "strength_attack_away",
		"strength_defence_home", "strength_defence_away").
		Suffix("ON CONFLICT (team_id, season_id) DO UPDATE SET " + "team_code=EXCLUDED.team_code," + " name = EXCLUDED.name, " + "short_name = EXCLUDED.short_name, " +
			"strength = EXCLUDED.strength, " + "form = EXCLUDED.form, " + "position = EXCLUDED.position, " +
			"points = EXCLUDED.points, " + "played = EXCLUDED.played, " + "win = EXCLUDED.win, " +
			"draw = EXCLUDED.draw, " + "loss = EXCLUDED.loss, " + "team_division = EXCLUDED.team_division, " +
			"unavailable = EXCLUDED.unavailable, " + "pulse_id = EXCLUDED.pulse_id, " +
			"strength_overall_home = EXCLUDED.strength_overall_home, " + "strength_overall_away = EXCLUDED.strength_overall_away, " +
			"strength_attack_home = EXCLUDED.strength_attack_home, " + "strength_attack_away = EXCLUDED.strength_attack_away, " +
			"strength_defence_home = EXCLUDED.strength_defence_home, " + "strength_defence_away = EXCLUDED.strength_defence_away," +
			"updated_at = CURRENT_TIMESTAMP").
		PlaceholderFormat(sq.Dollar)

	for _, team := range teams {
		query = query.Values(
			team.SeasonID, team.Team.ID, team.Team.Code, team.Team.Name, team.Team.ShortName, team.Team.Strength,
			team.Team.Form, team.Team.Position, team.Team.Points, team.Team.Played, team.Team.Win, team.Team.Draw, team.Team.Loss, team.Team.TeamDivision, team.Team.Unavailable,
			team.Team.PulseID, team.Team.StrengthOverallHome, team.Team.StrengthOverallAway, team.Team.StrengthAttackHome, team.Team.StrengthAttackAway,
			team.Team.StrengthDefenceHome, team.Team.StrengthDefenceAway,
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
