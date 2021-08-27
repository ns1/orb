// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package producer

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
)

const (
	streamID  = "orb.fleet"
	streamLen = 1000
)

var _ fleet.Service = (*eventStore)(nil)

type eventStore struct {
	svc    fleet.Service
	client *redis.Client
	logger *zap.Logger
}

func (es eventStore) ViewAgentByID(ctx context.Context, token string, thingID string) (fleet.Agent, error) {
	return es.svc.ViewAgentByID(ctx, token, thingID)
}

func (es eventStore) EditAgent(ctx context.Context, token string, agent fleet.Agent) (fleet.Agent, error) {
	return es.svc.EditAgent(ctx, token, agent)
}

func (es eventStore) ViewAgentGroupByIDInternal(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	return es.svc.ViewAgentGroupByIDInternal(ctx, groupID, ownerID)
}

func (es eventStore) ViewAgentGroupByID(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	return es.svc.ViewAgentGroupByID(ctx, groupID, ownerID)
}

func (es eventStore) ListAgentGroups(ctx context.Context, token string, pm fleet.PageMetadata) (fleet.PageAgentGroup, error) {
	return es.svc.ListAgentGroups(ctx, token, pm)
}

func (es eventStore) EditAgentGroup(ctx context.Context, token string, ag fleet.AgentGroup) (fleet.AgentGroup, error) {
	return es.svc.EditAgentGroup(ctx, token, ag)
}

func (es eventStore) ListAgents(ctx context.Context, token string, pm fleet.PageMetadata) (fleet.Page, error) {
	return es.svc.ListAgents(ctx, token, pm)
}

func (es eventStore) CreateAgent(ctx context.Context, token string, a fleet.Agent) (fleet.Agent, error) {
	return es.svc.CreateAgent(ctx, token, a)
}

func (es eventStore) CreateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (fleet.AgentGroup, error) {
	return es.svc.CreateAgentGroup(ctx, token, s)
}

func (es eventStore) RemoveAgentGroup(ctx context.Context, token string, groupID string) (err error) {
	err = es.svc.RemoveAgentGroup(ctx, token, groupID)
	if err != nil {
		return err
	}

	event := removeAgentGroupEvent{
		groupID: groupID,
		token:   token,
	}
	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       event.encode(),
	}
	err = es.client.XAdd(ctx, record).Err()
	if err != nil {
		es.logger.Error("error sending event to event store", zap.Error(err))
		return err
	}

	return nil

}

func (es eventStore) ValidateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (fleet.AgentGroup, error) {
	return es.svc.ValidateAgentGroup(ctx, token, s)
}

func (es eventStore) RemoveAgent(ctx context.Context, token, thingID string) (err error) {
	return es.svc.RemoveAgent(ctx, token, thingID)
}

// NewEventStoreMiddleware returns wrapper around fleet service that sends
// events to event store.
func NewEventStoreMiddleware(svc fleet.Service, client *redis.Client) fleet.Service {
	return eventStore{
		svc:    svc,
		client: client,
	}
}
