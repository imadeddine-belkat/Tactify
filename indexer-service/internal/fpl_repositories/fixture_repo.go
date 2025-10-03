package fpl_repositories

import (
	"database/sql"

	models "github.com/imadbelkat1/shared/models"
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

func (r *FixtureRepo) InsertFixtures(fixtures []models.Fixture) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
		INSERT INTO fixtures (fixture_id, event, kickoff_time,  team_h, team_a, team_h_score, team_a_score, finished, minutes, provisional_start_time, started, finished_provisional, pulse_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (fixture_id) DO UPDATE SET
			event = EXCLUDED.event,
			kickoff_time = EXCLUDED.kickoff_time,
			team_h = EXCLUDED.team_h,
			team_a = EXCLUDED.team_a,
			team_h_score = EXCLUDED.team_h_score,
			team_a_score = EXCLUDED.team_a_score,
			finished = EXCLUDED.finished,
			minutes = EXCLUDED.minutes,
			provisional_start_time = EXCLUDED.provisional_start_time,
			started = EXCLUDED.started,
			finished_provisional = EXCLUDED.finished_provisional,
			pulse_id = EXCLUDED.pulse_id
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, fixture := range fixtures {
		_, err := stmt.Exec(
			fixture.ID,
			fixture.Event,
			fixture.KickoffTime,
			fixture.TeamH,
			fixture.TeamA,
			fixture.TeamHScore,
			fixture.TeamAScore,
			fixture.Finished,
			fixture.Minutes,
			fixture.ProvisionalStartTime,
			fixture.Started,
			fixture.FinishedProvisional,
			fixture.PulseID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
