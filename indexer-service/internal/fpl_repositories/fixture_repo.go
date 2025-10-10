package fpl_repositories

import (
	"database/sql"

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

func (r *FixtureRepo) InsertFixtures(fixtures []models.FixtureMessage) error {
	if len(fixtures) == 0 {
		return nil
	}

	fixtureInsert := sq.Insert("fixtures").
		Columns("fixture_id", "season_id", "fixture_code", "event", "team_h", "team_a", "kickoff_time",
			"team_h_score", "team_a_score", "finished", "minutes", "provisional_start_time",
			"team_h_difficulty", "team_a_difficulty", "pulse_id").
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

	for _, f := range fixtures {
		fixtureInsert = fixtureInsert.Values(
			f.Fixture.ID, f.SeasonID, f.Fixture.Code, f.Fixture.Event,
			f.Fixture.TeamH, f.Fixture.TeamA, nullIfEmpty(f.Fixture.KickoffTime),
			f.Fixture.TeamHScore, f.Fixture.TeamAScore, f.Fixture.Finished,
			f.Fixture.Minutes, f.Fixture.ProvisionalStartTime,
			f.Fixture.TeamHDifficulty, f.Fixture.TeamADifficulty, f.Fixture.PulseID,
		)

		for _, stat := range f.Fixture.Stats {
			for _, p := range stat.H {
				fixtureStatsInsert = fixtureStatsInsert.Values(
					f.Fixture.ID, f.SeasonID, p.Element, stat.Identifier, p.Value,
				)
			}
			for _, p := range stat.A {
				fixtureStatsInsert = fixtureStatsInsert.Values(
					f.Fixture.ID, f.SeasonID, p.Element, stat.Identifier, p.Value,
				)
			}
		}
	}

	q1, args1, err := fixtureInsert.ToSql()
	if err != nil {
		return err
	}
	if _, err = r.db.Exec(q1, args1...); err != nil {
		return err
	}

	q2, args2, err := fixtureStatsInsert.ToSql()
	if err != nil {
		return err
	}
	if len(args2) > 0 {
		_, err = r.db.Exec(q2, args2...)
		if err != nil {
			return err
		}
	}

	return nil
}
