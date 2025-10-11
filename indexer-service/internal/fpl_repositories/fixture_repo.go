package fpl_repositories

import (
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/imadbelkat1/shared/models"
)

type FixtureRepo struct {
	db           *sql.DB
	FixtureModel *models.Fixture
}

func NewFixtureRepo(db *sql.DB, fixtureModel *models.Fixture) *FixtureRepo {
	return &FixtureRepo{
		db:           db,
		FixtureModel: fixtureModel,
	}
}

// In fixture_repo.go
func (r *FixtureRepo) InsertFixtures(fixtures []models.FixtureMessage) error {
	if len(fixtures) == 0 {
		return nil
	}

	fixtureInsert := sq.Insert("fixtures").
		Columns("fixture_id", "season_id", "fixture_code", "event", "team_h", "team_a",
			"kickoff_time", "team_h_score", "team_a_score", "finished", "minutes",
			"provisional_start_time", "team_h_difficulty", "team_a_difficulty", "pulse_id").
		Suffix(`ON CONFLICT (fixture_id, season_id)
                DO UPDATE SET
                fixture_code = EXCLUDED.fixture_code,
                event = EXCLUDED.event,
                team_h = EXCLUDED.team_h,
                team_a = EXCLUDED.team_a,
                kickoff_time = EXCLUDED.kickoff_time,
                team_h_score = EXCLUDED.team_h_score,
                team_a_score = EXCLUDED.team_a_score,
                finished = EXCLUDED.finished,
                minutes = EXCLUDED.minutes,
                provisional_start_time = EXCLUDED.provisional_start_time,
                team_h_difficulty = EXCLUDED.team_h_difficulty,
                team_a_difficulty = EXCLUDED.team_a_difficulty,
                pulse_id = EXCLUDED.pulse_id,
                updated_at = CURRENT_TIMESTAMP`).
		PlaceholderFormat(sq.Dollar)

	fixtureStatsInsert := sq.Insert("fixture_stats").
		Columns("fixture_id", "season_id", "player_id", "identifier", "value").
		Suffix(`ON CONFLICT (fixture_id, player_id, identifier, season_id)
                DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP`).
		PlaceholderFormat(sq.Dollar)

	hasStats := false

	// Add all fixtures to the batch
	for _, fixtureMsg := range fixtures {
		fixtureInsert = fixtureInsert.Values(
			fixtureMsg.Fixture.ID, fixtureMsg.SeasonID, fixtureMsg.Fixture.Code,
			fixtureMsg.Fixture.Event, fixtureMsg.Fixture.TeamH, fixtureMsg.Fixture.TeamA,
			nullIfEmpty(fixtureMsg.Fixture.KickoffTime), fixtureMsg.Fixture.TeamHScore,
			fixtureMsg.Fixture.TeamAScore, fixtureMsg.Fixture.Finished,
			fixtureMsg.Fixture.Minutes, fixtureMsg.Fixture.ProvisionalStartTime,
			fixtureMsg.Fixture.TeamHDifficulty, fixtureMsg.Fixture.TeamADifficulty,
			fixtureMsg.Fixture.PulseID,
		)

		// Add all stats for this fixture
		for _, stat := range fixtureMsg.Fixture.Stats {
			for _, p := range stat.H {
				fixtureStatsInsert = fixtureStatsInsert.Values(
					fixtureMsg.Fixture.ID, fixtureMsg.SeasonID, p.Element,
					stat.Identifier, p.Value,
				)
				hasStats = true
			}
			for _, p := range stat.A {
				fixtureStatsInsert = fixtureStatsInsert.Values(
					fixtureMsg.Fixture.ID, fixtureMsg.SeasonID, p.Element,
					stat.Identifier, p.Value,
				)
				hasStats = true
			}
		}
	}

	// Execute fixture batch insert
	q1, args1, err := fixtureInsert.ToSql()
	if err != nil {
		return err
	}
	result, err := r.db.Exec(q1, args1...)
	if err != nil {
		return err
	}
	fixtureRowsAffected, _ := result.RowsAffected()

	// Execute stats batch insert if we have stats
	statsRowsAffected := int64(0)
	if hasStats {
		q2, args2, err := fixtureStatsInsert.ToSql()
		if err != nil {
			return err
		}
		result2, err := r.db.Exec(q2, args2...)
		if err != nil {
			return err
		}
		statsRowsAffected, _ = result2.RowsAffected()
	}

	log.Printf("✅ Batch inserted %d fixtures: %d fixture rows, %d stats rows",
		len(fixtures), fixtureRowsAffected, statsRowsAffected)

	return nil
}
