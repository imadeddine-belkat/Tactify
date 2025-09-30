package models

type Fixture struct {
	Code                 int           `json:"code" db:"code"`
	Event                int           `json:"event" db:"event"`
	Finished             bool          `json:"finished" db:"finished"`
	FinishedProvisional  bool          `json:"finished_provisional" db:"finished_provisional"`
	ID                   int           `json:"id" db:"id"`
	KickoffTime          string        `json:"kickoff_time" db:"kickoff_time"`
	Minutes              int           `json:"minutes" db:"minutes"`
	ProvisionalStartTime bool          `json:"provisional_start_time" db:"provisional_start_time"`
	Started              bool          `json:"started" db:"started"`
	TeamA                int           `json:"team_a" db:"team_a"`
	TeamAScore           int           `json:"team_a_score" db:"team_a_score"`
	TeamH                int           `json:"team_h" db:"team_h"`
	TeamHScore           int           `json:"team_h_score" db:"team_h_score"`
	Stats                []FixtureStat `json:"stats,omitempty"`
	TeamHDifficulty      int           `json:"team_h_difficulty" db:"team_h_difficulty"`
	TeamADifficulty      int           `json:"team_a_difficulty" db:"team_a_difficulty"`
	PulseID              int           `json:"pulse_id" db:"pulse_id"`
}

type FixtureStat struct {
	Identifier string        `json:"identifier" db:"identifier"`
	A          []StatElement `json:"a"`
	H          []StatElement `json:"h"`
}

type StatElement struct {
	Value   int `json:"value" db:"value"`
	Element int `json:"element" db:"player_id"`
}

type Fixtures []Fixture
