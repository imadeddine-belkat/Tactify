package models

type MatchStats struct {
	Event           int               `json:"event"`
	MatchID         int               `json:"match_id"`
	HomeTeamID      int               `json:"homeTeamId"`
	AwayTeamID      int               `json:"awayTeamId"`
	MatchStatistics []MatchStatistics `json:"statistics"`
}

type MatchStatistics struct {
	Period string   `json:"period"`
	Groups []Groups `json:"groups"`
}

type Groups struct {
	GroupName  string      `json:"groupName"`
	StatsItems []StatsItem `json:"statisticsItems"`
}

type StatsItem struct {
	Key            string  `json:"key"`
	Name           string  `json:"name"`
	HomeValue      float64 `json:"homeValue"`
	AwayValue      float64 `json:"awayValue"`
	CompareCode    int     `json:"compareCode"` //Which team performed better 1-home, 2-away, 3-equal 0-unknown
	HomeTotal      *int    `json:"homeTotal,omitempty"`
	AwayTotal      *int    `json:"awayTotal,omitempty"`
	StatisticsType string  `json:"statisticsType"` //"positive" higher is better,"negative" higher is worse
	RenderType     int     `json:"renderType"`     // renderType: 1 = raw numbers (e.g., "11" vs "11"), 2 = percentage bar (e.g., "54%" vs "46%"), 3 = fraction with percentage (e.g., "32/78 (41%)"), 4 = percentage only from fraction (e.g., "42%" calculated from 8/19)
}
