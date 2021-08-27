// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package postgres_test

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/postgres"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAgentGroupSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentGroupRepository := postgres.NewAgentGroupRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-group")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	group := fleet.AgentGroup{
		Name:        nameID,
		MFOwnerID:   oID.String(),
		MFChannelID: chID.String(),
		Tags:        types.Tags{"testkey": "testvalue"},
	}

	cases := map[string]struct {
		agentGroup fleet.AgentGroup
		err        error
	}{
		"create new group": {
			agentGroup: group,
			err:        nil,
		},
		"create group that already exist": {
			agentGroup: group,
			err:        errors.ErrConflict,
		},
		"create group with invalid name": {
			agentGroup: fleet.AgentGroup{MFOwnerID: oID.String()},
			err:        errors.ErrMalformedEntity,
		},
		"create group with invalid owner ID": {
			agentGroup: fleet.AgentGroup{Name: nameID, MFOwnerID: "invalid"},
			err:        errors.ErrMalformedEntity,
		},
	}

	for desc, tc := range cases {
		_, err := agentGroupRepository.Save(context.Background(), tc.agentGroup)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected '%s' got '%s'", desc, tc.err, err))
	}

}

func TestAgentGroupRetrieve(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentGroupRepo := postgres.NewAgentGroupRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-group")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	group := fleet.AgentGroup{
		Name:        nameID,
		Description: "a example",
		MFOwnerID:   oID.String(),
		MFChannelID: chID.String(),
		Tags:        types.Tags{"testkey": "testvalue"},
	}

	id, err := agentGroupRepo.Save(context.Background(), group)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	cases := map[string]struct {
		groupID string
		ownerID string
		err     error
		tags    types.Tags
	}{
		"retrieve existing agent group by groupID and ownerID": {
			groupID: id,
			ownerID: group.MFOwnerID,
			tags:    types.Tags{"testkey": "testvalue"},
			err:     nil,
		},
		"retrieve non-existent agent group by groupID and ownerID": {
			ownerID: id,
			groupID: group.MFOwnerID,
			err:     errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		ag, err := agentGroupRepo.RetrieveByID(context.Background(), tc.groupID, tc.ownerID)
		if err == nil {
			assert.Equal(t, nameID, ag.Name, fmt.Sprintf("%s: expected %s got %s\n", desc, nameID, ag.Name))
		}
		if len(tc.tags) > 0 {
			assert.Equal(t, tc.tags, ag.Tags)
		}
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestMultiAgentGroupRetrieval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	agentGroupRepo := postgres.NewAgentGroupRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	wrongID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	n := uint64(10)
	for i := uint64(0); i < n; i++ {

		nameID, err := types.NewIdentifier(fmt.Sprintf("ue-agent-group-%d", i))
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		group := fleet.AgentGroup{
			Name:        nameID,
			Description: "a example",
			MFOwnerID:   oID.String(),
			MFChannelID: chID.String(),
			Tags:        types.Tags{"testkey": "testvalue"},
		}

		ag, err := agentGroupRepo.Save(context.Background(), group)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))
		fmt.Sprint(ag)
	}

	cases := map[string]struct {
		owner        string
		pageMetadata fleet.PageMetadata
		size         uint64
	}{
		"retrieve all groups with existing owner": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  n,
			},
			size: n,
		},
		"retrieve subset of agent groups with existing owner": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: n / 2,
				Limit:  n,
				Total:  n,
			},
			size: n / 2,
		},
		"retrieve agent groups with no-existing owner": {
			owner: wrongID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Total:  0,
			},
			size: 0,
		},
		"retrieve agent groups with no-existing name": {
			owner: oID.String(),
			pageMetadata: fleet.PageMetadata{
				Offset: 0,
				Limit:  n,
				Name:   "wrong",
				Total:  0,
			},
			size: 0,
		},
		"retrieve agent groups sorted by name ascendent": {
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
		page, err := agentGroupRepo.RetrieveAllAgentGroupsByOwner(context.Background(), tc.owner, tc.pageMetadata)
		require.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s\n", desc, err))
		size := uint64(len(page.AgentGroups))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected size %d got %d", desc, tc.size, size))
		assert.Equal(t, tc.pageMetadata.Total, page.Total, fmt.Sprintf("%s: expected total %d got %d", desc, tc.pageMetadata.Total, page.Total))

		if size > 0 {
			testSortAgentGroups(t, tc.pageMetadata, page.AgentGroups)
		}
	}
}

func TestAgentGroupUpdate(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	groupRepo := postgres.NewAgentGroupRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-group")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	invalideOwnerID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	invalideID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	group := fleet.AgentGroup{
		Name:        nameID,
		MFOwnerID:   oID.String(),
		MFChannelID: chID.String(),
		Tags:        types.Tags{"testkey": "testvalue"},
	}

	groupID, err := groupRepo.Save(context.Background(), group)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	group.ID = groupID

	cases := map[string]struct {
		group fleet.AgentGroup
		err   error
	}{
		"update a existing group": {
			group: group,
			err:   nil,
		},
		"update a non-existing group with a existing user": {
			group: fleet.AgentGroup{
				ID:        invalideID.String(),
				MFOwnerID: oID.String(),
			},
			err: fleet.ErrNotFound,
		},
		"update a existing group with a non-existing user": {
			group: fleet.AgentGroup{
				ID:        groupID,
				MFOwnerID: invalideOwnerID.String(),
			},
			err: fleet.ErrNotFound,
		},
		"update a non-existing group with a non-existing user": {
			group: fleet.AgentGroup{
				ID:        invalideID.String(),
				MFOwnerID: invalideOwnerID.String(),
			},
			err: fleet.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		_, err := groupRepo.Update(context.Background(), tc.group.MFOwnerID, tc.group)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestAgentGroupDelete(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	groupRepo := postgres.NewAgentGroupRepository(dbMiddleware, logger)

	oID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	invalidID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	chID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nameID, err := types.NewIdentifier("my-group")
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	group := fleet.AgentGroup{
		Name:        nameID,
		MFOwnerID:   oID.String(),
		MFChannelID: chID.String(),
		Tags:        types.Tags{"testkey": "testvalue"},
	}

	groupID, err := groupRepo.Save(context.Background(), group)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s\n", err))

	group.ID = groupID

	cases := map[string]struct {
		ID      string
		ownerID string
		err     error
	}{
		"remove a existing agent group": {
			ID:      group.ID,
			ownerID: group.MFOwnerID,
			err:     nil,
		},
		"remove a non-existing agent group": {
			ID:      invalidID.String(),
			ownerID: group.MFOwnerID,
			err:     nil,
		},
	}

	for desc, tc := range cases {
		err := groupRepo.Delete(context.Background(), tc.ID, tc.ownerID)
		require.Nil(t, err, fmt.Sprintf("%s: failed to remove agent group due to: %s", desc, err))

		_, err = groupRepo.RetrieveByID(context.Background(), tc.ID, tc.ownerID)
		require.True(t, errors.Contains(err, fleet.ErrNotFound), fmt.Sprintf("%s: expected %s got %s", desc, fleet.ErrNotFound, err))
	}
}

func testSortAgentGroups(t *testing.T, pm fleet.PageMetadata, ags []fleet.AgentGroup) {
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
