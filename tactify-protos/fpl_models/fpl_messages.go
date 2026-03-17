package fpl_models

// API message wrappers for event-driven communication

type TeamMessage struct {
	Team     Team `json:"team"`
	SeasonID int  `json:"season_id"`
}

type PlayerBootstrapMessage struct {
	Player   PlayerBootstrap `json:"player"`
	SeasonID int             `json:"season_id"`
}

type PlayerHistoryMessage struct {
	PlayerID int             `json:"player_id"`
	SeasonID int             `json:"season_id"`
	History  []PlayerHistory `json:"history"`
}

type PlayerPastHistoryMessage struct {
	PlayerCode        int                 `json:"element_code" db:"player_code"`
	PlayerPastHistory []PlayerPastHistory `json:"past_history" db:"past_history"`
}

type FixtureMessage struct {
	Fixture  Fixture `json:"fixture"`
	SeasonID int     `json:"season_id"`
}

type LiveEventMessage struct {
	PlayerID int           `json:"player_id"`
	Event    int           `json:"event"`
	SeasonID int           `json:"season_id"`
	Stats    LiveStats     `json:"stats"`
	Explain  []ExplainItem `json:"explain"`
	Modified bool          `json:"modified" db:"modified"`
}

type EntryMessage struct {
	Entry    Entry `json:"entry"`
	SeasonId int   `json:"season_id" db:"season_id"`
}

type EntryEventPicksMessage struct {
	EntryId  int             `json:"entry_id" db:"manager_id"`
	EventId  int             `json:"event" db:"event"`
	SeasonId int             `json:"season_id" db:"season_id"`
	Picks    EntryEventPicks `json:"picks"`
}

type EntryHistoryMessage struct {
	EntryHistory EntryHistory `json:"entry_history"`
	EntryId      int          `json:"entry_id" db:"manager_id"`
	SeasonId     int          `json:"season_id" db:"season_id"`
}

type EntryTransfersMessage struct {
	EntryId   int        `json:"entry_id" db:"manager_id"`
	SeasonId  int        `json:"season_id" db:"season_id"`
	Transfers []Transfer `json:"transfers"`
}
