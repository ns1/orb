// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package consumer

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/pkg/types"
	"go.uber.org/zap"
	"strings"
)

const (
	stream = "orb.policies"
	group  = "orb.fleet"

	datasetPrefix = "dataset."
	datasetCreate = datasetPrefix + "create"
	policyPrefix  = "policy."
	policyCreate  = policyPrefix + "create"
	policyUpdate  = policyPrefix + "update"
	policyRemove  = policyPrefix + "remove"

	exists = "BUSYGROUP Consumer Group name already exists"
)

type Subscriber interface {
	Subscribe(context context.Context) error
}

type eventStore struct {
	fleetService fleet.Service
	commsService fleet.AgentCommsService
	client       *redis.Client
	esconsumer   string
	logger       *zap.Logger
}

// NewEventStore returns new event store instance.
func NewEventStore(fleetService fleet.Service, commsService fleet.AgentCommsService, client *redis.Client, esconsumer string, log *zap.Logger) Subscriber {
	return eventStore{
		fleetService: fleetService,
		commsService: commsService,
		client:       client,
		esconsumer:   esconsumer,
		logger:       log,
	}
}

func (es eventStore) Subscribe(context context.Context) error {
	err := es.client.XGroupCreateMkStream(context, stream, group, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := es.client.XReadGroup(context, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: es.esconsumer,
			Streams:  []string{stream, ">"},
			Count:    100,
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}

		for _, msg := range streams[0].Messages {
			event := msg.Values

			var err error
			switch event["operation"] {
			case datasetCreate:
				rte := decodeDatasetCreate(event)
				err = es.handleDatasetCreate(context, rte)
			case policyUpdate:
				rte, derr := decodePolicyUpdate(event)
				if derr != nil {
					err = derr
					break
				}
				err = es.handlePolicyUpdate(context, rte)
			case policyRemove:
				rte := decodePolicyRemove(event)
				err = es.handlePolicyRemove(context, rte)
			}
			if err != nil {
				es.logger.Error("Failed to handle event", zap.String("operation", event["operation"].(string)), zap.Error(err))
				break
			}
			es.client.XAck(context, stream, group, msg.ID)
		}
	}
}

func decodeDatasetCreate(event map[string]interface{}) createDatasetEvent {
	return createDatasetEvent{
		id:           read(event, "id", ""),
		ownerID:      read(event, "owner_id", ""),
		name:         read(event, "name", ""),
		agentGroupID: read(event, "group_id", ""),
		policyID:     read(event, "policy_id", ""),
		sinkID:       read(event, "sink_id", ""),
	}
}

// the policy service is notifying that a new dataset has been created
// notify all agents in the AgentGroup specified in the dataset about the new agent policy
func (es eventStore) handleDatasetCreate(ctx context.Context, e createDatasetEvent) error {

	ag, err := es.fleetService.ViewAgentGroupByIDInternal(ctx, e.agentGroupID, e.ownerID)
	if err != nil {
		return err
	}

	return es.commsService.NotifyGroupNewDataset(ctx, ag, "", e.policyID, e.ownerID)
}

func decodePolicyUpdate(event map[string]interface{}) (updatePolicyEvent, error) {
	val := updatePolicyEvent{
		id:      read(event, "id", ""),
		ownerID: read(event, "owner_id", ""),
	}

	strgroups := read(event, "groups_ids", "")
	val.groupsIDs = strings.Split(strgroups, ",")

	strpolicy := read(event, "policy", "")
	var policy types.Metadata
	if err := json.Unmarshal([]byte(strpolicy), &policy); err != nil {
		return updatePolicyEvent{}, err
	}

	return val, nil
}

// the policy service is notifying that a policy has been updated
// notify all agents in the AgentGroup specified in the dataset about the policy update
func (es eventStore) handlePolicyUpdate(ctx context.Context, e updatePolicyEvent) error {
	for _, a := range e.groupsIDs {
		ag, err := es.fleetService.ViewAgentGroupByIDInternal(ctx, a, e.ownerID)
		if err != nil {
			return err
		}
		err = es.commsService.NotifyGroupNewDataset(ctx, ag, "", e.id, e.ownerID)
		if err != nil {
			return err
		}
	}
	return nil
}

func decodePolicyRemove(event map[string]interface{}) removePolicyEvent {
	val := removePolicyEvent{
		id:      read(event, "id", ""),
		ownerID: read(event, "owner_id", ""),
	}

	strgroups := read(event, "groups_ids", "")
	val.groupsIDs = strings.Split(strgroups, ",")
	return val
}

// the policy service is notifying that a policy has been removed
// notify all agents in the AgentGroup specified in the dataset about the policy removal
func (es eventStore) handlePolicyRemove(ctx context.Context, e removePolicyEvent) error {
	for _, a := range e.groupsIDs {
		ag, err := es.fleetService.ViewAgentGroupByIDInternal(ctx, a, e.ownerID)
		if err != nil {
			return err
		}
		err = es.commsService.NotifyPolicyRemoval(e.id, ag)
		if err != nil {
			return err
		}
	}
	return nil
}

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}

	return val
}
