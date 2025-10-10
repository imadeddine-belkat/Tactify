package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/imadbelkat1/fpl-service/config"
	fpl_api "github.com/imadbelkat1/fpl-service/internal/api"
	"github.com/imadbelkat1/kafka"
	"github.com/imadbelkat1/shared/models"
)

type ManagersApiService struct {
	Config   *config.FplConfig
	Client   *fpl_api.FplApiClient
	Producer *kafka.Producer
}

func (s *ManagersApiService) UpdateManager(ctx context.Context, managerId int, eventId int) error {
	info, err := s.GetManagerInfo(ctx, managerId)
	if err != nil {
		return fmt.Errorf("fetching manager info data: %w", err)
	}

	picks, err := s.GetManagerPicks(ctx, managerId, eventId)
	if err != nil {
		return fmt.Errorf("fetching manager picks data: %w", err)
	}

	history, err := s.GetManagerHistory(ctx, managerId)
	if err != nil {
		return fmt.Errorf("fetching manager history data: %w", err)
	}

	transfers, err := s.GetManagerTransfers(ctx, managerId)
	if err != nil {
		return fmt.Errorf("fetching manager transfers data: %w", err)
	}

	if err := s.publishManager(ctx, info, picks, history, transfers); err != nil {
		return fmt.Errorf("update manager: %w", err)
	}

	return nil
}

func (s *ManagersApiService) publishManager(ctx context.Context, info *models.EntryMessage, picks *models.EntryEventPicksMessage, history *models.EntryHistoryMessage, transfers *models.EntryTransfersMessage) error {
	entryTopic := s.Config.KafkaConfig.TopicsName.FplEntry
	entryEventTopic := s.Config.KafkaConfig.TopicsName.FplEntryPicks
	entryHistoryTopic := s.Config.KafkaConfig.TopicsName.FplEntryHistory
	entryTransfersTopic := s.Config.KafkaConfig.TopicsName.FplEntryTransfers

	var publishWg sync.WaitGroup
	publishWg.Add(1)
	go func() {
		if info != nil {
			value, err := json.Marshal(info)
			if err == nil {
				key := []byte(fmt.Sprintf("%d", info.Entry.ID))
				err := s.Producer.Publish(ctx, entryTopic, key, value)
				if err != nil {
					fmt.Printf("Failed to publish entry message: %v\n", err)
				}
			}
		}

		if picks != nil {
			value, err := json.Marshal(picks)
			if err == nil {
				key := []byte(fmt.Sprintf("%d_%d", picks.EntryId, picks.EventId))
				err := s.Producer.Publish(ctx, entryEventTopic, key, value)
				if err != nil {
					fmt.Printf("Failed to publish entry event picks message: %v\n", err)
				}
			}
		}

		if history != nil {
			value, err := json.Marshal(history)
			if err == nil {
				key := []byte(fmt.Sprintf("%d-%d", history.EntryId, history.SeasonId))
				err := s.Producer.Publish(ctx, entryHistoryTopic, key, value)
				if err != nil {
					fmt.Printf("Failed to publish entry history message: %v\n", err)
				}
			}
		}

		if transfers != nil {
			value, err := json.Marshal(transfers)
			if err == nil {
				key := []byte(fmt.Sprintf("%d-%d", transfers.EntryId, transfers.SeasonId))
				err := s.Producer.Publish(ctx, entryTransfersTopic, key, value)
				if err != nil {
					fmt.Printf("Failed to publish entry transfers message: %v\n", err)
				}
			}
		}
		defer publishWg.Done()
	}()

	publishWg.Wait()
	return nil
}

func (s *ManagersApiService) GetManagerInfo(ctx context.Context, managerId int) (*models.EntryMessage, error) {
	var entry models.EntryMessage

	entryEndpoint := s.Config.FplApi.Entry
	endpoint := fmt.Sprintf(entryEndpoint, managerId)
	log.Println(endpoint)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &entry.Entry); err != nil {
		return nil, err
	}

	entry.SeasonId = s.Config.FplApi.CurrentSeasonID

	return &entry, nil
}

func (s *ManagersApiService) GetManagerPicks(ctx context.Context, managerId int, eventId int) (*models.EntryEventPicksMessage, error) {
	var entryEvent models.EntryEventPicksMessage

	entryEventEndpoint := s.Config.FplApi.EntryPicks
	endpoint := fmt.Sprintf(entryEventEndpoint, managerId, eventId)
	log.Println(endpoint)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &entryEvent.Picks); err != nil {
		return nil, err
	}

	entryEvent.EventId = eventId
	entryEvent.EntryId = managerId

	return &entryEvent, nil
}

func (s *ManagersApiService) GetManagerHistory(ctx context.Context, managerId int) (*models.EntryHistoryMessage, error) {
	var entryHistory models.EntryHistoryMessage

	entryHistoryEndpoint := s.Config.FplApi.EntryHistory
	endpoint := fmt.Sprintf(entryHistoryEndpoint, managerId)
	log.Println(endpoint)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &entryHistory.EntryHistory); err != nil {
		return nil, err
	}

	entryHistory.EntryId = managerId
	entryHistory.SeasonId = s.Config.FplApi.CurrentSeasonID

	return &entryHistory, nil
}

func (s *ManagersApiService) GetManagerTransfers(ctx context.Context, managerId int) (*models.EntryTransfersMessage, error) {
	var entryTransfers models.EntryTransfersMessage

	entryTransfersEndpoint := s.Config.FplApi.EntryTransfers
	endpoint := fmt.Sprintf(entryTransfersEndpoint, managerId)
	log.Println(endpoint)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &entryTransfers.Transfers); err != nil {
		return nil, err
	}

	entryTransfers.EntryId = managerId
	entryTransfers.SeasonId = s.Config.FplApi.CurrentSeasonID

	return &entryTransfers, nil
}
