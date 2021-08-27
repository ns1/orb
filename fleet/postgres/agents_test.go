// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/postgres"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"strings"
	"testing"
)

const maxNameSize = 1024

var (
	invalidName = strings.Repeat("m", maxNameSize+1)
	logger, _   = zap.NewDevelopment()
	wrongValue  = "wrong-value"
)

func TestAgentSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("myagent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	cases := []struct {
		desc  string
		agent fleet.Agent
		err   error
	}{
		{
			desc:  "create new agent",
			agent: agent,
			err:   nil,
		},
		{
			desc:  "create agent that already exist",
			agent: agent,
			err:   errors.ErrConflict,
		},
	}

	for _, tc := range cases {
		err := agentRepo.Save(context.Background(), tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", tc.desc, tc.err, err))
	}
}

func TestAgentRetrieve(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("myagent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		thingID   string
		channelID string
		err       error
		tags      types.Tags
	}{
		"retrieve existing agent by thingID and channelID": {
			thingID:   agent.MFThingID,
			channelID: agent.MFChannelID,
			tags:      types.Tags{"testkey": "testvalue"},
			err:       nil,
		},
		"retrieve non-existent agent by thingID and channelID": {
			thingID:   agent.MFOwnerID,
			channelID: agent.MFChannelID,
			err:       errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		ag, err := agentRepo.RetrieveByIDWithChannel(context.Background(), tc.thingID, tc.channelID)
		if err == nil {
			assert.Equal(t, nameID, ag.Name, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
		}
		if len(tc.tags) > 0 {
			assert.Equal(t, tc.tags, ag.OrbTags)
			assert.Equal(t, tc.tags, ag.AgentTags)
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestAgentUpdateData(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("myagent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		agent fleet.Agent
		err   error
	}{
		"update existing agent data by thingID and channelID": {
			agent: fleet.Agent{
				MFThingID:     thID.String(),
				MFChannelID:   chID.String(),
				AgentMetadata: types.Metadata{"newkey": "newvalue"},
			},
			err: nil,
		},
		"update non-existent agent data by thingID and channelID": {
			agent: fleet.Agent{
				MFThingID:     chID.String(),
				MFChannelID:   thID.String(),
				AgentMetadata: types.Metadata{"newkey": "newvalue"},
			},
			err: errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		err = agentRepo.UpdateDataByIDWithChannel(context.Background(), tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		if err == nil {
			ag, err := agentRepo.RetrieveByIDWithChannel(context.Background(), tc.agent.MFThingID, tc.agent.MFChannelID)
			assert.Nil(t, err)
			assert.Equal(t, tc.agent.AgentMetadata, ag.AgentMetadata, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
		}
	}
}

func TestAgentUpdateHeartbeat(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("myagent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:        nameID,
		MFThingID:   thID.String(),
		MFOwnerID:   oID.String(),
		MFChannelID: chID.String(),
		LastHBData:  types.Metadata{"heartbeatdata": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		agent fleet.Agent
		err   error
	}{
		"update existing agent heartbeat by thingID and channelID": {
			agent: fleet.Agent{
				MFThingID:   thID.String(),
				MFChannelID: chID.String(),
				LastHBData:  types.Metadata{"heartbeatdata2": "newvalue"},
			},
			err: nil,
		},
		"update non-existent agent heart beat by thingID and channelID": {
			agent: fleet.Agent{
				MFThingID:   chID.String(),
				MFChannelID: thID.String(),
				LastHBData:  types.Metadata{"heartbeatdata2": "newvalue"},
			},
			err: errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		err = agentRepo.UpdateHeartbeatByIDWithChannel(context.Background(), tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		if err == nil {
			ag, err := agentRepo.RetrieveByIDWithChannel(context.Background(), tc.agent.MFThingID, tc.agent.MFChannelID)
			assert.Nil(t, err)
			assert.Equal(t, tc.agent.LastHBData, ag.LastHBData, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
		}
	}
}

func TestMultiAgentRetrieval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	wrongoID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	name := "agent_name"
	metaStr := `{"field1":"value1","field2":{"subfield11":"value2","subfield12":{"subfield121":"value3","subfield122":"value4"}}}`
	subMetaStr := `{"field2":{"subfield12":{"subfield121":"value3"}}}`
	tagsStr := `{"region": "EU", "node_type": "dns"}`
	subTagsStr := `{"region": "EU"}`

	metadata := types.Metadata{}
	json.Unmarshal([]byte(metaStr), &metadata)

	subMeta := types.Metadata{}
	json.Unmarshal([]byte(subMetaStr), &subMeta)

	tags := types.Tags{}
	json.Unmarshal([]byte(tagsStr), &tags)

	subTags := types.Tags{}
	json.Unmarshal([]byte(subTagsStr), &subTags)

	wrongMeta := types.Metadata{
		"field": "value1",
	}
	wrongTags := types.Tags{
		"field": "value1",
	}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		thID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		chID, err := uuid.NewV4()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		th := fleet.Agent{
			MFOwnerID:   oID.String(),
			MFThingID:   thID.String(),
			MFChannelID: chID.String(),
		}

		th.Name, err = types.NewIdentifier(fmt.Sprintf("%s-%d", name, i))
		require.True(t, th.Name.IsValid(), "invalid Identifier name: %s")
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		th.AgentMetadata = metadata
		th.AgentTags = tags
		th.OrbTags = subTags

		err = agentRepo.Save(context.Background(), th)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
	}

	cases := map[string]struct {
		owner        string
		pageMetadata fleet.PageMetadata
		size         uint64
	}{
		"retrieve all agents with existing owner": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
			},
			size: n,
		},
		"retrieve subset of agents with existing owner": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: n / 2,
				Limit:  n,
				Total:  n,
			},
			size: n / 2,
		},
		"retrieve agents with existing metadata": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset:   0,
				Limit:    n,
				Total:    n,
				Metadata: subMeta,
			},
			size: n,
		},
		"retrieve agents with existing tags": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Tags:   subTags,
			},
			size: n,
		},
		"retrieve agents with non-existing owner": {
			owner: wrongoID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  0,
			},
			size: 0,
		},
		"retrieve agents with non-existing name": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Name:   "wrong",
				Total:  0,
			},
			size: 0,
		},
		"retrieve agents with non-existing metadata": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset:   0,
				Limit:    n,
				Total:    0,
				Metadata: wrongMeta,
			},
			size: 0,
		},
		"retrieve agents with non-existing tags": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  0,
				Tags:   wrongTags,
			},
			size: 0,
		},
		"retrieve agents sorted by name ascendent": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Order:  "name",
				Dir:    "asc",
			},
			size: n,
		},
		"retrieve agents sorted by name descendent": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
				Order:  "name",
				Dir:    "desc",
			},
			size: n,
		},
	}

	for desc, tc := range cases {
		page, err := agentRepo.RetrieveAll(context.Background(), tc.owner, tc.pageMetadata)
		size := uint64(len(page.Agents))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected size %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.pageMetadata.Total, page.Total, fmt.Sprintf("%s: expected total %d got %d\n", desc, tc.pageMetadata.Total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))

		// Check if Agents list have been sorted properly
		if size > 0 {
			testSortAgents(t, tc.pageMetadata, page.Agents)
		}
	}
}

