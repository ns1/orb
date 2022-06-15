// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet_test

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/fleet"
	flmocks "github.com/ns1labs/orb/fleet/mocks"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	agent = fleet.Agent{
		Name:          types.Identifier{},
		MFOwnerID:     "",
		MFThingID:     "",
		MFKeyID:       "",
		MFChannelID:   "",
		Created:       time.Time{},
		OrbTags:       nil,
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: nil,
		State:         0,
		LastHBData:    nil,
		LastHB:        time.Time{},
	}
)

func TestViewAgent(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	ag, err := createAgent(t, "my-agent1", fleetService)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		id    string
		token string
		err   error
	}{
		"view a existing agent": {
			id:    ag.MFThingID,
			token: token,
			err:   nil,
		},
		"view agent with wrong credentials": {
			id:    ag.MFThingID,
			token: "wrong",
			err:   fleet.ErrUnauthorizedAccess,
		},
		"view non-existing agent": {
			id:    "9bb1b244-a199-93c2-aa03-28067b431e2c",
			token: token,
			err:   fleet.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := fleetService.ViewAgentByID(context.Background(), tc.token, tc.id)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestViewAgentMatchingGroups(t *testing.T) {

	//Setup
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	ag, err := createAgent(t, "my-agent1", fleetService)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	aGroup, err := createAgentGroup(t, "my-group1", fleetService)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	// Test cases
	cases := map[string]struct {
		id    string
		token string
		err   error
		mg    fleet.MatchingGroups
	}{
		"view a existing agent": {
			id:    ag.MFThingID,
			token: token,
			err:   nil,
			mg: fleet.MatchingGroups{
				OwnerID: aGroup.MFOwnerID,
				Groups:  []fleet.Group{{GroupID: aGroup.ID, GroupName: aGroup.Name}},
			},
		},
		"view matching groups with wrong credentials": {
			id:    ag.MFThingID,
			token: "wrong",
			err:   fleet.ErrUnauthorizedAccess,
			mg:    fleet.MatchingGroups{},
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			matchingGroups, err := fleetService.ViewAgentMatchingGroupsByID(context.Background(), tc.token, tc.id)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
			assert.Equal(t, tc.mg, matchingGroups, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestListAgents(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	var agents []fleet.Agent
	for i := 0; i < limit; i++ {
		ag, err := createAgent(t, fmt.Sprintf("my-agent-%d", i), fleetService)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
		agents = append(agents, ag)
	}

	cases := map[string]struct {
		token string
		pm    fleet.PageMetadata
		size  uint64
		err   error
	}{
		"retrieve a list of agents": {
			token: token,
			pm: fleet.PageMetadata{
				Limit:  limit,
				Offset: 0,
			},
			size: limit,
			err:  nil,
		},
		"list half": {
			token: token,
			pm: fleet.PageMetadata{
				Offset: limit / 2,
				Limit:  limit,
			},
			size: limit / 2,
			err:  nil,
		},
		"list last agent": {
			token: token,
			pm: fleet.PageMetadata{
				Offset: limit - 1,
				Limit:  limit,
			},
			size: 1,
			err:  nil,
		},
		"list empty set": {
			token: token,
			pm: fleet.PageMetadata{
				Offset: limit + 1,
				Limit:  limit,
			},
			size: 0,
			err:  nil,
		},
		"list with zero limit": {
			token: token,
			pm: fleet.PageMetadata{
				Offset: 1,
				Limit:  0,
			},
			size: 0,
			err:  nil,
		},
		"list with wrong credentials": {
			token: "wrong",
			pm: fleet.PageMetadata{
				Offset: 0,
				Limit:  0,
			},
			size: 0,
			err:  fleet.ErrUnauthorizedAccess,
		},
		"list all agents sorted by name ascendent": {
			token: token,
			pm: fleet.PageMetadata{
				Offset: 0,
				Limit:  limit,
				Order:  "name",
				Dir:    "asc",
			},
			size: limit,
			err:  nil,
		},
		"list all agents sorted by name descendent": {
			token: token,
			pm: fleet.PageMetadata{
				Offset: 0,
				Limit:  limit,
				Order:  "name",
				Dir:    "desc",
			},
			size: limit,
			err:  nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			page, err := fleetService.ListAgents(context.Background(), tc.token, tc.pm)
			size := uint64(len(page.Agents))
			assert.Equal(t, size, tc.size, fmt.Sprintf("%s: expected %d got %d", desc, tc.size, size))
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
			testSortAgents(t, tc.pm, page.Agents)
		})

	}
}

func TestUpdateAgent(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	validAgentName, err := types.NewIdentifier("group")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	ag, err := fleetService.CreateAgent(context.Background(), "token", fleet.Agent{
		Name:      validAgentName,
		AgentTags: map[string]string{"test": "true"},
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	validName, err := types.NewIdentifier("group")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	_, _ = fleetService.CreateAgentGroup(context.Background(), "token", fleet.AgentGroup{
		Name: validName,
		Tags: map[string]string{"test": "true"},
	})

	wrongAgentGroup := fleet.Agent{MFThingID: wrongID}
	cases := map[string]struct {
		group fleet.Agent
		token string
		err   error
	}{
		"update existing agent": {
			group: ag,
			token: token,
			err:   nil,
		},
		"update group with wrong credentials": {
			group: ag,
			token: invalidToken,
			err:   fleet.ErrUnauthorizedAccess,
		},
		"update a non-existing group": {
			group: wrongAgentGroup,
			token: token,
			err:   fleet.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := fleetService.EditAgent(context.Background(), tc.token, tc.group)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %d got %d", desc, tc.err, err))
		})
	}
}

func TestValidateAgent(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	nameID, err := types.NewIdentifier("eu-agents")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	validAgent := fleet.Agent{
		MFOwnerID: ownerID.String(),
		Name:      nameID,
		OrbTags:   make(map[string]string),
	}
	validAgent.OrbTags = map[string]string{
		"region":    "eu",
		"node_type": "dns",
	}
	cases := map[string]struct {
		agent fleet.Agent
		token string
		err   error
	}{
		"validate a valid agent": {
			agent: validAgent,
			token: token,
			err:   nil,
		},
		"validate a valid agent with an invalid token": {
			agent: validAgent,
			token: invalidToken,
			err:   fleet.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := fleetService.ValidateAgent(context.Background(), tc.token, tc.agent)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestCreateAgent(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})
	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	ownerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	nameID, err := types.NewIdentifier("eu-agents")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	conflictCase, err := createAgent(t, "agent", fleetService)

	validAgent := fleet.Agent{
		MFOwnerID: ownerID.String(),
		Name:      nameID,
		OrbTags:   make(map[string]string),
		Created:   time.Time{},
	}
	validAgent.OrbTags = map[string]string{
		"region":    "eu",
		"node_type": "dns",
	}
	cases := map[string]struct {
		agent fleet.Agent
		token string
		err   error
	}{
		"add a valid agent": {
			agent: validAgent,
			token: token,
			err:   nil,
		},
		"add a valid agent with an invalid token": {
			agent: validAgent,
			token: invalidToken,
			err:   fleet.ErrUnauthorizedAccess,
		},
		"add a conflict agent": {
			agent: conflictCase,
			token: token,
			err:   fleet.ErrConflict,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := fleetService.CreateAgent(context.Background(), tc.token, tc.agent)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestRemoveAgent(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	ag, err := createAgent(t, "my-agent", fleetService)
	require.Nil(t, err, fmt.Sprintf("unexpetec error: %s", err))

	cases := map[string]struct {
		id    string
		token string
		err   error
	}{
		"remove existing agent": {
			id:    ag.MFThingID,
			token: token,
			err:   nil,
		},
		"remove agent with wrong credentials": {
			id:    ag.MFThingID,
			token: invalidToken,
			err:   fleet.ErrUnauthorizedAccess,
		},
		"remove non-existing agent": {
			id:    wrongID,
			token: token,
			err:   nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := fleetService.RemoveAgent(context.Background(), tc.token, tc.id)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		})
	}
}

func TestListBackends(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	cases := map[string]struct {
		token string
		err   error
	}{
		"Retrieve a list of backends": {
			token: token,
			err:   nil,
		},
		"Retrieve a list of backends with a invalid token": {
			token: invalidToken,
			err:   fleet.ErrUnauthorizedAccess,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := fleetService.ListAgentBackends(context.Background(), tc.token)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestViewAgentBackend(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	cases := map[string]struct {
		name  string
		token string
		err   error
	}{
		"view backend not registered": {
			name:  "invalid",
			token: token,
			err:   errors.ErrNotFound,
		},
		"view backend with invalid token": {
			name:  "pktvisor",
			token: invalidToken,
			err:   errors.ErrUnauthorizedAccess,
		},
		"view registered backend": {
			name:  "pktvisor",
			token: token,
			err:   nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			_, err := fleetService.ViewAgentBackend(context.Background(), tc.token, tc.name)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestViewAgentInfoByChannelIDInternal(t *testing.T) {
	users := flmocks.NewAuthService(map[string]string{token: email})

	thingsServer := newThingsServer(newThingsService(users))
	fleetService := newService(users, thingsServer.URL)

	ag, err := createAgent(t, "agent", fleetService)

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := map[string]struct {
		channelID string
		agent     fleet.Agent
		err       error
	}{
		"view agent info by existent channelID": {
			channelID: ag.MFChannelID,
			agent:     ag,
			err:       nil,
		},
		"view agent info by non-existent channelID": {
			channelID: chID.String(),
			agent:     fleet.Agent{},
			err:       errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			agent, err := fleetService.ViewAgentInfoByChannelIDInternal(context.Background(), tc.channelID)
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
			assert.Equal(t, tc.agent, agent, fmt.Sprintf("%s: expected %s got %s", desc, tc.agent, agent))
		})
	}
}

func createAgent(t *testing.T, name string, svc fleet.Service) (fleet.Agent, error) {
	t.Helper()
	aCopy := agent
	validName, err := types.NewIdentifier(name)
	if err != nil {
		return fleet.Agent{}, err
	}
	aCopy.Name = validName
	ag, err := svc.CreateAgent(context.Background(), token, aCopy)
	if err != nil {
		return fleet.Agent{}, err
	}
	return ag, nil
}

func testSortAgents(t *testing.T, pm fleet.PageMetadata, ags []fleet.Agent) {
	t.Helper()
	switch pm.Order {
	case "name":
		current := ags[0]
		for _, res := range ags {
			if pm.Dir == "asc" {
				assert.GreaterOrEqual(t, res.Name.String(), current.Name.String())
			}
			if pm.Dir == "desc" {
				assert.GreaterOrEqual(t, current.Name.String(), res.Name.String())
			}
			current = res
		}
	default:
		break
	}
}
