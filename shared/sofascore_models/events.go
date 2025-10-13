package sofascore_models

type Events struct {
	Events      []Event `json:"events,omitempty"`
	HasNextPage bool    `json:"hasNextPage,omitempty"`
}

type Event struct {
	ID             int        `json:"id,omitempty"`
	Tournament     Tournament `json:"tournament,omitempty"`
	Season         Season     `json:"season,omitempty"`
	RoundInfo      RoundInfo  `json:"roundInfo,omitempty"`
	HomeTeam       Team       `json:"homeTeam,omitempty"`
	AwayTeam       Team       `json:"awayTeam,omitempty"`
	Status         Status     `json:"status,omitempty"`
	StartTimestamp int64      `json:"startTimestamp,omitempty"`
}

type Tournament struct {
	ID               int              `json:"id,omitempty"`
	Name             string           `json:"name,omitempty"`
	UniqueTournament UniqueTournament `json:"uniqueTournament,omitempty"`
}

type UniqueTournament struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Season struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Year string `json:"year,omitempty"`
}

type RoundInfo struct {
	Round int `json:"round,omitempty"`
}

type Team struct {
	ID        int        `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	ShortName string     `json:"shortName,omitempty"`
	NameCode  string     `json:"nameCode,omitempty"`
	Country   Country    `json:"country,omitempty"`
	Colors    TeamColors `json:"teamColors,omitempty"`
}

type Status struct {
	Code        int    `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

type Country struct {
	Name string `json:"name,omitempty"`
}

type TeamColors struct {
	PrimaryColor   string `json:"primary,omitempty"`
	SecondaryColor string `json:"secondary,omitempty"`
	TextColor      string `json:"text,omitempty"`
}
