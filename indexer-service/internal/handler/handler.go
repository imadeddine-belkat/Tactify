package handler

import (
	"github.com/imadbelkat1/indexer-service/internal/fpl_repositories"
	kafkaConfig "github.com/imadbelkat1/kafka/config"
)

type Handler struct {
	kafkaConfig *kafkaConfig.KafkaConfig
	playerRepo  *fpl_repositories.PlayerRepo
	teamRepo    *fpl_repositories.TeamRepo
	fixtureRepo *fpl_repositories.FixtureRepo
	managerRepo *fpl_repositories.ManagerRepo
}

func NewHandler(
	kafkaConfig *kafkaConfig.KafkaConfig,
	playerRepo *fpl_repositories.PlayerRepo,
	teamRepo *fpl_repositories.TeamRepo,
	fixtureRepo *fpl_repositories.FixtureRepo,
	managerRepo *fpl_repositories.ManagerRepo,
) *Handler {
	return &Handler{
		kafkaConfig: kafkaConfig,
		playerRepo:  playerRepo,
		teamRepo:    teamRepo,
		fixtureRepo: fixtureRepo,
		managerRepo: managerRepo,
	}
}

func (h *Handler) Route(topic string, message []byte) {
	switch topic {
	case h.kafkaConfig.TopicsName.FplLiveEvent:
	case h.kafkaConfig.TopicsName.FplTeams:
	case h.kafkaConfig.TopicsName.FplFixtures:
	case h.kafkaConfig.TopicsName.FplFixtureDetails:
	case h.kafkaConfig.TopicsName.FplPlayersStats:
	case h.kafkaConfig.TopicsName.FplPlayerMatchStats:
	case h.kafkaConfig.TopicsName.FplPlayerHistoryStats:
	case h.kafkaConfig.TopicsName.FplEntry:
	case h.kafkaConfig.TopicsName.FplEntryEvent:
	case h.kafkaConfig.TopicsName.FplEntryHistory:
	case h.kafkaConfig.TopicsName.FplEntryPicks:
	case h.kafkaConfig.TopicsName.FplEntryTransfers:
	case h.kafkaConfig.TopicsName.FplLeagueClassicStanding:
	case h.kafkaConfig.TopicsName.FplLeagueH2hStanding:
	default:
		return
	}
}
