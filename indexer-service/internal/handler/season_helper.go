package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
)

// SeasonContext manages the current season for the application
type SeasonContext struct {
	db              *sql.DB
	currentSeasonID int
	mu              sync.RWMutex
}

func NewSeasonContext(db *sql.DB) (*SeasonContext, error) {
	sc := &SeasonContext{db: db}
	if err := sc.loadCurrentSeason(); err != nil {
		return nil, fmt.Errorf("loading current season: %w", err)
	}
	return sc, nil
}

func (sc *SeasonContext) loadCurrentSeason() error {
	query := `SELECT season_id FROM seasons WHERE is_current = true LIMIT 1`

	var seasonID int
	err := sc.db.QueryRow(query).Scan(&seasonID)
	if err != nil {
		return fmt.Errorf("querying current season: %w", err)
	}

	sc.mu.Lock()
	sc.currentSeasonID = seasonID
	sc.mu.Unlock()

	log.Printf("✅ Current season ID loaded: %d", seasonID)
	return nil
}

func (sc *SeasonContext) GetCurrentSeasonID() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.currentSeasonID
}

// EnsureSeasonExists creates a season if it doesn't exist
func (sc *SeasonContext) EnsureSeasonExists(ctx context.Context, seasonName, startDate string) (int, error) {
	query := `
		INSERT INTO seasons (season_name, start_date, is_current)
		VALUES ($1, $2, true)
		ON CONFLICT (season_name) DO UPDATE SET is_current = true
		RETURNING season_id
	`

	var seasonID int
	err := sc.db.QueryRowContext(ctx, query, seasonName, startDate).Scan(&seasonID)
	if err != nil {
		return 0, fmt.Errorf("ensuring season exists: %w", err)
	}

	sc.mu.Lock()
	sc.currentSeasonID = seasonID
	sc.mu.Unlock()

	return seasonID, nil
}
