package sofascore_repositories

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	sofascore "github.com/imadeddine-belkat/tactify-protos/go/sofascore/v1"
)

type PlayerRepo struct {
	db     *sql.DB
	Player *sofascore.PlayerMessage
}

func NewPlayerRepo(
	db *sql.DB,
	Player *sofascore.PlayerMessage,
) *PlayerRepo {
	return &PlayerRepo{
		db:     db,
		Player: Player,
	}
}

func (p *PlayerRepo) InsertPlayerInfo(playerMessage *sofascore.PlayerMessage) error {
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
		playerMessage.Player.Id,
		playerMessage.SeasonId,
		playerMessage.TeamId,
		playerMessage.TeamName,
		playerMessage.LeagueId,
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
