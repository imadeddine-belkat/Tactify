package models

type EntryHistoryMessage struct {
	EntryHistory EntryHistory `json:"entry_history"`
	EntryId      int          `json:"entry_id" db:"manager_id"`
	SeasonId     int          `json:"season_id" db:"season_id"`
}
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
	Value              int `json:"value" db:"value"` //value is in tenths of £m
	EventTransfers     int `json:"event_transfers" db:"transfers"`
	EventTransfersCost int `json:"event_transfers_cost" db:"transfers_cost"`
	PointsOnBench      int `json:"points_on_bench" db:"points_on_bench"`
}

type EntryHistoryPast struct {
	SeasonName  string `json:"season_name" db:"season_name"`
	TotalPoints int    `json:"total_points" db:"total_points"`
	Rank        int    `json:"rank" db:"rank"`
}

type EntryHistoryChip struct {
	Name  string `json:"name" db:"name"`
	Time  string `json:"time" db:"time"`
	Event int    `json:"event" db:"event"`
}
