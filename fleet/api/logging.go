/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package api

import (
	"context"
	"github.com/ns1labs/orb/fleet"
	"go.uber.org/zap"
	"time"
)

var _ fleet.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *zap.Logger
	svc    fleet.Service
}

func (l loggingMiddleware) ViewAgentByID(ctx context.Context, token string, thingID string) (_ fleet.Agent, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: view_agent_by_id",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: view_agent_by_id",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ViewAgentByID(ctx, token, thingID)
}

func (l loggingMiddleware) EditAgent(ctx context.Context, token string, agent fleet.Agent) (_ fleet.Agent, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: edit_agent_by_id",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: edit_agent_by_id",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.EditAgent(ctx, token, agent)
}

func (l loggingMiddleware) ViewAgentGroupByIDInternal(ctx context.Context, groupID string, ownerID string) (_ fleet.AgentGroup, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: view_agent_group_by_id_internal",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: view_agent_group_by_id_internal",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ViewAgentGroupByIDInternal(ctx, groupID, ownerID)
}

func (l loggingMiddleware) ViewAgentGroupByID(ctx context.Context, groupID string, ownerID string) (_ fleet.AgentGroup, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: view_agent_group_by_id",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: view_agent_group_by_id",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ViewAgentGroupByID(ctx, groupID, ownerID)
}

func (l loggingMiddleware) ListAgentGroups(ctx context.Context, token string, pm fleet.PageMetadata) (_ fleet.PageAgentGroup, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: list_agent_groups",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: list_agent_groups",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ListAgentGroups(ctx, token, pm)
}

func (l loggingMiddleware) EditAgentGroup(ctx context.Context, token string, ag fleet.AgentGroup) (_ fleet.AgentGroup, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: edit_agent_groups",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: edit_agent_groups",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.EditAgentGroup(ctx, token, ag)
}

func (l loggingMiddleware) ListAgents(ctx context.Context, token string, pm fleet.PageMetadata) (_ fleet.Page, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: list_agents",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: list_agents",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ListAgents(ctx, token, pm)
}

func (l loggingMiddleware) CreateAgent(ctx context.Context, token string, a fleet.Agent) (_ fleet.Agent, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: create_agent",
				zap.String("name", a.Name.String()),
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: create_agent",
				zap.String("name", a.Name.String()),
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.CreateAgent(ctx, token, a)
}

func (l loggingMiddleware) CreateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (_ fleet.AgentGroup, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: create_agent_group",
				zap.String("name", s.Name.String()),
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: create_agent_group",
				zap.String("name", s.Name.String()),
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.CreateAgentGroup(ctx, token, s)
}

func (l loggingMiddleware) RemoveAgentGroup(ctx context.Context, token, groupID string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: delete_agent_groups",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: delete_agent_groups",
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.RemoveAgentGroup(ctx, token, groupID)
}

func (l loggingMiddleware) ValidateAgentGroup(ctx context.Context, token string, s fleet.AgentGroup) (_ fleet.AgentGroup, err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: validate_agent_group",
				zap.String("name", s.Name.String()),
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: validate_agent_group",
				zap.String("name", s.Name.String()),
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())
	return l.svc.ValidateAgentGroup(ctx, token, s)
}

func (l loggingMiddleware) RemoveAgent(ctx context.Context, token, thingID string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			l.logger.Warn("method call: delete_agent",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		} else {
			l.logger.Info("method call: delete_agent",
				zap.Error(err),
				zap.Duration("duration", time.Since(begin)))
		}
	}(time.Now())

	return l.svc.RemoveAgent(ctx, token, thingID)
}

func NewLoggingMiddleware(svc fleet.Service, logger *zap.Logger) fleet.Service {
	return &loggingMiddleware{logger, svc}
}
