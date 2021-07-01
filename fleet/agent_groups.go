/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

type AgentGroup struct {
	ID          string
	Name        types.Identifier
	MFOwnerID   string
	MFChannelID string
	Tags        types.Tags
	Created     time.Time
}

type AgentGroupService interface {
	// CreateAgentGroup creates new AgentGroup, associated channel, applies to Agents as appropriate
	CreateAgentGroup(ctx context.Context, token string, s AgentGroup) (AgentGroup, error)
}

type AgentGroupRepository interface {
	// Save persists the AgentGroup. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, group AgentGroup) (string, error)
	// RetrieveAllByAgent get all AgentGroup which an Agent belongs to.
	RetrieveAllByAgent(ctx context.Context, a Agent) ([]AgentGroup, error)
}
