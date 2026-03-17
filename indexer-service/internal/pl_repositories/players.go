package pl_repositories

import (
	"context"
	"database/sql"

	"github.com/imadeddine-belkat/indexer-service/internal/db/helpers"
	fpl "github.com/imadeddine-belkat/tactify-protos/go/fpl/v1"
)

type PlayerRepo struct {
	db     *sql.DB
	Player *fpl.PlayerBootstrapMessage
}

// InsertPlayers inserts/updates main player information
func (r *PlayerRepo) InsertPlayers(ctx context.Context, players []*fpl.PlayerBootstrapMessage) error {
	rows := make([][]any, len(players))
	for i, p := range players {
		rows[i] = []any{
			p.Player.GetCode(),
			p.Player.GetFirstName(),
			p.Player.GetSecondName(),
			p.Player.GetWebName(),
			p.Player.GetBirthDate(),
			p.Player.GetOptaCode(),
		}
	}

	return helpers.BatchUpsert(ctx, r.db, helpers.UpsertOpts{
		Table:        "players",
		Columns:      []string{"code", "first_name", "second_name", "web_name", "birth_date", "opta_code"},
		ConflictCols: []string{"code"},
		Rows:         rows,
	})
}

func (r *PlayerRepo) InsertPlayerSeasons(ctx context.Context, players []*fpl.PlayerBootstrapMessage) error {
	rows := make([][]any, len(players))
	for i, p := range players {
		rows[i] = []any{
			p.Player.GetCode(),
			p.Player.GetFirstName(),
			p.Player.GetSecondName(),
			p.Player.GetWebName(),
			p.Player.GetBirthDate(),
			p.Player.GetOptaCode(),
		}
	}

	return helpers.BatchUpsert(ctx, r.db, helpers.UpsertOpts{
		Table:        "players",
		Columns:      []string{"code", "first_name", "second_name", "web_name", "birth_date", "opta_code"},
		ConflictCols: []string{"code"},
		Rows:         rows,
	})
}
