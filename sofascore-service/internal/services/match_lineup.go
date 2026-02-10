package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	sofascore "github.com/imadeddine-belkat/tactify-protos/go/sofascore/v1"
	"golang.org/x/sync/errgroup"
)

type MatchLineupService struct {
	Event    *EventsService
	Config   *config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

func (l *MatchLineupService) GetMatchLineup(ctx context.Context, matchID int) (*sofascore.MatchLineup, error) {
	matchLineup := &sofascore.MatchLineup{}

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

	for _, event := range roundMatches {
		event := event
		g.Go(func() error {
			matchLineup, err := l.GetMatchLineup(ctx, int(event.Id))
			if err != nil {
				if isNotFoundError(err) {
					fmt.Printf("Warning: Lineup not found for match %d (this is normal for future/cancelled matches)\n", event.Id)
					return nil
				}
				return fmt.Errorf("failed to get lineup for match %d: %w", event.Id, err)
			}

			players, err := l.processMatchLineup(matchLineup, seasonId, leagueId, round, int(event.Id))
			if err != nil {
				return fmt.Errorf("failed to process match lineup for match %d: %w", event.Id, err)
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

func (l *MatchLineupService) processMatchLineup(lineup *sofascore.MatchLineup, seasonId, leagueId, round, eventId int) ([]*sofascore.PlayerMatchStatsMessage, error) {
	players := make([]*sofascore.PlayerMatchStatsMessage, 0, len(lineup.Home.Players)+len(lineup.Away.Players))

	for _, player := range lineup.Home.Players {
		players = append(players, &sofascore.PlayerMatchStatsMessage{
			PlayerName: player.Player.Name,
			SeasonId:   int32(seasonId),
			LeagueId:   int32(leagueId),
			MatchId:    int32(eventId),
			Round:      int32(round),
			Player:     player,
		})
	}

	for _, player := range lineup.Away.Players {
		players = append(players, &sofascore.PlayerMatchStatsMessage{
			PlayerName: player.Player.Name,
			SeasonId:   int32(seasonId),
			LeagueId:   int32(leagueId),
			MatchId:    int32(eventId),
			Round:      int32(round),
			Player:     player,
		})
	}

	return players, nil
}

func (l *MatchLineupService) publishPlayer(ctx context.Context, player *sofascore.PlayerMatchStatsMessage) error {
	playersStatsTopic := l.Config.KafkaConfig.TopicsName.SofascorePlayerMatchStats.Name

	key := []byte(fmt.Sprintf("%s-%d-%d-%d",
		player.Player.Player.Name,
		player.SeasonId,
		player.LeagueId,
		player.Round,
	))

	return l.Producer.PublishWithProcess(ctx, player, playersStatsTopic, key)
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
