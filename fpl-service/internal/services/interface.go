package services

import (
	"github.com/imadbelkat1/shared/models"
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
	UpdateTeams(teams []models.Team) error
	GetTeam(id int) (*models.Team, error)
	GetAllTeams() ([]models.Team, error)
	GetTeamsByStrength(minStrength int) ([]models.Team, error)
	GetLeagueTable() ([]models.Team, error)
	GetTeamFixtures(teamID int, upcoming bool) ([]models.Fixture, error)
	GetTeamForm(teamID int, gameweeks int) ([]models.Team, error)
}

// PlayerService handles player data and statistics
type PlayerService interface {
	GetPlayer(id int) (*models.Player, error)
	GetAllPlayers() ([]models.Player, error)
	GetPlayersByTeam(teamID int) ([]models.Player, error)
	GetPlayersByPosition(elementType int) ([]models.Player, error)
	GetPlayersByPriceRange(minPrice, maxPrice int) ([]models.Player, error)
	GetAvailablePlayers() ([]models.Player, error)
	GetPlayersByOwnership(minPercent float64) ([]models.Player, error)
	GetTopScoringPlayers(position int, gameweeks int) ([]models.Player, error)
	SearchPlayersByName(name string) ([]models.Player, error)
	GetPlayersByForm(minForm float64) ([]models.Player, error)
}

// FixtureService handles match fixtures and results
type FixtureService interface {
	UpdateFixtures(fixtures []models.Fixture) error
	GetFixture(id int) (*models.Fixture, error)
	GetFixturesByGameweek(gameweek int) ([]models.Fixture, error)
	GetFixturesByTeam(teamID int) ([]models.Fixture, error)
	GetUpcomingFixtures(limit int) ([]models.Fixture, error)
	GetCompletedFixtures(gameweek int) ([]models.Fixture, error)
	GetCurrentGameweekFixtures() ([]models.Fixture, error)
}

// LiveEventService handles live gameweek data
type LiveEventService interface {
	GetLiveEvent(gameweek int) (*models.LiveEvent, error)
	GetLivePlayerStats(playerID, gameweek int) (*models.LiveElement, error)
}

// PlayerHistoryService handles player gameweek performance
type PlayerHistoryService interface {
	GetPlayerHistory(playerID int) ([]models.PlayerHistory, error)
	GetPlayerFixtureHistory(playerID, fixtureID int) (*models.PlayerHistory, error)
	GetPlayerGameweekHistory(playerID, gameweek int) ([]models.PlayerHistory, error)
	GetPlayerHomeAwayStats(playerID int) ([]models.PlayerHistory, []models.PlayerHistory, error)
	GetPlayerRecentForm(playerID int, gameweeks int) ([]models.PlayerHistory, error)
}

// PlayerPastHistoryService handles player season history
type PlayerPastHistoryService interface {
	GetPlayerPastHistory(playerID int) ([]models.PlayerPastHistory, error)
	GetPlayerSeasonHistory(playerID int, season string) (*models.PlayerPastHistory, error)
	GetPlayerCareerStats(playerID int) ([]models.PlayerPastHistory, error)
}

// ScoringService handles FPL scoring rules and calculations
type ScoringService interface {
	GetScoringRules() (*models.Scoring, error)
	CalculatePlayerGameweekPoints(playerHistory *models.PlayerHistory) (int, error)
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
	GetElementType(id int) (*models.ElementType, error)
	GetAllElementTypes() ([]models.ElementType, error)
}

// EntryService handles operations for manager entries/teams
type EntryService interface {
	GetEntry(id int) (*models.Entry, error)
	GetEntriesByLeague(leagueID int) ([]models.Entry, error)
	UpdateEntry(entry *models.Entry) error
	GetEntryRankings(limit int) ([]models.Entry, error)
}

// EntryEventService handles gameweek-specific team data
type EntryEventService interface {
	GetEntryEvent(entryID, eventID int) (*models.EntryEventPicks, error)
	GetEntryPicks(entryID, eventID int) ([]models.Pick, error)
	GetAutomaticSubs(entryID, eventID int) ([]models.AutomaticSub, error)
}

// EntryHistoryService handles manager historical performance
type EntryHistoryService interface {
	GetEntryHistory(entryID int) (*models.EntryHistory, error)
	GetEntryCurrentHistory(entryID int) ([]models.EntryHistoryCurrent, error)
	GetEntryPastHistory(entryID int) ([]models.EntryHistoryPast, error)
	GetEntryChips(entryID int) ([]models.EntryHistoryChip, error)
}
