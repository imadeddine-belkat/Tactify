package models

import "time"

type EntryTransfersMessage struct {
	EntryId   int        `json:"entry_id" db:"manager_id"`
	SeasonId  int        `json:"season_id" db:"season_id"`
	Transfers []Transfer `json:"transfers"`
}
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
