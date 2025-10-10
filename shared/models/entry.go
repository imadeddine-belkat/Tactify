package models

type EntryMessage struct {
	Entry    Entry `json:"entry"`
	SeasonId int   `json:"season_id" db:"season_id"`
}
type Entry struct {
	ID                int    `json:"id" db:"manager_id"`
	Name              string `json:"name" db:"manager_name"`
	PlayerFirstName   string `json:"player_first_name" db:"player_first_name"`
	PlayerLastName    string `json:"player_last_name" db:"player_last_name"`
	PlayerRegionID    int    `json:"player_region_id" db:"player_region_id"`
	PlayerRegionName  string `json:"player_region_name" db:"player_region_name"`
	PlayerRegionShort string `json:"player_region_iso_code_short" db:"player_region_iso_code_short"`
	PlayerRegionLong  string `json:"player_region_iso_code_long" db:"player_region_iso_code_long"`
	FavouriteTeam     int    `json:"favourite_team" db:"favourite_team"`

	JoinedTime   string `json:"joined_time" db:"joined_time"`
	StartedEvent int    `json:"started_event" db:"started_event"`
	YearsActive  int    `json:"years_active" db:"years_active"`

	SummaryOverallPoints int `json:"summary_overall_points" db:"summary_overall_points"`
	SummaryOverallRank   int `json:"summary_overall_rank" db:"summary_overall_rank"`
	SummaryEventPoints   int `json:"summary_event_points" db:"summary_event_points"`
	SummaryEventRank     int `json:"summary_event_rank" db:"summary_event_rank"`
	CurrentEvent         int `json:"current_event" db:"current_event"`

	NameChangeBlocked bool `json:"name_change_blocked" db:"name_change_blocked"`

	LastDeadlineBank      int    `json:"last_deadline_bank" db:"last_deadline_bank"`
	LastDeadlineValue     int    `json:"last_deadline_value" db:"last_deadline_value"`
	LastDeadlineTransfers int    `json:"last_deadline_total_transfers" db:"last_deadline_total_transfers"`
	ClubBadgeSrc          string `json:"club_badge_src" db:"club_badge_src"`
}
