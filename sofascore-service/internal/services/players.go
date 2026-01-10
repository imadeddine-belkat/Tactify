package services

import (
	"context"
	"fmt"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/shared/sofascore_models"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascoreapi "github.com/imadeddine-belkat/sofascore-service/internal/api"
)

type PlayersService struct {
	Config   config.SofascoreConfig
	Client   *sofascoreapi.SofascoreApiClient
	Standing *LeagueStandingService
	Producer *kafka.Producer
}

func (p *PlayersService) UpdateLeaguePlayersInfo(ctx context.Context, seasonId, leagueId int) error {
	teams, err := p.GetTeamIDs(ctx, seasonId, leagueId)
	if err != nil {
		return fmt.Errorf("error getting team ids for league %d, season %d: %w", leagueId, seasonId, err)
	}

	for _, team := range *teams {
		players, err := p.GetPlayersInfo(ctx, team.ID)
		if err != nil {
			return fmt.Errorf("error getting players info for team %d: %w", team.ID, err)
		}

		for _, player := range players.TeamPlayers {
			pl := player
			playerMessage := &sofascore_models.PlayerMessage{
				SeasonID: seasonId,
				LeagueID: leagueId,
				TeamID:   team.ID,
				Player:   pl.Player,
			}

			if err := p.publishPlayerInfo(ctx, playerMessage); err != nil {
				return fmt.Errorf("error publishing player info for player %d: %w", pl.Player.ID, err)
			}
		}
	}

	return nil
}

func (p *PlayersService) GetTeamIDs(ctx context.Context, seasonId, leagueId int) (*[]sofascore_models.Team, error) {
	teams := &[]sofascore_models.Team{}

	standing, err := p.Standing.GetLeagueStanding(ctx, seasonId, leagueId)
	if err != nil {
		return nil, fmt.Errorf("error getting league standing for league %d, season %d: %w", leagueId, seasonId, err)
	}

	for _, s := range standing.Standings {
		for _, row := range s.Rows {
			team := sofascore_models.Team{
				ID:   row.Team.ID,
				Name: row.Team.Name,
			}
			*teams = append(*teams, team)
		}
	}

	return teams, nil
}
func (p *PlayersService) GetPlayersInfo(ctx context.Context, teamId int) (*sofascore_models.TeamPlayers, error) {
	players := &sofascore_models.TeamPlayers{}

	playerEndpoint := p.Config.SofascoreApi.PlayerEndpoints.PlayersInfo
	endpoint := fmt.Sprintf(playerEndpoint, teamId)

	if err := p.Client.GetAndUnmarshal(ctx, endpoint, players); err != nil {
		return nil, fmt.Errorf("fetching players info data: %w", err)
	}

	return players, nil
}

func (p *PlayersService) publishPlayerInfo(ctx context.Context, playerMessage *sofascore_models.PlayerMessage) error {
	topic := p.Config.KafkaConfig.TopicsName.SofascorePlayerInfo.Name

	key := []byte(fmt.Sprintf("%d-%d", playerMessage.TeamID, playerMessage.Player.ID))

	if err := p.Producer.PublishWithProcess(ctx, playerMessage, topic, key); err != nil {
		return fmt.Errorf("error publishing player info for player %d: %w", playerMessage.Player.ID, err)
	}

	return nil
}
