/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"github.com/ns1labs/orb/fleet"
)

var _ fleet.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     fleet.Service
}

func (m metricsMiddleware) ViewAgentByID(ctx context.Context, token string, thingID string) (fleet.Agent, error) {
	return m.svc.ViewAgentByID(ctx, token, thingID)
}

func (m metricsMiddleware) EditAgent(ctx context.Context, token string, agent fleet.Agent) (fleet.Agent, error) {
	return m.svc.EditAgent(ctx, token, agent)
}

func (m metricsMiddleware) ViewAgentGroupByIDInternal(ctx context.Context, groupID string, ownerID string) (fleet.AgentGroup, error) {
	return m.svc.ViewAgentGroupByIDInternal(ctx, groupID, ownerID)
}

func (m metricsMiddleware) ViewAgentGroupByID(ctx context.Context, token string, groupID string) (fleet.AgentGroup, error) {
	return m.svc.ViewAgentGroupByID(ctx, token, groupID)
}

func (m metricsMiddleware) ListAgentGroups(ctx context.Context, token string, pm fleet.PageMetadata) (fleet.PageAgentGroup, error) {
	return m.svc.ListAgentGroups(ctx, token, pm)
}

func (m metricsMiddleware) EditAgentGroup(ctx context.Context, token string, ag fleet.AgentGroup) (fleet.AgentGroup, error) {
	return m.svc.EditAgentGroup(ctx, token, ag)
}

func (m metricsMiddleware) ListAgents(ctx context.Context, token string, pm fleet.PageMetadata) (fleet.Page, error) {
	return m.svc.ListAgents(ctx, token, pm)
}

func (m metricsMiddleware) CreateAgent(ctx context.Context, token string, a fleet.Agent) (fleet.Agent, error) {
	return m.svc.CreateAgent(ctx, token, a)
}

func (m metricsMiddleware) CreateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (fleet.AgentGroup, error) {
	return m.svc.CreateAgentGroup(ctx, token, s)
}

func (m metricsMiddleware) RemoveAgentGroup(ctx context.Context, token string, groupID string) error {
	return m.svc.RemoveAgentGroup(ctx, token, groupID)
}

func (m metricsMiddleware) ValidateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (fleet.AgentGroup, error) {
	return m.svc.ValidateAgentGroup(ctx, token, s)
}

func (m metricsMiddleware) RemoveAgent(ctx context.Context, token string, thingID string) error {
	return m.svc.RemoveAgent(ctx, token, thingID)
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc fleet.Service, counter metrics.Counter, latency metrics.Histogram) fleet.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}
