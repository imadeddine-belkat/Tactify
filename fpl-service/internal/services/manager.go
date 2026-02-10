package services

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/imadeddine-belkat/fpl-service/config"
	fpl_api "github.com/imadeddine-belkat/fpl-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	fpl "github.com/imadeddine-belkat/tactify-protos/go/fpl/v1"
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

func (s *ManagersApiService) publishManager(ctx context.Context, info *fpl.EntryMessage, picks *fpl.EntryEventPicksMessage, history *fpl.EntryHistoryMessage, transfers *fpl.EntryTransfersMessage) error {
	entryTopic := s.Config.KafkaConfig.TopicsName.FplEntry.Name
	entryEventTopic := s.Config.KafkaConfig.TopicsName.FplEntryPicks.Name
	entryHistoryTopic := s.Config.KafkaConfig.TopicsName.FplEntryHistory.Name
	entryTransfersTopic := s.Config.KafkaConfig.TopicsName.FplEntryTransfers.Name

	var publishWg sync.WaitGroup
	publishWg.Add(1)
	go func() {
		if &info != nil {
			key := []byte(fmt.Sprintf("%d-%d", info.Entry.Id, info.SeasonId))
			err := s.Producer.PublishWithProcess(ctx, info, entryTopic, key)
			if err != nil {
				fmt.Printf("Failed to publish entry message: %v\n", err)
			}
		}

		if &picks != nil {
			key := []byte(fmt.Sprintf("%d-%d", picks.EntryId, picks.Event))
			err := s.Producer.PublishWithProcess(ctx, picks, entryEventTopic, key)
			if err != nil {
				fmt.Printf("Failed to publish entry picks message: %v\n", err)
			}
		}

		if &history != nil {
			key := []byte(fmt.Sprintf("%d-%d", history.EntryId, history.SeasonId))
			err := s.Producer.PublishWithProcess(ctx, history, entryHistoryTopic, key)
			if err != nil {
				fmt.Printf("Failed to publish entry history message: %v\n", err)
			}
		}

		if &transfers != nil {
			key := []byte(fmt.Sprintf("%d-%d", transfers.EntryId, transfers.SeasonId))
			err := s.Producer.PublishWithProcess(ctx, transfers, entryTransfersTopic, key)
			if err != nil {
				fmt.Printf("Failed to publish entry transfers message: %v\n", err)
			}
		}
		defer publishWg.Done()
	}()

	publishWg.Wait()
	return nil
}

func (s *ManagersApiService) GetManagerInfo(ctx context.Context, managerId int) (*fpl.EntryMessage, error) {
	var entry fpl.EntryMessage

	entryEndpoint := s.Config.FplApi.Entry
	endpoint := fmt.Sprintf(entryEndpoint, managerId)
	log.Println(endpoint)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &entry.Entry); err != nil {
		return nil, err
	}

	entry.SeasonId = s.Config.FplApi.CurrentSeasonID

	return &entry, nil
}

func (s *ManagersApiService) GetManagerPicks(ctx context.Context, managerId int, eventId int) (*fpl.EntryEventPicksMessage, error) {
	var entryEvent fpl.EntryEventPicksMessage

	entryEventEndpoint := s.Config.FplApi.EntryPicks
	endpoint := fmt.Sprintf(entryEventEndpoint, managerId, eventId)
	log.Println(endpoint)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &entryEvent.Picks); err != nil {
		return nil, err
	}

	entryEvent.Event = int32(eventId)
	entryEvent.EntryId = int32(managerId)
	entryEvent.SeasonId = s.Config.FplApi.CurrentSeasonID

	return &entryEvent, nil
}

func (s *ManagersApiService) GetManagerHistory(ctx context.Context, managerId int) (*fpl.EntryHistoryMessage, error) {
	var entryHistory fpl.EntryHistoryMessage

	entryHistoryEndpoint := s.Config.FplApi.EntryHistory
	endpoint := fmt.Sprintf(entryHistoryEndpoint, managerId)
	log.Println(endpoint)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &entryHistory.EntryHistory); err != nil {
		return nil, err
	}

	entryHistory.EntryId = int32(managerId)
	entryHistory.SeasonId = s.Config.FplApi.CurrentSeasonID
	for i := range entryHistory.EntryHistory.Past {
		seasonID := s.Config.MapSeasonNameToID(entryHistory.EntryHistory.Past[i].SeasonName)
		entryHistory.EntryHistory.Past[i].SeasonId = seasonID
	}

	return &entryHistory, nil
}

func (s *ManagersApiService) GetManagerTransfers(ctx context.Context, managerId int) (*fpl.EntryTransfersMessage, error) {
	var entryTransfers fpl.EntryTransfersMessage

	entryTransfersEndpoint := s.Config.FplApi.EntryTransfers
	endpoint := fmt.Sprintf(entryTransfersEndpoint, managerId)
	log.Println(endpoint)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &entryTransfers.Transfers); err != nil {
		return nil, err
	}

	entryTransfers.EntryId = int32(managerId)
	entryTransfers.SeasonId = s.Config.FplApi.CurrentSeasonID

	return &entryTransfers, nil
}
