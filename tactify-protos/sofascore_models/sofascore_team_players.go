package sofascore_models

type TeamPlayers struct {
	TeamPlayers []Players `json:"players"`
}

type Players struct {
	Player Player `json:"player"`
}
