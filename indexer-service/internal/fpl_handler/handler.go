package fpl_handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/imadeddine-belkat/indexer-service/config"
	"github.com/imadeddine-belkat/indexer-service/internal/fpl_repositories"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	kafkaConfig "github.com/imadeddine-belkat/tactify-kafka/config"
	fpl "github.com/imadeddine-belkat/tactify-protos/go/fpl/v1"
)

type Handler struct {
	config      *config.IndexerConfig
	kafkaConfig *kafkaConfig.KafkaConfig
	consumers   map[string]*kafka.Consumer
	playerRepo  *fpl_repositories.PlayerRepo
	teamRepo    *fpl_repositories.TeamRepo
	fixtureRepo *fpl_repositories.FixtureRepo
	managerRepo *fpl_repositories.ManagerRepo
}

func NewHandler(
	config *config.IndexerConfig,
	kafkaConfig *kafkaConfig.KafkaConfig,
	playerRepo *fpl_repositories.PlayerRepo,
	teamRepo *fpl_repositories.TeamRepo,
	fixtureRepo *fpl_repositories.FixtureRepo,
	managerRepo *fpl_repositories.ManagerRepo,
) *Handler {
	h := &Handler{
		config:      config,
		kafkaConfig: kafkaConfig,
		playerRepo:  playerRepo,
		teamRepo:    teamRepo,
		fixtureRepo: fixtureRepo,
		managerRepo: managerRepo,
		consumers:   make(map[string]*kafka.Consumer),
	}

	// Pre-create consumers only for non-nil repositories
	if fixtureRepo != nil {
		h.consumers[kafkaConfig.TopicsName.FplFixtures.Name] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplFixtures.Name,
			kafkaConfig.ConsumersGroupID.FplFixtures,
		)
	}

	if teamRepo != nil {
		h.consumers[kafkaConfig.TopicsName.FplTeams.Name] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplTeams.Name,
			kafkaConfig.ConsumersGroupID.FplTeams,
		)
	}

	if playerRepo != nil {
		h.consumers[kafkaConfig.TopicsName.FplPlayersBootstrap.Name] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplPlayersBootstrap.Name,
			kafkaConfig.ConsumersGroupID.FplPlayers,
		)

		h.consumers[kafkaConfig.TopicsName.FplPlayersStats.Name] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplPlayersStats.Name,
			kafkaConfig.ConsumersGroupID.FplPlayersStats,
		)

		h.consumers[kafkaConfig.TopicsName.FplPlayerMatchStats.Name] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplPlayerMatchStats.Name,
			kafkaConfig.ConsumersGroupID.FplPlayersStats,
		)

		h.consumers[kafkaConfig.TopicsName.FplPlayerHistoryStats.Name] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplPlayerHistoryStats.Name,
			kafkaConfig.ConsumersGroupID.FplPlayersStats,
		)

		h.consumers[kafkaConfig.TopicsName.FplLiveEvent.Name] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplLiveEvent.Name,
			kafkaConfig.ConsumersGroupID.FplLive,
		)
	}

	if managerRepo != nil {
		h.consumers[kafkaConfig.TopicsName.FplEntry.Name] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplEntry.Name,
			kafkaConfig.ConsumersGroupID.FplEntries,
		)

		h.consumers[kafkaConfig.TopicsName.FplEntryPicks.Name] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplEntryPicks.Name,
			kafkaConfig.ConsumersGroupID.FplEntriesPicks,
		)

		h.consumers[kafkaConfig.TopicsName.FplEntryTransfers.Name] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplEntryTransfers.Name,
			kafkaConfig.ConsumersGroupID.FplEntriesTransfers,
		)

		h.consumers[kafkaConfig.TopicsName.FplEntryHistory.Name] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplEntryHistory.Name,
			kafkaConfig.ConsumersGroupID.FplEntriesHistory,
		)
	}

	return h
}

type HandlerFunc func(ctx context.Context)

