package fpl_models

// BootstrapResponse represents the main FPL API bootstrap-static response
// This aggregates all core game data in a single response
type BootstrapResponse struct {
	Events       []Event           `json:"events"`
	GameSettings GameSettings      `json:"game_settings"`
	Phases       []Phase           `json:"phases"`
	Teams        []Team            `json:"teams"`
	TotalPlayers int               `json:"total_players"`
	Elements     []PlayerBootstrap `json:"elements"`
	ElementStats []ElementStat     `json:"element_stats"`
	ElementTypes []ElementType     `json:"element_types"`
}
