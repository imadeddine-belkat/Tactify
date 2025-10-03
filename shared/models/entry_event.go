package models

type EntryEvent struct {
	ActiveChip    *string             `json:"active_chip"`
	AutomaticSubs []AutomaticSub      `json:"automatic_subs"`
	EntryHistory  EntryHistoryCurrent `json:"entry_history"`
	Picks         []Pick              `json:"picks"`
}

type AutomaticSub struct {
	Entry      int `json:"entry" db:"entry_id"`
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