func (h *Handler) Route(ctx context.Context, topic string) {
	handlers := map[string]HandlerFunc{
		h.kafkaConfig.TopicsName.FplFixtures.Name:           h.handleFixtures,
		h.kafkaConfig.TopicsName.FplTeams.Name:              h.handleTeams,
		h.kafkaConfig.TopicsName.FplPlayersBootstrap.Name:   h.handlePlayerBootstrap,
		h.kafkaConfig.TopicsName.FplPlayersStats.Name:       h.handlePlayerStats,
		h.kafkaConfig.TopicsName.FplPlayerMatchStats.Name:   h.handlePlayerMatchStats,
		h.kafkaConfig.TopicsName.FplPlayerHistoryStats.Name: h.handlePlayerPastHistory,
		h.kafkaConfig.TopicsName.FplLiveEvent.Name:          h.handlePlayerExplain,
		h.kafkaConfig.TopicsName.FplEntry.Name:              h.handleManagers,
		h.kafkaConfig.TopicsName.FplEntryPicks.Name:         h.handleManagerPicks,
		h.kafkaConfig.TopicsName.FplEntryTransfers.Name:     h.handleManagerTransfers,
		h.kafkaConfig.TopicsName.FplEntryHistory.Name:       h.handleManagerHistory,
	}

	if fn, ok := handlers[topic]; ok {
		go fn(ctx)
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
				log.Printf("Error unmarshaling %s message: %v, raw: %s\n", topicName, err, string(msg.Value))
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

// Generic batch processor with slice conversion - for handlers that need to convert map to slice
func batchProcessWithSlice[T any, K comparable](
	ctx context.Context,
	consumer *kafka.Consumer,
	batchSize int,
	flushInterval time.Duration,
	topicName string,
	getKey func(T) K,
	processBatch func([]T) error,
) {
	messages, errors := consumer.Subscribe(ctx)
	batch := make(map[K]T)
	flushTicker := time.NewTicker(flushInterval)
	defer flushTicker.Stop()

	flushBatch := func(logContext string) {
		if len(batch) == 0 {
			return
		}

		items := make([]T, 0, len(batch))
		for _, item := range batch {
			items = append(items, item)
		}

		if err := processBatch(items); err != nil {
			log.Printf("Error %s batch for %s: %v\n", logContext, topicName, err)
		}
		batch = make(map[K]T)
	}

	for {
		select {
		case msg := <-messages:
			var item T
			if err := json.Unmarshal(msg.Value, &item); err != nil {
				log.Printf("Error unmarshaling %s message: %v, raw: %s\n", topicName, err, string(msg.Value))
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

func (h *Handler) handleFixtures(ctx context.Context) {
	batchProcessWithSlice(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplFixtures.Name],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplFixtures.Name,
		func(f *fpl.FixtureMessage) int { return int(f.Fixture.Id) },
		h.fixtureRepo.InsertFixtures,
	)
}

func (h *Handler) handleTeams(ctx context.Context) {
	batchProcessWithSlice(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplTeams.Name],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplTeams.Name,
		func(t *fpl.TeamMessage) int { return int(t.Team.Id) },
		h.teamRepo.InsertTeams,
	)

}

func (h *Handler) handlePlayerBootstrap(ctx context.Context) {
	consumer := h.consumers[h.kafkaConfig.TopicsName.FplPlayersBootstrap.Name]
	messages, errors := consumer.Subscribe(ctx)
	batch := make(map[int]*fpl.PlayerBootstrapMessage)

	flushTicker := time.NewTicker(h.config.FlushInterval)
	verifyTicker := time.NewTicker(30 * time.Second)
	defer flushTicker.Stop()
	defer verifyTicker.Stop()

	totalReceived := 0
	totalProcessed := 0

	flushBatch := func(logContext string) {
		if len(batch) == 0 {
			return
		}

		players := make([]*fpl.PlayerBootstrapMessage, 0, len(batch))
		for _, p := range batch {
			players = append(players, p)
		}

		if err := h.playerRepo.InsertPlayerBootstrapComplete(players); err != nil {
			log.Printf("❌ Error %s player bootstrap batch: %v", logContext, err)
		} else {
			totalProcessed += len(players)
			if logContext == "inserting" {
				log.Printf("✅ Batch processed. Total received: %d, Total processed: %d", totalReceived, totalProcessed)
			} else if logContext == "flushing" {
				log.Printf("✅ Flush completed. Total received: %d, Total processed: %d", totalReceived, totalProcessed)
			}
		}
		batch = make(map[int]*fpl.PlayerBootstrapMessage)
	}

	for {
		select {
		case msg := <-messages:
			var player *fpl.PlayerBootstrapMessage
			if err := json.Unmarshal(msg.Value, &player); err != nil {
				log.Printf("Error unmarshaling player bootstrap message: %v, raw: %s\n", err, string(msg.Value))
				continue
			}
			batch[int(player.Player.Id)] = player
			totalReceived++

			if len(batch) >= h.config.BatchSize {
				flushBatch("inserting")
			}

		case <-flushTicker.C:
			flushBatch("flushing")

		case err := <-errors:
			if err != nil {
				log.Println("Error consuming player bootstrap message:", err)
			}

		case <-ctx.Done():
			flushBatch("inserting on shutdown")
			return
		}
	}
}

func (h *Handler) handlePlayerStats(ctx context.Context) {
	batchProcessWithSlice(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplPlayersStats.Name],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplPlayersStats.Name,
		func(p *fpl.PlayerBootstrapMessage) int { return int(p.Player.Id) },
		h.playerRepo.InsertPlayerBootstrapComplete,
	)
}

func (h *Handler) handlePlayerMatchStats(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplPlayerMatchStats.Name],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplPlayerMatchStats.Name,
		func(msg *fpl.PlayerHistoryMessage) int { return int(msg.PlayerId) },
		func(histMsg *fpl.PlayerHistoryMessage) error {
			return h.playerRepo.InsertPlayerGameweekStats(histMsg)
		},
	)
}

func (h *Handler) handlePlayerPastHistory(ctx context.Context) {

	batchProcessWithSlice(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplPlayerHistoryStats.Name],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplPlayerHistoryStats.Name,
		func(p *fpl.PlayerPastHistoryMessage) int { return int(p.ElementCode) },
		h.playerRepo.InsertPlayerPastSeasons,
	)
}

