/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"database/sql/driver"
	"github.com/ns1labs/orb/pkg/types"
	"time"
)

const (
	New State = iota
	Online
	Offline
	Stale
	Removed
	UpgradeRequired
)

type State int

var stateMap = [...]string{
	"new",
	"online",
	"offline",
	"stale",
	"removed",
	"upgrade_required",
}

var stateRevMap = map[string]State{
	"new":              New,
	"online":           Online,
	"offline":          Offline,
	"stale":            Stale,
	"removed":          Removed,
	"upgrade_required": UpgradeRequired,
}

func (s State) String() string {
	return stateMap[s]
}

func (s *State) Scan(value interface{}) error { *s = stateRevMap[string(value.([]byte))]; return nil }
func (s State) Value() (driver.Value, error)  { return s.String(), nil }

type Agent struct {
	Name           types.Identifier
	MFOwnerID      string
	MFThingID      string
	MFKeyID        string
	MFChannelID    string
	Created        time.Time
	OrbTags        types.Tags
	AgentTags      types.Tags
	AgentMetadata  types.Metadata
	State          State
	LastHBData     types.Metadata
	LastHB         time.Time
	MatchingGroups types.Metadata
}

// Page contains page related metadata as well as list of agents that
// belong to this page.
type Page struct {
	PageMetadata
	Agents []Agent
}

type Group struct {
	GroupID   string
	GroupName types.Identifier
}

type MatchingGroups struct {
	OwnerID string
	Groups  []Group
}

// AgentService Agent CRUD interface
type AgentService interface {
	// CreateAgent creates new agent
	CreateAgent(ctx context.Context, token string, a Agent) (Agent, error)
	// ViewAgentByID retrieves a Agent by provided thingID
	ViewAgentByID(ctx context.Context, token string, thingID string) (Agent, error)
	// ViewAgentMatchingGroupsByID Groups this Agent currently belongs to, according to matching agent and group tags
	ViewAgentMatchingGroupsByID(ctx context.Context, token string, thingID string) (MatchingGroups, error)
	// ViewAgentByIDInternal retrieves a Agent by provided thingID
	ViewAgentByIDInternal(ctx context.Context, ownerID string, thingID string) (Agent, error)
	// ListAgents retrieves data about subset of agents that belongs to the
	// user identified by the provided key.
	ListAgents(ctx context.Context, token string, pm PageMetadata) (Page, error)
	// EditAgent edit a Agent by provided thingID
	EditAgent(ctx context.Context, token string, agent Agent) (Agent, error)
	// ValidateAgent validates agent
	ValidateAgent(ctx context.Context, token string, a Agent) (Agent, error)
	// RemoveAgent removes an existing agent by owner and id
	RemoveAgent(ctx context.Context, token string, thingID string) error
	// ListAgentBackends List the available backends from fleet agents
	ListAgentBackends(ctx context.Context, token string) ([]string, error)
	// ViewAgentBackend retrieves a Backend by provided backend name
	ViewAgentBackend(ctx context.Context, token string, name string) (interface{}, error)
	//ViewAgentInfoByChannelIDInternal return a correspondent ownerID, name and agent tags by a provided channel id
	ViewAgentInfoByChannelIDInternal(ctx context.Context, channelID string) (Agent, error)
	// ResetAgent reset a agent on edge by a provided agent
	ResetAgent(ct context.Context, token string, agentID string) error
}

type AgentRepository interface {
	AgentHeartbeatRepository // may move this out so it can be in e.g. redis

	// Save persists the Agent. Successful operation is indicated by non-nil
	// error response.
	Save(ctx context.Context, agent Agent) error
	// UpdateDataByIDWithChannel update the tags and metadata for the Agent having the provided ID and owner
	UpdateDataByIDWithChannel(ctx context.Context, agent Agent) error
	// RetrieveByIDWithChannel retrieves the Agent having the provided ID and channelID access (i.e. from a Message)
	RetrieveByIDWithChannel(ctx context.Context, thingID string, channelID string) (Agent, error)
	// RetrieveAll retrieves the subset of Agents owned by the specified user
	RetrieveAll(ctx context.Context, owner string, pm PageMetadata) (Page, error)
	// RetrieveAllByAgentGroupID retrieves Agents in the specified group
	RetrieveAllByAgentGroupID(ctx context.Context, owner string, agentGroupID string, onlinishOnly bool) ([]Agent, error)
	// RetrieveMatchingAgents retrieve the matching agents by tags
	RetrieveMatchingAgents(ctx context.Context, owner string, tags types.Tags) (types.Metadata, error)
	// UpdateAgentByID update the the tags and name for the Agent having provided ID and owner
	UpdateAgentByID(ctx context.Context, ownerID string, agent Agent) error
	// RetrieveByID retrieves the Agent having the provided ID and owner
	RetrieveByID(ctx context.Context, ownerID string, thingID string) (Agent, error)
	// Delete an existing agent by owner and id
	Delete(ctx context.Context, ownerID string, thingID string) error
	// RetrieveAgentMetadataByOwner retrieves the Metadata having the OwnerID
	RetrieveAgentMetadataByOwner(ctx context.Context, ownerID string) ([]types.Metadata, error)
	// SetStaleStatus change status to stale according provided duration without heartbeats
	SetStaleStatus(ctx context.Context, minutes time.Duration) (int64, error)
	// RetrieveAgentInfoByChannelID gRPC version to retrieve ownerID, name and agent tags by a provided channelID
	RetrieveAgentInfoByChannelID(ctx context.Context, channelID string) (Agent, error)
}

type AgentHeartbeatRepository interface {
	// UpdateHeartbeatByIDWithChannel update the heartbeat data for the Agent having the provided ID and owner
	UpdateHeartbeatByIDWithChannel(ctx context.Context, agent Agent) error
}
