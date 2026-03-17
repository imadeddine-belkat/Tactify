package sofascore_models

// Team match statistics

type MatchStats struct {
	MatchPeriods []MatchPeriods `json:"statistics"`
}

type MatchPeriods struct {
	Period string   `json:"period"` // first half -> 1, second half -> 2
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
	CompareCode    int     `json:"compareCode"` // Which team performed better: 1-home, 2-away, 3-equal, 0-unknown
	HomeTotal      *int    `json:"homeTotal,omitempty"`
	AwayTotal      *int    `json:"awayTotal,omitempty"`
	StatisticsType string  `json:"statisticsType"` // "positive" higher is better, "negative" higher is worse
	RenderType     int     `json:"renderType"`     // 1=raw numbers, 2=percentage bar, 3=fraction with %, 4=% only from fraction
}
