package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/imadbelkat1/indexer-service/config"
	"github.com/imadbelkat1/indexer-service/internal/fpl_repositories"
	"github.com/imadbelkat1/kafka"
	kafkaConfig "github.com/imadbelkat1/kafka/config"
	"github.com/imadbelkat1/shared/models"
)

type Handler struct {
	config      *config.IndexerConfig
	kafkaConfig *kafkaConfig.KafkaConfig
	consumer    *kafka.Consumer
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
	return &Handler{
		config:      config,
		kafkaConfig: kafkaConfig,
		playerRepo:  playerRepo,
		teamRepo:    teamRepo,
		fixtureRepo: fixtureRepo,
		managerRepo: managerRepo,
	}
}

type HandlerFunc func(ctx context.Context)

func (h *Handler) Route(ctx context.Context, topic string) {
	handlers := map[string]HandlerFunc{
		h.kafkaConfig.TopicsName.FplFixtures:           h.handleFixtures,
		h.kafkaConfig.TopicsName.FplTeams:              h.handleTeams,
		h.kafkaConfig.TopicsName.FplPlayersBootstrap:   h.handlePlayerBootstrap,
		h.kafkaConfig.TopicsName.FplPlayersStats:       h.handlePlayerStats,
		h.kafkaConfig.TopicsName.FplPlayerMatchStats:   h.handlePlayerMatchStats,
		h.kafkaConfig.TopicsName.FplPlayerHistoryStats: h.handlePlayerPastHistory,
	}

	if fn, ok := handlers[topic]; ok {
		go fn(ctx)
	}
}

func (h *Handler) handleFixtures(ctx context.Context) {
	batch := make(map[int]models.FixtureMessage)
	fixtureConsumer := kafka.NewConsumer(h.kafkaConfig, h.kafkaConfig.TopicsName.FplFixtures, h.kafkaConfig.ConsumersGroupID.FplFixtures)
	messages, errors := fixtureConsumer.Subscribe(ctx)

	flushTicker := time.NewTicker(h.config.FlushInterval)
	defer flushTicker.Stop()

	for {
		select {
		case msg := <-messages:
			var fixture models.FixtureMessage
			err := json.Unmarshal(msg.Value, &fixture)
			if err != nil {
				log.Printf("Error unmarshaling fixture message: %v, raw: %s\n", err, string(msg.Value))
				continue
			}
			batch[fixture.Fixture.ID] = fixture

			if len(batch) >= h.config.BatchSize {
				fixtures := make([]models.FixtureMessage, 0, len(batch))
				for _, f := range batch {
					fixtures = append(fixtures, f)
				}
				err := h.fixtureRepo.InsertFixtures(fixtures)
				if err != nil {
					log.Println("Error inserting fixtures batch:", err)
				}
				batch = make(map[int]models.FixtureMessage)
			}

		case <-flushTicker.C:
			if len(batch) > 0 {
				fixtures := make([]models.FixtureMessage, 0, len(batch))
				for _, f := range batch {
					fixtures = append(fixtures, f)
				}
				if err := h.fixtureRepo.InsertFixtures(fixtures); err != nil {
					log.Println("Error flushing fixtures batch:", err)
				}
				batch = make(map[int]models.FixtureMessage)
			}

		case err := <-errors:
			if err != nil {
				log.Println("Error consuming fixture message:", err)
			}

		case <-ctx.Done():
			if len(batch) > 0 {
				fixtures := make([]models.FixtureMessage, 0, len(batch))
				for _, f := range batch {
					fixtures = append(fixtures, f)
				}
				err := h.fixtureRepo.InsertFixtures(fixtures)
				if err != nil {
					log.Println("Error inserting fixtures batch:", err)
				}
			}
			return
		}
	}
}

