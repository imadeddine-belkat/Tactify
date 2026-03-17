package sofascore_models

// League and tournament structures

type LeagueCategories struct {
	Categories []LeagueCategory `json:"categories"`
}

type LeagueUniqueTournaments struct {
	Groups []UniqueTournamentGroups `json:"groups"`
}

type UniqueTournamentGroups struct {
	UniqueTournament []UniqueTournament `json:"uniqueTournaments"`
}
