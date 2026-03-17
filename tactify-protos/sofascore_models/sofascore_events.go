package sofascore_models

// Event (Match) related types

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
