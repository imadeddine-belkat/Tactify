package services

import (
	"context"
	"fmt"
	"log"

	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	sofascore "github.com/imadeddine-belkat/tactify-protos/go/sofascore/v1"
	"golang.org/x/sync/errgroup"
)

type TeamMatchStatsService struct {
	Event    *EventsService
	Config   config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

const MaxConcurrentMatches = 20

func (s *TeamMatchStatsService) UpdateLeagueMatchStats(ctx context.Context, seasonId int, leagueId int, round int) error {
	log.Printf("Starting UpdateLeagueMatchStats for league=%d, season=%d, round=%d", leagueId, seasonId, round)

	roundMatches, err := s.Event.GetRoundMatches(ctx, seasonId, leagueId, round)
	if err != nil {
		return fmt.Errorf("getting round matches: %w", err)
	}

	if err = s.publishTeamMatchStats(ctx, leagueId, seasonId, round, roundMatches); err != nil {
		return fmt.Errorf("publishing team match stats: %w", err)
	}

	log.Printf("Successfully updated team match stats for league %d, season %d, round %d", leagueId, seasonId, round)

	return nil
}

func (s *TeamMatchStatsService) GetTeamMatchStats(ctx context.Context, matchId int) (*sofascore.MatchStats, error) {
	teamMatchStats := &sofascore.MatchStats{}

	matchStats := s.Config.SofascoreApi.TeamEndpoints.TeamMatchStats
	endpoint := fmt.Sprintf(matchStats, matchId)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, teamMatchStats); err != nil {
		return nil, err
	}

	return teamMatchStats, nil
}

func (s *TeamMatchStatsService) publishTeamMatchStats(ctx context.Context, leagueId, seasonId, round int, roundMatches []*sofascore.Event) error {
	teamMatchStatsTopic := s.Config.KafkaConfig.TopicsName.SofascoreTeamMatchStats.Name

	if len(roundMatches) == 0 {
		log.Printf("No matches found for league %d, season %d, round %d", leagueId, seasonId, round)
		return nil
	}

	matchCount := len(roundMatches)
	log.Printf("Starting processing for %d matches (concurrency limit: %d)", matchCount, MaxConcurrentMatches)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(MaxConcurrentMatches)

	for _, match := range roundMatches {
		match := match

		g.Go(func() error {
			teamMatchStats, err := s.GetTeamMatchStats(ctx, int(match.Id))
			if err != nil {
				log.Printf("Error fetching team match stats for match %d: %v", match.Id, err)
				return nil
			}

			for _, period := range teamMatchStats.Statistics {
				for _, group := range period.Groups {
					for _, stat := range group.StatisticsItems {

						if ctx.Err() != nil {
							return ctx.Err()
						}

						matchStatMsg := &sofascore.MatchStatsMessage{
							SeasonId:   int32(seasonId),
							LeagueId:   int32(leagueId),
							MatchId:    match.Id,
							EventId:    int32(round),
							HomeTeamId: match.HomeTeam.Id,
							AwayTeamId: match.AwayTeam.Id,
							GroupName:  group.GroupName,
							Statistics: &sofascore.StatsItem{
								Period:         period.Period,
								Key:            stat.Key,
								Name:           stat.Name,
								HomeValue:      stat.HomeValue,
								AwayValue:      stat.AwayValue,
								HomeTotal:      stat.HomeTotal,
								AwayTotal:      stat.AwayTotal,
								CompareCode:    stat.CompareCode,
								StatisticsType: stat.StatisticsType,
								RenderType:     stat.RenderType,
							},
						}

						key := []byte(fmt.Sprintf("%d-%s-%s", match.Id, period.Period, stat.Name))

						if err := s.Producer.PublishWithProcess(ctx, matchStatMsg, teamMatchStatsTopic, key); err != nil {
							log.Printf("Error publishing stat '%s' for match %d: %v", stat.Name, match.Id, err)
							continue
						}
					}
				}
			}

			log.Printf("Successfully processed team match stats for match %d", match.Id)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error in match stats processing group: %w", err)
	}

	log.Printf("All workers completed processing for league %d, season %d, round %d", leagueId, seasonId, round)
	return nil
}