func (h *Handler) handleTeams(ctx context.Context) {
	batch := make(map[int]models.TeamMessage)
	teamConsumer := kafka.NewConsumer(h.kafkaConfig, h.kafkaConfig.TopicsName.FplTeams, h.kafkaConfig.ConsumersGroupID.FplTeams)
	messages, errors := teamConsumer.Subscribe(ctx)

	flushTicker := time.NewTicker(h.config.FlushInterval)
	defer flushTicker.Stop()

	for {
		select {
		case msg := <-messages:
			var team models.TeamMessage
			err := json.Unmarshal(msg.Value, &team)
			if err != nil {
				log.Printf("Error unmarshaling team message: %v, raw: %s\n", err, string(msg.Value))
				continue
			}
			batch[team.Team.ID] = team

			if len(batch) >= h.config.BatchSize {
				teams := make([]models.TeamMessage, 0, len(batch))
				for _, t := range batch {
					teams = append(teams, t)
				}
				err := h.teamRepo.InsertTeams(teams)
				if err != nil {
					log.Println("Error inserting teams batch:", err)
				}
				batch = make(map[int]models.TeamMessage)
			}

		case <-flushTicker.C:
			if len(batch) > 0 {
				teams := make([]models.TeamMessage, 0, len(batch))
				for _, t := range batch {
					teams = append(teams, t)
				}
				if err := h.teamRepo.InsertTeams(teams); err != nil {
					log.Println("Error flushing teams batch:", err)
				}
				batch = make(map[int]models.TeamMessage)
			}

		case err := <-errors:
			if err != nil {
				log.Println("Error consuming team message:", err)
			}

		case <-ctx.Done():
			if len(batch) > 0 {
				teams := make([]models.TeamMessage, 0, len(batch))
				for _, t := range batch {
					teams = append(teams, t)
				}
				err := h.teamRepo.InsertTeams(teams)
				if err != nil {
					log.Println("Error inserting teams batch:", err)
				}
			}
			return
		}
	}
}

func (h *Handler) handlePlayerBootstrap(ctx context.Context) {
	batch := make(map[int]models.PlayerBootstrapMessage)
	consumer := kafka.NewConsumer(h.kafkaConfig, h.kafkaConfig.TopicsName.FplPlayersBootstrap, h.kafkaConfig.ConsumersGroupID.FplPlayers)
	messages, errors := consumer.Subscribe(ctx)

	flushTicker := time.NewTicker(h.config.FlushInterval)
	verifyTicker := time.NewTicker(30 * time.Second) // Verify every 30 seconds
	defer flushTicker.Stop()
	defer verifyTicker.Stop()

	totalReceived := 0
	totalProcessed := 0

	for {
		select {
		case msg := <-messages:
			var player models.PlayerBootstrapMessage
			err := json.Unmarshal(msg.Value, &player)
			if err != nil {
				log.Printf("Error unmarshaling player bootstrap message: %v, raw: %s\n", err, string(msg.Value))
				continue
			}
			batch[player.Player.ID] = player
			totalReceived++

			if len(batch) >= h.config.BatchSize {
				players := make([]models.PlayerBootstrapMessage, 0, len(batch))
				for _, p := range batch {
					players = append(players, p)
				}
				if err := h.playerRepo.InsertPlayerBootstrapComplete(players); err != nil {
					log.Printf("❌ Error inserting player bootstrap batch: %v", err)
				} else {
					totalProcessed += len(players)
					log.Printf("✅ Batch processed. Total received: %d, Total processed: %d", totalReceived, totalProcessed)
				}
				batch = make(map[int]models.PlayerBootstrapMessage)
			}

		case <-flushTicker.C:
			if len(batch) > 0 {
				players := make([]models.PlayerBootstrapMessage, 0, len(batch))
				for _, p := range batch {
					players = append(players, p)
				}
				if err := h.playerRepo.InsertPlayerBootstrapComplete(players); err != nil {
					log.Printf("❌ Error flushing player bootstrap batch: %v", err)
				} else {
					totalProcessed += len(players)
					log.Printf("✅ Flush completed. Total received: %d, Total processed: %d", totalReceived, totalProcessed)
				}
				batch = make(map[int]models.PlayerBootstrapMessage)
			}

		case <-verifyTicker.C:
			// Verify how many players are actually in the database
			if h.playerRepo != nil {
				count, err := h.playerRepo.CountPlayers()
				if err != nil {
					log.Printf("⚠️  Error counting players: %v", err)
				} else {
					log.Printf("📊 Players in database: %d (received: %d, processed: %d)", count, totalReceived, totalProcessed)
				}
			}

		case err := <-errors:
			if err != nil {
				log.Println("Error consuming player bootstrap message:", err)
			}

		case <-ctx.Done():
			if len(batch) > 0 {
				players := make([]models.PlayerBootstrapMessage, 0, len(batch))
				for _, p := range batch {
					players = append(players, p)
				}
				if err := h.playerRepo.InsertPlayerBootstrapComplete(players); err != nil {
					log.Printf("❌ Error inserting player bootstrap batch on shutdown: %v", err)
				} else {
					totalProcessed += len(players)
				}
			}

			// Final count
			if h.playerRepo != nil {
				count, _ := h.playerRepo.CountPlayers()
				log.Printf("🏁 Final stats - Players in DB: %d, Received: %d, Processed: %d", count, totalReceived, totalProcessed)
			}
			return
		}
	}
}

