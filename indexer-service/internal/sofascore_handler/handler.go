package sofascore_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/imadeddine-belkat/indexer-service/config"
	"github.com/imadeddine-belkat/indexer-service/internal/sofascore_repositories"
	"github.com/imadeddine-belkat/kafka"
	kafkaConfig "github.com/imadeddine-belkat/kafka/config"
	"github.com/imadeddine-belkat/shared/sofascore_models"
)

type Handler struct {
	config      *config.IndexerConfig
	kafkaConfig *kafkaConfig.KafkaConfig
	consumers   map[string]*kafka.Consumer
	teamRepo    *sofascore_repositories.TeamRepo
	matchRepo   *sofascore_repositories.MatchRepo
	leagueRepo  *sofascore_repositories.LeagueRepo
}

func NewHandler(
	cfg *config.IndexerConfig,
	kafkaCfg *kafkaConfig.KafkaConfig,
	teamRepo *sofascore_repositories.TeamRepo,
	matchRepo *sofascore_repositories.MatchRepo,
	leagueRepo *sofascore_repositories.LeagueRepo,
) *Handler {
	h := &Handler{
		config:      cfg,
		kafkaConfig: kafkaCfg,
		consumers:   make(map[string]*kafka.Consumer),
		teamRepo:    teamRepo,
		matchRepo:   matchRepo,
		leagueRepo:  leagueRepo,
	}

	if teamRepo != nil {
		h.consumers[kafkaCfg.TopicsName.SofascoreLeagueStandings.Name] = kafka.NewConsumer(
			kafkaCfg,
			kafkaCfg.TopicsName.SofascoreLeagueStandings.Name,
			kafkaCfg.ConsumersGroupID.SofascoreLeagueStanding,
		)

		h.consumers[kafkaCfg.TopicsName.SofascoreTeamOverallStats.Name] = kafka.NewConsumer(
			kafkaCfg,
			kafkaCfg.TopicsName.SofascoreTeamOverallStats.Name,
			kafkaCfg.ConsumersGroupID.SofascoreTeamOverallStats,
		)

		h.consumers[kafkaCfg.TopicsName.SofascoreTeamMatchStats.Name] = kafka.NewConsumer(
			kafkaCfg,
			kafkaCfg.TopicsName.SofascoreTeamMatchStats.Name,
			kafkaCfg.ConsumersGroupID.SofascoreTeamMatchStats,
		)
	}

	if matchRepo != nil {
		h.consumers[kafkaCfg.TopicsName.SofascoreLeagueRoundMatches.Name] = kafka.NewConsumer(
			kafkaCfg,
			kafkaCfg.TopicsName.SofascoreLeagueRoundMatches.Name,
			kafkaCfg.ConsumersGroupID.SofascoreLeagueRoundMatches,
		)
	}

	if leagueRepo != nil {
		h.consumers[kafkaCfg.TopicsName.SofascoreLeagueIDs.Name] = kafka.NewConsumer(
			kafkaCfg,
			kafkaCfg.TopicsName.SofascoreLeagueIDs.Name,
			kafkaCfg.ConsumersGroupID.SofascoreLeagueRoundMatches,
		)

		h.consumers[kafkaCfg.TopicsName.SofascoreLeagueSeasons.Name] = kafka.NewConsumer(
			kafkaCfg,
			kafkaCfg.TopicsName.SofascoreLeagueSeasons.Name,
			kafkaCfg.ConsumersGroupID.SofascoreLeagueRoundMatches,
		)
	}
	return h
}

type HandlerFunc func(ctx context.Context)