func (h *Handler) handlePlayerExplain(ctx context.Context) {
	batchProcessWithSlice(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplLiveEvent.Name],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplLiveEvent.Name,
		func(p *fpl.LiveEventMessage) int { return int(p.PlayerId) },
		h.playerRepo.InsertPlayerGameweekExplain,
	)
}

func (h *Handler) handleManagers(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplEntry.Name],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplEntry.Name,
		func(e *fpl.EntryMessage) int { return int(e.Entry.Id) },
		func(e *fpl.EntryMessage) error {
			return h.managerRepo.InsertManagerInfo(e)
		},
	)
}

func (h *Handler) handleManagerPicks(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplEntryPicks.Name],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplEntryPicks.Name,
		func(p *fpl.EntryEventPicksMessage) int { return int(p.EntryId) },
		func(p *fpl.EntryEventPicksMessage) error {
			return h.managerRepo.InsertManagerPicks(p)
		},
	)
}

func (h *Handler) handleManagerTransfers(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplEntryTransfers.Name],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplEntryTransfers.Name,
		func(t *fpl.EntryTransfersMessage) int { return int(t.EntryId) },
		func(t *fpl.EntryTransfersMessage) error {
			return h.managerRepo.InsertManagerTransfers(t)
		},
	)
}

func (h *Handler) handleManagerHistory(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplEntryHistory.Name],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplEntryHistory.Name,
		func(hi *fpl.EntryHistoryMessage) int { return int(hi.EntryId) },
		func(hi *fpl.EntryHistoryMessage) error {
			return h.managerRepo.InsertManagerFullHistory(hi)
		},
	)
}
