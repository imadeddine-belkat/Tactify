package sofascore_models

// League standings types

type Standings struct {
	Standings []Standing `json:"standings"`
}

type Standing struct {
	Rows []Row `json:"rows"`
}

type Row struct {
	Team          Team   `json:"team"`
	Position      int    `json:"position"`
	Matches       int    `json:"matches"`
	Wins          int    `json:"wins"`
	Losses        int    `json:"losses"`
	Draws         int    `json:"draws"`
	Points        int    `json:"points"`
	ScoresFor     int    `json:"scoresFor"`
	ScoresAgainst int    `json:"scoresAgainst"`
	Diff          string `json:"scoreDiffFormatted"`
}
