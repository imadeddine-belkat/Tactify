package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/imadeddine-belkat/indexer-service/config"
)

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type UpsertOpts struct {
	Table        string
	Columns      []string
	ConflictCols []string
	SkipUpdate   []string
	Rows         [][]any
	config       *config.IndexerConfig
}

func BatchUpsert(ctx context.Context, db Executor, opts UpsertOpts) error {
	if len(opts.Rows) == 0 {
		return nil
	}

	batchSize := opts.config.BatchSize

	suffix := buildUpsertSuffix(opts.Columns, opts.ConflictCols, opts.SkipUpdate)

	for i := 0; i < len(opts.Rows); i += batchSize {
		end := min(i+batchSize, len(opts.Rows))
		if err := execChunk(ctx, db, opts.Table, opts.Columns, suffix, opts.Rows[i:end]); err != nil {
			return fmt.Errorf("batch [%d:%d]: %w", i, end, err)
		}
	}
	return nil
}

func execChunk(ctx context.Context, db Executor, table string, cols []string, suffix string, rows [][]any) error {
	query := sq.Insert(table).
		Columns(cols...).
		PlaceholderFormat(sq.Dollar)

	if suffix != "" {
		query = query.Suffix(suffix)
	}

	for _, row := range rows {
		query = query.Values(row...)
	}

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("building query: %w", err)
	}

	_, err = db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("exec (%d rows, %d args): %w", len(rows), len(args), err)
	}
	return nil
}

func buildUpsertSuffix(cols, conflictCols, skipCols []string) string {
	if len(conflictCols) == 0 {
		return ""
	}

	exclude := make(map[string]struct{}, len(conflictCols)+len(skipCols))
	for _, c := range conflictCols {
		exclude[c] = struct{}{}
	}
	for _, c := range skipCols {
		exclude[c] = struct{}{}
	}

	var setClauses []string
	for _, c := range cols {
		if _, skip := exclude[c]; skip {
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = EXCLUDED.%s", c, c))
	}

	if len(setClauses) == 0 {
		return fmt.Sprintf("ON CONFLICT (%s) DO NOTHING", strings.Join(conflictCols, ", "))
	}

	return fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s",
		strings.Join(conflictCols, ", "),
		strings.Join(setClauses, ", "),
	)
}