func (h *Handler) handlePlayerStats(ctx context.Context) {
	batch := make(map[int]models.PlayerBootstrapMessage)
	consumer := kafka.NewConsumer(h.kafkaConfig, h.kafkaConfig.TopicsName.FplPlayersStats, h.kafkaConfig.ConsumersGroupID.FplPlayersStats)
	messages, errors := consumer.Subscribe(ctx)

	flushTicker := time.NewTicker(h.config.FlushInterval)
	defer flushTicker.Stop()

	for {
		select {
		case msg := <-messages:
			var player models.PlayerBootstrapMessage
			err := json.Unmarshal(msg.Value, &player)
			if err != nil {
				log.Printf("Error unmarshaling player stats message: %v, raw: %s\n", err, string(msg.Value))
				continue
			}
			batch[player.Player.ID] = player

			if len(batch) >= h.config.BatchSize {
				players := make([]models.PlayerBootstrapMessage, 0, len(batch))
				for _, p := range batch {
					players = append(players, p)
				}
				if err := h.playerRepo.InsertPlayerBootstrapComplete(players); err != nil {
					log.Println("Error inserting player stats batch:", err)
				}
				batch = make(map[int]models.PlayerBootstrapMessage)
			}

		case <-flushTicker.C:
			if len(batch) > 0 {
				players := make([]models.PlayerBootstrapMessage, 0, len(batch))
				for _, p := range batch {
					players = append(players, p)
				}
				if err := h.playerRepo.InsertPlayerBootstrapComplete(players); err != nil {
					log.Println("Error flushing player stats batch:", err)
				}
				batch = make(map[int]models.PlayerBootstrapMessage)
			}

		case err := <-errors:
			if err != nil {
				log.Println("Error consuming player stats message:", err)
			}

		case <-ctx.Done():
			if len(batch) > 0 {
				players := make([]models.PlayerBootstrapMessage, 0, len(batch))
				for _, p := range batch {
					players = append(players, p)
				}
				if err := h.playerRepo.InsertPlayerBootstrapComplete(players); err != nil {
					log.Println("Error inserting player stats batch on shutdown:", err)
				}
			}
			return
		}
	}
}

