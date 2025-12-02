package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/shared/sofascore_models"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	"golang.org/x/sync/errgroup"
)

type MatchLineupService struct {
	Event    *EventsService
	Config   *config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

func (l *MatchLineupService) GetMatchLineup(ctx context.Context, matchID int) (*sofascore_models.MatchLineup, error) {
	matchLineup := &sofascore_models.MatchLineup{}

	lineup := l.Config.SofascoreApi.MatchEndpoints.MatchLineups //event/%d/lineups
	endpoint := fmt.Sprintf(lineup, matchID)

	if err := l.Client.GetAndUnmarshal(ctx, endpoint, matchLineup); err != nil {
		return nil, err
	}

	return matchLineup, nil
}

func (l *MatchLineupService) UpdatePlayersStats(ctx context.Context, seasonId, leagueId, round int) error {
	g, ctx := errgroup.WithContext(ctx)

	roundMatches, err := l.Event.GetRoundMatches(ctx, seasonId, leagueId, round)
	if err != nil {
		return fmt.Errorf("failed to get round matches for season %d, league %d, round %d: %w", seasonId, leagueId, round, err)
	}

	for _, event := range roundMatches.Events {
		event := event
		g.Go(func() error {
			matchLineup, err := l.GetMatchLineup(ctx, event.ID)
			if err != nil {
				if isNotFoundError(err) {
					// Log but don't fail - some matches may not have lineups yet
					fmt.Printf("Warning: Lineup not found for match %d (this is normal for future/cancelled matches)\n", event.ID)
					return nil // Continue processing other matches
				}
				return fmt.Errorf("failed to get lineup for match %d: %w", event.ID, err)
			}

			players, err := l.processMatchLineup(matchLineup, seasonId, leagueId, round, event.ID)
			if err != nil {
				return fmt.Errorf("failed to process match lineup for match %d: %w", event.ID, err)
			}

			publishGroup, publishCtx := errgroup.WithContext(ctx)

			for _, player := range players {
				player := player
				publishGroup.Go(func() error {
					return l.publishPlayer(publishCtx, player)
				})
			}
			return publishGroup.Wait()
		})
	}

	return g.Wait()
}

func (l *MatchLineupService) processMatchLineup(lineup *sofascore_models.MatchLineup, seasonId, leagueId, round, eventId int) ([]*sofascore_models.PlayerMatchStatsMessage, error) {
	players := make([]*sofascore_models.PlayerMatchStatsMessage, 0, len(lineup.Home.Players)+len(lineup.Away.Players))

	for _, player := range lineup.Home.Players {
		players = append(players, &sofascore_models.PlayerMatchStatsMessage{
			PlayerName:  player.Player.Name,
			SeasonID:    seasonId,
			LeagueID:    leagueId,
			MatchID:     eventId,
			Round:       round,
			MatchPlayer: player,
		})
	}

	for _, player := range lineup.Away.Players {
		players = append(players, &sofascore_models.PlayerMatchStatsMessage{
			PlayerName:  player.Player.Name,
			SeasonID:    seasonId,
			LeagueID:    leagueId,
			MatchID:     eventId,
			Round:       round,
			MatchPlayer: player,
		})
	}

	return players, nil
}

func (l *MatchLineupService) publishPlayer(ctx context.Context, player *sofascore_models.PlayerMatchStatsMessage) error {
	playersStatsTopic := l.Config.KafkaConfig.TopicsName.SofascorePlayerMatchStats

	value, err := json.Marshal(player)
	if err != nil {
		return fmt.Errorf("failed to marshal player for match %d: %w", player.MatchID, err)
	}

	key := []byte(fmt.Sprintf("%s-%d-%d-%d",
		player.MatchPlayer.Player.Name,
		player.SeasonID,
		player.LeagueID,
		player.Round))

	return l.Producer.Publish(ctx, playersStatsTopic, key, value)
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()

	// Check for the exact error format: {"error":{"code":404,"message":"Not Found"}}
	return strings.Contains(errMsg, `"code":404`) ||
		strings.Contains(errMsg, "404") && strings.Contains(errMsg, "Not Found")
}
