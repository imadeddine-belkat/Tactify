package sofascore_repositories

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/imadeddine-belkat/tactify-protos/sofascore_models"
)

type PlayerRepo struct {
	db     *sql.DB
	Player *sofascore_models.PlayerMessage
}

func NewPlayerRepo(
	db *sql.DB,
	Player *sofascore_models.PlayerMessage,
) *PlayerRepo {
	return &PlayerRepo{
		db:     db,
		Player: Player,
	}
}

func (p *PlayerRepo) InsertPlayerInfo(playerMessage sofascore_models.PlayerMessage) error {
	query := sq.Insert("players").Columns("player_id", "season_id", "team_id", "team_name", "league_id", "player_name", "player_short_name",
		"position", "height", "preferred_foot",
	).Suffix("ON CONFLICT (player_id, season_id, team_id, league_id) DO UPDATE SET " +
		"season_id=EXCLUDED.season_id, " +
		"team_id=EXCLUDED.team_id, " +
		"team_name=EXCLUDED.team_name, " +
		"league_id=EXCLUDED.league_id, " +
		"player_name=EXCLUDED.player_name, " +
		"player_short_name=EXCLUDED.player_short_name, " +
		"position=EXCLUDED.position, " +
		"height=EXCLUDED.height, " +
		"preferred_foot=EXCLUDED.preferred_foot, " +
		"updated_at= CURRENT_TIMESTAMP",
	).PlaceholderFormat(sq.Dollar)

	query = query.Values(
		playerMessage.Player.ID,
		playerMessage.SeasonID,
		playerMessage.TeamID,
		playerMessage.TeamName,
		playerMessage.LeagueID,
		playerMessage.Player.Name,
		playerMessage.Player.ShortName,
		playerMessage.Player.Position,
		playerMessage.Player.Height,
		playerMessage.Player.PreferredFoot,
	)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = p.db.Exec(sqlQuery, args...)
	if err != nil {
		return err
	}

	return nil
}
