package services

import fpl "github.com/imadeddine-belkat/tactify-protos/go/fpl/v1"

// Services aggregates all service interfaces for FPL data
type Services struct {
	ElementType       ElementTypeService
	Entry             EntryService
	EntryEvent        EntryEventService
	EntryHistory      EntryHistoryService
	Fixture           FixtureService
	LiveEvent         LiveEventService
	Player            PlayerService
	PlayerHistory     PlayerHistoryService
	PlayerPastHistory PlayerPastHistoryService
	Scoring           ScoringService
	Team              TeamService
	Gameweek          GameweekService
}

// TeamService handles Premier League team data
type TeamService interface {
	UpdateTeams(teams []fpl.Team) error
	GetTeam(id int) (*fpl.Team, error)
	GetAllTeams() ([]fpl.Team, error)
	GetTeamsByStrength(minStrength int) ([]fpl.Team, error)
	GetLeagueTable() ([]fpl.Team, error)
	GetTeamFixtures(teamID int, upcoming bool) ([]fpl.Fixture, error)
	GetTeamForm(teamID int, gameweeks int) ([]fpl.Team, error)
}

// PlayerService handles player data and statistics
type PlayerService interface {
	GetPlayer(id int) (*fpl.Player, error)
	GetAllPlayers() ([]fpl.Player, error)
	GetPlayersByTeam(teamID int) ([]fpl.Player, error)
	GetPlayersByPosition(elementType int) ([]fpl.Player, error)
	GetPlayersByPriceRange(minPrice, maxPrice int) ([]fpl.Player, error)
	GetAvailablePlayers() ([]fpl.Player, error)
	GetPlayersByOwnership(minPercent float64) ([]fpl.Player, error)
	GetTopScoringPlayers(position int, gameweeks int) ([]fpl.Player, error)
	SearchPlayersByName(name string) ([]fpl.Player, error)
	GetPlayersByForm(minForm float64) ([]fpl.Player, error)
}

// FixtureService handles match fixtures and results
type FixtureService interface {
	UpdateFixtures(fixtures []fpl.Fixture) error
	GetFixture(id int) (*fpl.Fixture, error)
	GetFixturesByGameweek(gameweek int) ([]fpl.Fixture, error)
	GetFixturesByTeam(teamID int) ([]fpl.Fixture, error)
	GetUpcomingFixtures(limit int) ([]fpl.Fixture, error)
	GetCompletedFixtures(gameweek int) ([]fpl.Fixture, error)
	GetCurrentGameweekFixtures() ([]fpl.Fixture, error)
}

// LiveEventService handles live gameweek data
type LiveEventService interface {
	GetLiveEvent(gameweek int) (*fpl.LiveEvent, error)
	GetLivePlayerStats(playerID, gameweek int) (*fpl.LiveElement, error)
}

// PlayerHistoryService handles player gameweek performance
type PlayerHistoryService interface {
	GetPlayerHistory(playerID int) ([]fpl.PlayerHistory, error)
	GetPlayerFixtureHistory(playerID, fixtureID int) (*fpl.PlayerHistory, error)
	GetPlayerGameweekHistory(playerID, gameweek int) ([]fpl.PlayerHistory, error)
	GetPlayerHomeAwayStats(playerID int) ([]fpl.PlayerHistory, []fpl.PlayerHistory, error)
	GetPlayerRecentForm(playerID int, gameweeks int) ([]fpl.PlayerHistory, error)
}

// PlayerPastHistoryService handles player season history
type PlayerPastHistoryService interface {
	GetPlayerPastHistory(playerID int) ([]fpl.PlayerPastHistory, error)
	GetPlayerSeasonHistory(playerID int, season string) (*fpl.PlayerPastHistory, error)
	GetPlayerCareerStats(playerID int) ([]fpl.PlayerPastHistory, error)
}

// ScoringService handles FPL scoring rules and calculations
type ScoringService interface {
	GetScoringRules() (*fpl.Scoring, error)
	CalculatePlayerGameweekPoints(playerHistory *fpl.PlayerHistory) (int, error)
	CalculateBonusPoints(fixtureID int) (map[int]int, error)
}

// GameweekService handles current gameweek information
type GameweekService interface {
	GetCurrentGameweek() (int, error)
	GetGameweekDeadline(gameweek int) (string, error)
	GetGameweekStatus(gameweek int) (string, error)
	IsGameweekLive(gameweek int) (bool, error)
}

// ElementTypeService handles operations for ElementType (positions: GKP, DEF, MID, FWD)
type ElementTypeService interface {
	GetElementType(id int) (*fpl.ElementType, error)
	GetAllElementTypes() ([]fpl.ElementType, error)
}

// EntryService handles operations for manager entries/teams
type EntryService interface {
	GetEntry(id int) (*fpl.Entry, error)
	GetEntriesByLeague(leagueID int) ([]fpl.Entry, error)
	UpdateEntry(entry *fpl.Entry) error
	GetEntryRankings(limit int) ([]fpl.Entry, error)
}

// EntryEventService handles gameweek-specific team data
type EntryEventService interface {
	GetEntryEvent(entryID, eventID int) (*fpl.EntryEventPicks, error)
	GetEntryPicks(entryID, eventID int) ([]fpl.Pick, error)
	GetAutomaticSubs(entryID, eventID int) ([]fpl.AutomaticSub, error)
}

// EntryHistoryService handles manager historical performance
type EntryHistoryService interface {
	GetEntryHistory(entryID int) (*fpl.EntryHistory, error)
	GetEntryCurrentHistory(entryID int) ([]fpl.EntryHistoryCurrent, error)
	GetEntryPastHistory(entryID int) ([]fpl.EntryHistoryPast, error)
	GetEntryChips(entryID int) ([]fpl.EntryHistoryChip, error)
}