func TestAgentUpdate(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-agent1")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	updatedNameID, err := types.NewIdentifier("my-agent2")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	agent := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), agent)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		agent fleet.Agent
		err   error
	}{
		"update existing agent data by thingID": {
			agent: fleet.Agent{
				MFThingID: thID.String(),
				MFOwnerID: oID.String(),
				Name:      updatedNameID,
				OrbTags:   types.Tags{"newkey": "newvalue"},
			},
			err: nil,
		},
		"update non-existent agent data by thingID": {
			agent: fleet.Agent{
				MFThingID: oID.String(),
				MFOwnerID: oID.String(),
				Name:      updatedNameID,
				OrbTags:   types.Tags{"newkey": "newvalue"},
			},
			err: errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		err = agentRepo.UpdateAgentByID(context.Background(), tc.agent.MFOwnerID, tc.agent)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
		if err == nil {
			ag, err := agentRepo.RetrieveByID(context.Background(), tc.agent.MFOwnerID, tc.agent.MFThingID)
			assert.Nil(t, err)
			assert.Equal(t, tc.agent.OrbTags, ag.OrbTags, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
		}
	}
}

func TestDeleteAgent(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentRepo := postgres.NewAgentRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	thID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	invalidID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-agent")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	ag := fleet.Agent{
		Name:          nameID,
		MFThingID:     thID.String(),
		MFOwnerID:     oID.String(),
		MFChannelID:   chID.String(),
		OrbTags:       types.Tags{"testkey": "testvalue"},
		AgentTags:     types.Tags{"testkey": "testvalue"},
		AgentMetadata: types.Metadata{"testkey": "testvalue"},
	}

	err = agentRepo.Save(context.Background(), ag)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		ID      string
		ownerID string
		err     error
	}{
		"remove existing agent": {
			ID:      ag.MFThingID,
			ownerID: ag.MFOwnerID,
			err:     nil,
		},
		"remove a non-existing agent": {
			ID:      invalidID.String(),
			ownerID: ag.MFOwnerID,
			err:     nil,
		},
	}

	for desc, tc := range cases {
		err := agentRepo.Delete(context.Background(), tc.ownerID, tc.ID)
		require.Nil(t, err, fmt.Sprintf("%s: failed to remove agent due to: %s", desc, err))
	}
}

func testSortAgents(t *testing.T, pm fleet.PageMetadata, ths []fleet.Agent) {
	switch pm.Order {
	case "name":
		current := ths[0]
		for _, res := range ths {
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
