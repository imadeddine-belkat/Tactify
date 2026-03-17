package fpl_models

// Team (Premier League club) types

type Team struct {
	ID        int    `json:"id" db:"team_id"`
	Code      int    `json:"code" db:"team_code"`
	Name      string `json:"name" db:"name"`
	ShortName string `json:"short_name" db:"short_name"`

	Strength int     `json:"strength" db:"strength"`
	Form     *string `json:"form" db:"form"`
	Position int     `json:"position" db:"position"`
	Points   int     `json:"points" db:"points"`
	Played   int     `json:"played" db:"played"`
	Win      int     `json:"win" db:"win"`
	Draw     int     `json:"draw" db:"draw"`
	Loss     int     `json:"loss" db:"loss"`

	TeamDivision *int `json:"team_division" db:"team_division"`
	Unavailable  bool `json:"unavailable" db:"unavailable"`
	PulseID      int  `json:"pulse_id" db:"pulse_id"`

	StrengthOverallHome int `json:"strength_overall_home" db:"strength_overall_home"`
	StrengthOverallAway int `json:"strength_overall_away" db:"strength_overall_away"`
	StrengthAttackHome  int `json:"strength_attack_home" db:"strength_attack_home"`
	StrengthAttackAway  int `json:"strength_attack_away" db:"strength_attack_away"`
	StrengthDefenceHome int `json:"strength_defence_home" db:"strength_defence_home"`
	StrengthDefenceAway int `json:"strength_defence_away" db:"strength_defence_away"`
}