func (h *Handler) handlePlayerMatchStats(ctx context.Context) {
	batch := make(map[int]models.PlayerHistoryMessage)

	consumer := kafka.NewConsumer(h.kafkaConfig, h.kafkaConfig.TopicsName.FplPlayerMatchStats, h.kafkaConfig.ConsumersGroupID.FplPlayersStats)
	messages, errors := consumer.Subscribe(ctx)

	flushTicker := time.NewTicker(h.config.FlushInterval)
	defer flushTicker.Stop()

	for {
		select {
		case msg := <-messages:
			var historyMsg models.PlayerHistoryMessage
			err := json.Unmarshal(msg.Value, &historyMsg)
			if err != nil {
				log.Printf("Error unmarshaling player match stats message: %v, raw: %s\n", err, string(msg.Value))
				continue
			}
			batch[historyMsg.PlayerID] = historyMsg

			if len(batch) >= h.config.BatchSize {
				for playerID, histMsg := range batch {
					if err := h.playerRepo.InsertPlayerGameweekStats(playerID, histMsg); err != nil {
						log.Printf("Error inserting player match stats for player %d: %v\n", playerID, err)
					}
				}
				batch = make(map[int]models.PlayerHistoryMessage)
			}

		case <-flushTicker.C:
			if len(batch) > 0 {
				for playerID, histMsg := range batch {
					if err := h.playerRepo.InsertPlayerGameweekStats(playerID, histMsg); err != nil {
						log.Printf("Error flushing player match stats for player %d: %v\n", playerID, err)
					}
				}
				batch = make(map[int]models.PlayerHistoryMessage)
			}

		case err := <-errors:
			if err != nil {
				log.Println("Error consuming player match stats message:", err)
			}

		case <-ctx.Done():
			if len(batch) > 0 {
				for playerID, histMsg := range batch {
					if err := h.playerRepo.InsertPlayerGameweekStats(playerID, histMsg); err != nil {
						log.Printf("Error inserting player match stats for player %d on shutdown: %v\n", playerID, err)
					}
				}
			}
			return
		}
	}
}

func (h *Handler) handlePlayerPastHistory(ctx context.Context) {
	type PlayerPastMessage struct {
		PlayerCode  int                        `json:"player_code"`
		PastHistory []models.PlayerPastHistory `json:"past_history"`
	}

	batch := make(map[int]PlayerPastMessage)
	consumer := kafka.NewConsumer(h.kafkaConfig, h.kafkaConfig.TopicsName.FplPlayerHistoryStats, h.kafkaConfig.ConsumersGroupID.FplPlayersStats)
	messages, errors := consumer.Subscribe(ctx)

	flushTicker := time.NewTicker(h.config.FlushInterval)
	defer flushTicker.Stop()

	for {
		select {
		case msg := <-messages:
			var pastMsg PlayerPastMessage
			err := json.Unmarshal(msg.Value, &pastMsg)
			if err != nil {
				log.Printf("Error unmarshaling player past history message: %v, raw: %s\n", err, string(msg.Value))
				continue
			}
			batch[pastMsg.PlayerCode] = pastMsg

			if len(batch) >= h.config.BatchSize {
				for playerCode, pMsg := range batch {
					if err := h.playerRepo.InsertPlayerPastSeasons(playerCode, pMsg.PastHistory); err != nil {
						log.Printf("Error inserting player past history for player code %d: %v\n", playerCode, err)
					}
				}
				batch = make(map[int]PlayerPastMessage)
			}

		case <-flushTicker.C:
			if len(batch) > 0 {
				for playerCode, pMsg := range batch {
					if err := h.playerRepo.InsertPlayerPastSeasons(playerCode, pMsg.PastHistory); err != nil {
						log.Printf("Error flushing player past history for player code %d: %v\n", playerCode, err)
					}
				}
				batch = make(map[int]PlayerPastMessage)
			}

		case err := <-errors:
			if err != nil {
				log.Println("Error consuming player past history message:", err)
			}

		case <-ctx.Done():
			if len(batch) > 0 {
				for playerCode, pMsg := range batch {
					if err := h.playerRepo.InsertPlayerPastSeasons(playerCode, pMsg.PastHistory); err != nil {
						log.Printf("Error inserting player past history for player code %d on shutdown: %v\n", playerCode, err)
					}
				}
			}
			return
		}
	}
}
