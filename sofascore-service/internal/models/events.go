package models

type Events struct {
	Events      []Event `json:"events"`
	HasNextPage bool    `json:"hasNextPage"`
}

type Event struct {
	ID             int        `json:"id"`
	Tournament     Tournament `json:"tournament"`
	Season         Season     `json:"season"`
	RoundInfo      RoundInfo  `json:"roundInfo"`
	HomeTeam       Team       `json:"homeTeam"`
	AwayTeam       Team       `json:"awayTeam"`
	Status         Status     `json:"status"`
	StartTimestamp int64      `json:"startTimestamp"`
}

type Tournament struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Season struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Year string `json:"year"`
}

type RoundInfo struct {
	Round int `json:"round"`
}

type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	NameCode  string `json:"nameCode"`
}

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Type        string `json:"type"`
}
