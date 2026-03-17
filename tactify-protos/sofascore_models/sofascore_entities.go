package sofascore_models

// Core domain entities
type Team struct {
	ID         int        `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"`
	ShortName  string     `json:"shortName,omitempty"`
	NameCode   string     `json:"nameCode,omitempty"`
	Country    Country    `json:"country,omitempty"`
	Colors     TeamColors `json:"teamColors,omitempty"`
	Tournament Tournament `json:"tournament,omitempty"`
}

type Player struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	ShortName        string   `json:"shortName"`
	Position         string   `json:"position"`
	PositionDetailed []string `json:"positionsDetailed"`
	JerseyNumber     string   `json:"jerseyNumber"`
	Height           int      `json:"height"`
	PreferredFoot    string   `json:"preferredFoot"`
	Country          Country  `json:"country"`
	Team             Team     `json:"team"`
}

type LineupPlayer struct {
	ID                     int     `json:"id"`
	Name                   string  `json:"name"`
	FirstName              string  `json:"firstName,omitempty"`
	LastName               string  `json:"lastName,omitempty"`
	Slug                   string  `json:"slug"`
	ShortName              string  `json:"shortName"`
	Position               string  `json:"position"`
	JerseyNumber           string  `json:"jerseyNumber"`
	Height                 int     `json:"height"`
	UserCount              int     `json:"userCount"`
	Gender                 string  `json:"gender"`
	SofascoreID            string  `json:"sofascoreId,omitempty"`
	Country                Country `json:"country"`
	MarketValueCurrency    string  `json:"marketValueCurrency,omitempty"`
	DateOfBirthTimestamp   int64   `json:"dateOfBirthTimestamp"`
	ProposedMarketValueRaw *struct {
		Value    int    `json:"value"`
		Currency string `json:"currency"`
	} `json:"proposedMarketValueRaw,omitempty"`
	FieldTranslations *struct {
		NameTranslation      map[string]string `json:"nameTranslation,omitempty"`
		ShortNameTranslation map[string]string `json:"shortNameTranslation,omitempty"`
	} `json:"fieldTranslations,omitempty"`
}

type Tournament struct {
	Name             string           `json:"name,omitempty"`
	UniqueTournament UniqueTournament `json:"uniqueTournament,omitempty"`
}

type UniqueTournament struct {
	ID       int            `json:"id,omitempty"`
	Name     string         `json:"name,omitempty"`
	Category LeagueCategory `json:"category,omitempty"`
}

type Season struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Year      string `json:"year,omitempty"`
	IsCurrent bool   `json:"is_current"`
}

type Seasons struct {
	LeagueID int      `json:"LeagueID,omitempty"`
	Seasons  []Season `json:"seasons,omitempty"`
}

type LeagueCategory struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	Flag string `json:"flag"`
}

type RoundInfo struct {
	Round int `json:"round,omitempty"`
}
