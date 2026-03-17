package fpl_models

// ElementType represents player positions (GKP, DEF, MID, FWD) and squad rules

type ElementType struct {
	ID                int    `json:"id" db:"id"`
	SingularName      string `json:"singular_name" db:"singular_name"`
	PluralName        string `json:"plural_name" db:"plural_name"`
	SingularNameShort string `json:"singular_name_short" db:"singular_name_short"`
	PluralNameShort   string `json:"plural_name_short" db:"plural_name_short"`

	SquadSelect    int  `json:"squad_select" db:"squad_select"`
	SquadMinSelect *int `json:"squad_min_select" db:"squad_min_select"`
	SquadMaxSelect *int `json:"squad_max_select" db:"squad_max_select"`
	SquadMinPlay   int  `json:"squad_min_play" db:"squad_min_play"`
	SquadMaxPlay   int  `json:"squad_max_play" db:"squad_max_play"`

	UIShirtSpecific    bool  `json:"ui_shirt_specific" db:"ui_shirt_specific"`
	SubPositionsLocked []int `json:"sub_positions_locked" db:"sub_positions_locked"`
	ElementCount       int   `json:"element_count" db:"element_count"`
}