func (h *Handler) Route(ctx context.Context, topic string) {

	handlers := map[string]HandlerFunc{
		h.kafkaConfig.TopicsName.SofascoreLeagueStandings.Name:    h.handleTeamsInfo,
		h.kafkaConfig.TopicsName.SofascoreTeamOverallStats.Name:   h.handleTeamOverallStats,
		h.kafkaConfig.TopicsName.SofascoreTeamMatchStats.Name:     h.handleTeamMatchStat,
		h.kafkaConfig.TopicsName.SofascoreLeagueRoundMatches.Name: h.handleLeagueRoundMatches,
		h.kafkaConfig.TopicsName.SofascoreLeagueIDs.Name:          h.handleLeagueInfo,
		h.kafkaConfig.TopicsName.SofascoreLeagueSeasons.Name:      h.handleLeagueSeasonsInfo,
	}

	if fn, ok := handlers[topic]; ok {
		go fn(ctx)
	} else {
		log.Printf("[WARN] No handler found for topic=%s\n", topic)
	}
}

// Generic batch processor - handles all the common batching logic
func batchProcess[T any, K comparable](
	ctx context.Context,
	consumer *kafka.Consumer,
	batchSize int,
	flushInterval time.Duration,
	topicName string,
	getKey func(T) K,
	process func(T) error,
) {

	messages, errors := consumer.Subscribe(ctx)
	batch := make(map[K]T)
	flushTicker := time.NewTicker(flushInterval)
	defer flushTicker.Stop()

	flushBatch := func(logContext string) {
		if len(batch) == 0 {
			return
		}

		for _, item := range batch {
			if err := process(item); err != nil {
				log.Printf("Error %s for %s: %v\n", logContext, topicName, err)
			}
		}
		batch = make(map[K]T)
	}

	for {
		select {
		case msg := <-messages:

			var item T
			if err := json.Unmarshal(msg.Value, &item); err != nil {
				continue
			}

			batch[getKey(item)] = item

			if len(batch) >= batchSize {
				flushBatch("inserting")
			}

		case <-flushTicker.C:
			flushBatch("flushing")

		case err := <-errors:
			if err != nil {
				log.Printf("Error consuming %s message: %v\n", topicName, err)
			}

		case <-ctx.Done():
			flushBatch("inserting on shutdown")
			return
		}
	}
}

func (h *Handler) handleTeamsInfo(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.SofascoreLeagueStandings.Name],
		1,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.SofascoreLeagueStandings.Name,
		func(t sofascore_models.StandingMessage) string {
			return fmt.Sprintf("%d", time.Now().UnixNano())
		},
		h.teamRepo.InsertTeamInfo,
	)
}

func (h *Handler) handleTeamOverallStats(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.SofascoreTeamOverallStats.Name],
		1,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.SofascoreTeamOverallStats.Name,
		func(t sofascore_models.TeamOverallStatsMessage) string {
			return fmt.Sprintf("%d", time.Now().UnixNano())
		},
		h.teamRepo.InsertTeamOverallStats,
	)
}

func (h *Handler) handleTeamMatchStat(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.SofascoreTeamMatchStats.Name],
		1,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.SofascoreTeamMatchStats.Name,
		func(t sofascore_models.MatchStatsMessage) string {
			return fmt.Sprintf("%d", time.Now().UnixNano())
		},
		h.teamRepo.InsertTeamMatchStats,
	)
}

func (h *Handler) handleLeagueRoundMatches(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.SofascoreLeagueRoundMatches.Name],
		1,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.SofascoreLeagueRoundMatches.Name,
		func(t sofascore_models.Event) string {
			return fmt.Sprintf("%d", time.Now().UnixNano())
		},
		h.matchRepo.InsertLeagueRoundMatches,
	)
}

func (h *Handler) handleLeagueInfo(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.SofascoreLeagueIDs.Name],
		1,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.SofascoreLeagueIDs.Name,
		func(l sofascore_models.LeagueUniqueTournaments) string {
			return fmt.Sprintf("%d", time.Now().UnixNano())
		},
		h.leagueRepo.InsertLeagueInfo,
	)
}

func (h *Handler) handleLeagueSeasonsInfo(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.SofascoreLeagueSeasons.Name],
		1,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.SofascoreLeagueSeasons.Name,
		func(l sofascore_models.Seasons) string {
			return fmt.Sprintf("%d", time.Now().UnixNano())
		},
		h.leagueRepo.InsertLeagueSeasonsInfo,
	)
}
