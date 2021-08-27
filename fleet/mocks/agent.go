package mocks

import (
	"context"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/types"
)

var _ fleet.AgentRepository = (*agentRepositoryMock)(nil)

type agentRepositoryMock struct {
	counter    uint64
	agentsMock map[string]fleet.Agent
}

func (a agentRepositoryMock) RetrieveByID(ctx context.Context, ownerID string, thingID string) (fleet.Agent, error) {
	if _, ok := a.agentsMock[thingID]; ok {
		if a.agentsMock[thingID].MFOwnerID != ownerID {
			return fleet.Agent{}, fleet.ErrNotFound
		}
		return a.agentsMock[thingID], nil
	}
	return fleet.Agent{}, fleet.ErrNotFound
}

func (a agentRepositoryMock) RetrieveMatchingAgents(ctx context.Context, ownerID string, tags types.Tags) (types.Metadata, error) {
	return nil, nil
}

func (a agentRepositoryMock) UpdateHeartbeatByIDWithChannel(ctx context.Context, agent fleet.Agent) error {
	panic("implement me")
}

func (a agentRepositoryMock) Save(ctx context.Context, agent fleet.Agent) error {
	for _, ag := range a.agentsMock {
		if ag.Name == agent.Name && ag.MFOwnerID == agent.MFOwnerID {
			return fleet.ErrConflict
		}
	}
	a.agentsMock[agent.MFThingID] = agent
	return nil
}

func (a agentRepositoryMock) UpdateAgentByID(ctx context.Context, ownerID string, agent fleet.Agent) error {
	if _, ok := a.agentsMock[agent.MFThingID]; ok {
		if a.agentsMock[agent.MFThingID].MFOwnerID != ownerID {
			return fleet.ErrUpdateEntity
		}
		a.agentsMock[agent.MFThingID] = agent
		return nil
	}
	return fleet.ErrNotFound
}

func (a agentRepositoryMock) UpdateDataByIDWithChannel(ctx context.Context, agent fleet.Agent) error {
	panic("implement me")
}

func (a agentRepositoryMock) RetrieveByIDWithChannel(ctx context.Context, thingID string, channelID string) (fleet.Agent, error) {
	panic("implement me")
}

func (a agentRepositoryMock) RetrieveAll(ctx context.Context, owner string, pm fleet.PageMetadata) (fleet.Page, error) {
	first := uint64(pm.Offset)
	last := first + uint64(pm.Limit)

	var agents []fleet.Agent
	id := uint64(0)
	for _, v := range a.agentsMock {
		if v.MFOwnerID == owner && id >= first && id < last {
			agents = append(agents, v)
		}
		id++
	}

	agents = sortAgents(pm, agents)

	pageAgentGroup := fleet.Page{
		PageMetadata: fleet.PageMetadata{
			Total: a.counter,
		},
		Agents: agents,
	}
	return pageAgentGroup, nil
}

func (a agentRepositoryMock) RetrieveAllByAgentGroupID(ctx context.Context, owner string, agentGroupID string, onlinishOnly bool) ([]fleet.Agent, error) {
	return []fleet.Agent{}, nil
}

func (a agentRepositoryMock) Delete(ctx context.Context, ownerID, thingID string) error {
	if _, ok := a.agentsMock[thingID]; ok {
		if a.agentsMock[thingID].MFOwnerID == ownerID {
			delete(a.agentsMock, thingID)
		}
	}

	return nil
}

func NewAgentRepositoryMock() fleet.AgentRepository {
	return &agentRepositoryMock{
		agentsMock: make(map[string]fleet.Agent),
	}
}
