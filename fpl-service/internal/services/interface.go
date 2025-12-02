package services

import (
	"github.com/imadeddine-belkat/shared/fpl_models"
)

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
	UpdateTeams(teams []fpl_models.Team) error
	GetTeam(id int) (*fpl_models.Team, error)
	GetAllTeams() ([]fpl_models.Team, error)
	GetTeamsByStrength(minStrength int) ([]fpl_models.Team, error)
	GetLeagueTable() ([]fpl_models.Team, error)
	GetTeamFixtures(teamID int, upcoming bool) ([]fpl_models.Fixture, error)
	GetTeamForm(teamID int, gameweeks int) ([]fpl_models.Team, error)
}

// PlayerService handles player data and statistics
type PlayerService interface {
	GetPlayer(id int) (*fpl_models.Player, error)
	GetAllPlayers() ([]fpl_models.Player, error)
	GetPlayersByTeam(teamID int) ([]fpl_models.Player, error)
	GetPlayersByPosition(elementType int) ([]fpl_models.Player, error)
	GetPlayersByPriceRange(minPrice, maxPrice int) ([]fpl_models.Player, error)
	GetAvailablePlayers() ([]fpl_models.Player, error)
	GetPlayersByOwnership(minPercent float64) ([]fpl_models.Player, error)
	GetTopScoringPlayers(position int, gameweeks int) ([]fpl_models.Player, error)
	SearchPlayersByName(name string) ([]fpl_models.Player, error)
	GetPlayersByForm(minForm float64) ([]fpl_models.Player, error)
}

// FixtureService handles match fixtures and results
type FixtureService interface {
	UpdateFixtures(fixtures []fpl_models.Fixture) error
	GetFixture(id int) (*fpl_models.Fixture, error)
	GetFixturesByGameweek(gameweek int) ([]fpl_models.Fixture, error)
	GetFixturesByTeam(teamID int) ([]fpl_models.Fixture, error)
	GetUpcomingFixtures(limit int) ([]fpl_models.Fixture, error)
	GetCompletedFixtures(gameweek int) ([]fpl_models.Fixture, error)
	GetCurrentGameweekFixtures() ([]fpl_models.Fixture, error)
}

// LiveEventService handles live gameweek data
type LiveEventService interface {
	GetLiveEvent(gameweek int) (*fpl_models.LiveEvent, error)
	GetLivePlayerStats(playerID, gameweek int) (*fpl_models.LiveElement, error)
}

// PlayerHistoryService handles player gameweek performance
type PlayerHistoryService interface {
	GetPlayerHistory(playerID int) ([]fpl_models.PlayerHistory, error)
	GetPlayerFixtureHistory(playerID, fixtureID int) (*fpl_models.PlayerHistory, error)
	GetPlayerGameweekHistory(playerID, gameweek int) ([]fpl_models.PlayerHistory, error)
	GetPlayerHomeAwayStats(playerID int) ([]fpl_models.PlayerHistory, []fpl_models.PlayerHistory, error)
	GetPlayerRecentForm(playerID int, gameweeks int) ([]fpl_models.PlayerHistory, error)
}

// PlayerPastHistoryService handles player season history
type PlayerPastHistoryService interface {
	GetPlayerPastHistory(playerID int) ([]fpl_models.PlayerPastHistory, error)
	GetPlayerSeasonHistory(playerID int, season string) (*fpl_models.PlayerPastHistory, error)
	GetPlayerCareerStats(playerID int) ([]fpl_models.PlayerPastHistory, error)
}

// ScoringService handles FPL scoring rules and calculations
type ScoringService interface {
	GetScoringRules() (*fpl_models.Scoring, error)
	CalculatePlayerGameweekPoints(playerHistory *fpl_models.PlayerHistory) (int, error)
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
	GetElementType(id int) (*fpl_models.ElementType, error)
	GetAllElementTypes() ([]fpl_models.ElementType, error)
}

// EntryService handles operations for manager entries/teams
type EntryService interface {
	GetEntry(id int) (*fpl_models.Entry, error)
	GetEntriesByLeague(leagueID int) ([]fpl_models.Entry, error)
	UpdateEntry(entry *fpl_models.Entry) error
	GetEntryRankings(limit int) ([]fpl_models.Entry, error)
}

// EntryEventService handles gameweek-specific team data
type EntryEventService interface {
	GetEntryEvent(entryID, eventID int) (*fpl_models.EntryEventPicks, error)
	GetEntryPicks(entryID, eventID int) ([]fpl_models.Pick, error)
	GetAutomaticSubs(entryID, eventID int) ([]fpl_models.AutomaticSub, error)
}

// EntryHistoryService handles manager historical performance
type EntryHistoryService interface {
	GetEntryHistory(entryID int) (*fpl_models.EntryHistory, error)
	GetEntryCurrentHistory(entryID int) ([]fpl_models.EntryHistoryCurrent, error)
	GetEntryPastHistory(entryID int) ([]fpl_models.EntryHistoryPast, error)
	GetEntryChips(entryID int) ([]fpl_models.EntryHistoryChip, error)
}
