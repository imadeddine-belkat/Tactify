package fpl_models

import "time"

// Manager (Entry) and team management types

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

// Entry event picks and lineup
type EntryEventPicks struct {
	ActiveChip    *string             `json:"active_chip"`
	AutomaticSubs []AutomaticSub      `json:"automatic_subs"`
	EntryHistory  EntryHistoryCurrent `json:"entry_history"`
	Picks         []Pick              `json:"picks"`
}

type AutomaticSub struct {
	Entry      int `json:"entry" db:"manager_id"`
	ElementIn  int `json:"element_in" db:"element_in"`   // player id in
	ElementOut int `json:"element_out" db:"element_out"` // player id out
	Event      int `json:"event" db:"event"`
}

type Pick struct {
	Element       int  `json:"element" db:"player_id"`     // player id
	Position      int  `json:"position" db:"position"`     // 1-15
	Multiplier    int  `json:"multiplier" db:"multiplier"` // 0 if bench, 2 if captain
	IsCaptain     bool `json:"is_captain" db:"is_captain"`
	IsViceCaptain bool `json:"is_vice_captain" db:"is_vice_captain"`
	ElementType   int  `json:"element_type" db:"element_type"` // GKP, DEF, MID, FWD
}

// Entry history
type EntryHistory struct {
	Current []EntryHistoryCurrent `json:"current"`
	Past    []EntryHistoryPast    `json:"past"`
	Chips   []EntryHistoryChip    `json:"chips"`
}

type EntryHistoryCurrent struct {
	Event              int `json:"event" db:"event"`
	Points             int `json:"points" db:"points"`
	TotalPoints        int `json:"total_points" db:"total_points"`
	Rank               int `json:"rank" db:"rank"`
	RankSort           int `json:"rank_sort" db:"rank_sort"`
	OverallRank        int `json:"overall_rank" db:"overall_rank"`
	PercentileRank     int `json:"percentile_rank" db:"percentile_rank"`
	Bank               int `json:"bank" db:"bank"`
	Value              int `json:"value" db:"value"` // value is in tenths of £m
	EventTransfers     int `json:"event_transfers" db:"transfers"`
	EventTransfersCost int `json:"event_transfers_cost" db:"transfers_cost"`
	PointsOnBench      int `json:"points_on_bench" db:"points_on_bench"`
}

type EntryHistoryPast struct {
	SeasonName  string `json:"season_name" db:"season_name"`
	SeasonId    int    `json:"season_id" db:"season_id"`
	TotalPoints int    `json:"total_points" db:"total_points"`
	Rank        int    `json:"rank" db:"rank"`
}

type EntryHistoryChip struct {
	Name  string `json:"name" db:"name"`
	Time  string `json:"time" db:"time"`
	Event int    `json:"event" db:"event"`
}

// Entry transfers
type EntryTransfers struct {
	Transfers []Transfer `json:"transfers"`
}

type Transfer struct {
	Entry          int       `json:"entry" db:"manager_id"`
	ElementIn      int       `json:"element_in" db:"element_in"`             // player id in
	ElementInCost  int       `json:"element_in_cost" db:"element_in_cost"`   // cost of player in
	ElementOut     int       `json:"element_out" db:"element_out"`           // player id out
	ElementOutCost int       `json:"element_out_cost" db:"element_out_cost"` // cost of player out
	Event          int       `json:"event" db:"event"`                       // gameweek
	Time           time.Time `json:"time" db:"time"`
}
