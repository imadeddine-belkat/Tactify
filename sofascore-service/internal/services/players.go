package services

import (
	"context"
	"fmt"

	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascoreapi "github.com/imadeddine-belkat/sofascore-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	sofascore "github.com/imadeddine-belkat/tactify-protos/go/sofascore/v1"
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

	for _, team := range teams {
		players, err := p.GetPlayersInfo(ctx, int(team.Id), leagueId, seasonId)
		if err != nil {
			return fmt.Errorf("error getting players info for team %d: %w", team.Id, err)
		}

		for _, player := range players.TopPlayers.Rating {
			pl := player
			playerMessage := &sofascore.PlayerMessage{
				SeasonId: int32(seasonId),
				LeagueId: int32(leagueId),
				TeamId:   team.Id,
				TeamName: team.Name,
				Player:   pl.Player,
			}

			if err := p.publishPlayerInfo(ctx, playerMessage); err != nil {
				return fmt.Errorf("error publishing player info for player %d: %w", pl.Player.Id, err)
			}
		}
	}

	return nil
}

func (p *PlayersService) GetTeamIDs(ctx context.Context, seasonId, leagueId int) ([]*sofascore.Team, error) {
	var teams []*sofascore.Team

	standing, err := p.Standing.GetLeagueStanding(ctx, seasonId, leagueId)
	if err != nil {
		return nil, fmt.Errorf("error getting league standing for league %d, season %d: %w", leagueId, seasonId, err)
	}

	for _, s := range standing.Standings {
		for _, row := range s.Rows {
			team := sofascore.Team{
				Id:   row.Team.Id,
				Name: row.Team.Name,
			}
			teams = append(teams, &team)
		}
	}

	return teams, nil
}
func (p *PlayersService) GetPlayersInfo(ctx context.Context, teamId, leagueId, seasonId int) (*sofascore.TopPlayers, error) {
	players := &sofascore.TopPlayers{}

	playerEndpoint := p.Config.SofascoreApi.TeamEndpoints.TeamTopPlayerStats
	endpoint := fmt.Sprintf(playerEndpoint, teamId, leagueId, seasonId)

	if err := p.Client.GetAndUnmarshal(ctx, endpoint, players); err != nil {
		return nil, fmt.Errorf("fetching players info data: %w", err)
	}

	return players, nil
}

func (p *PlayersService) publishPlayerInfo(ctx context.Context, playerMessage *sofascore.PlayerMessage) error {
	topic := p.Config.KafkaConfig.TopicsName.SofascorePlayerInfo.Name

	key := []byte(fmt.Sprintf("%d-%d", playerMessage.TeamId, playerMessage.Player.Id))

	if err := p.Producer.PublishWithProcess(ctx, playerMessage, topic, key); err != nil {
		return fmt.Errorf("error publishing player info for player %d: %w", playerMessage.Player.Id, err)
	}

	return nil
}
