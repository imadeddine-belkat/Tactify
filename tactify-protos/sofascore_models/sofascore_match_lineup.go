package sofascore_models

// Match lineup and player match data

type MatchLineup struct {
	Confirmed bool       `json:"confirmed"`
	Home      TeamLineup `json:"home"`
	Away      TeamLineup `json:"away"`
}

type TeamLineup struct {
	Players         []MatchPlayer   `json:"players"`
	SupportStaff    []interface{}   `json:"supportStaff"`
	Formation       string          `json:"formation"`
	PlayerColor     PlayerColor     `json:"playerColor"`
	GoalkeeperColor PlayerColor     `json:"goalkeeperColor"`
	MissingPlayers  []MissingPlayer `json:"missingPlayers,omitempty"`
}

type MatchPlayer struct {
	Player       Player           `json:"player"`
	TeamID       int              `json:"teamId"`
	ShirtNumber  int              `json:"shirtNumber"`
	JerseyNumber string           `json:"jerseyNumber"`
	Position     string           `json:"position"`
	Substitute   bool             `json:"substitute"`
	Captain      bool             `json:"captain,omitempty"`
	Statistics   PlayerStatistics `json:"statistics"`
}

type MissingPlayer struct {
	Player          Player `json:"player"`
	Type            string `json:"type"`
	Reason          int    `json:"reason"`
	Description     string `json:"description,omitempty"`
	ExternalType    int    `json:"externalType"`
	ExpectedEndDate string `json:"expectedEndDate,omitempty"`
}
