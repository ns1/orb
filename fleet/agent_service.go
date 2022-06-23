// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"encoding/json"
	"github.com/mainflux/mainflux"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/ns1labs/orb/fleet/backend"
	"github.com/ns1labs/orb/pkg/errors"
	"go.uber.org/zap"
	"strings"
)

var (
	ErrCreateAgent = errors.New("failed to create agent")

	// ErrThings indicates failure to communicate with Mainflux Things service.
	// It can be due to networking error or invalid/unauthorized request.
	ErrThings = errors.New("failed to receive response from Things service")

	errCreateThing   = errors.New("failed to create thing")
	errThingNotFound = errors.New("thing not found")
)

func (svc fleetService) addAgentToAgentGroupChannels(token string, a Agent) error {
	groupList, err := svc.agentGroupRepository.RetrieveAllByAgent(context.Background(), a)
	if err != nil {
		return err
	}

	if len(groupList) == 0 {
		return nil
	}

	var idList = make([]string, 1)
	idList[0] = a.MFThingID
	for _, group := range groupList {
		ids := mfsdk.ConnectionIDs{
			ChannelIDs: []string{group.MFChannelID},
			ThingIDs:   idList,
		}
		err = svc.mfsdk.Connect(ids, token)
		if err != nil {
			if strings.Contains(err.Error(), "409") {
				svc.logger.Warn("agent already connected, skipping...")
			} else {
				return err
			}
		}
	}

	return nil
}

func (svc fleetService) ViewAgentByID(ctx context.Context, token string, thingID string) (Agent, error) {
	ownerID, err := svc.identify(token)
	if err != nil {
		return Agent{}, err
	}
	return svc.agentRepo.RetrieveByID(ctx, ownerID, thingID)
}

func (svc fleetService) ViewAgentMatchingGroupsByID(ctx context.Context, token string, thingID string) (MatchingGroups, error) {
	ownerID, err := svc.identify(token)
	if err != nil {
		return MatchingGroups{}, err
	}

	matchingGroups, err := svc.agentGroupRepository.RetrieveMatchingGroups(ctx, ownerID, thingID)
	if err != nil {
		return MatchingGroups{}, err
	}

	return matchingGroups, nil
}

func (svc fleetService) ResetAgent(ctx context.Context, token string, agentID string) error {
	ownerID, err := svc.identify(token)
	if err != nil {
		return err
	}

	agent, err := svc.agentRepo.RetrieveByID(ctx, ownerID, agentID)
	if err != nil {
		return err
	}

	return svc.agentComms.NotifyAgentReset(agent, true, "Reset initiated from control plane")
}

func (svc fleetService) ViewAgentByIDInternal(ctx context.Context, ownerID string, id string) (Agent, error) {
	return svc.agentRepo.RetrieveByID(ctx, ownerID, id)
}

func (svc fleetService) ListAgents(ctx context.Context, token string, pm PageMetadata) (Page, error) {
	res, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return Page{}, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	return svc.agentRepo.RetrieveAll(ctx, res.GetId(), pm)
}

func (svc fleetService) CreateAgent(ctx context.Context, token string, a Agent) (Agent, error) {
	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return Agent{}, err
	}

	a.MFOwnerID = mfOwnerID

	md := map[string]interface{}{"type": "orb_agent"}

	// create new Thing
	mfThing, err := svc.thing(token, "", a.Name.String(), md)
	if err != nil {
		return Agent{}, errors.Wrap(ErrCreateAgent, err)
	}

	a.MFThingID = mfThing.ID
	a.MFKeyID = mfThing.Key

	// create main Agent RPC Channel
	mfChannelID, err := svc.mfsdk.CreateChannel(mfsdk.Channel{
		Name:     a.Name.String(),
		Metadata: md,
	}, token)
	if err != nil {
		if errT := svc.mfsdk.DeleteThing(mfThing.ID, token); errT != nil {
			err = errors.Wrap(err, errT)
		}
		return Agent{}, errors.Wrap(ErrCreateAgent, err)
	}

	a.MFChannelID = mfChannelID

	// RPC Channel to Agent
	err = svc.mfsdk.Connect(mfsdk.ConnectionIDs{
		ChannelIDs: []string{mfChannelID},
		ThingIDs:   []string{mfThing.ID},
	}, token)
	if err != nil {
		if errT := svc.mfsdk.DeleteThing(mfThing.ID, token); errT != nil {
			err = errors.Wrap(err, errT)
			// fall through
		}
		if errT := svc.mfsdk.DeleteChannel(mfChannelID, token); errT != nil {
			err = errors.Wrap(err, errT)
		}
		return Agent{}, errors.Wrap(ErrCreateAgent, err)
	}

	err = svc.agentRepo.Save(ctx, a)
	if err != nil {
		if errT := svc.mfsdk.DeleteThing(mfThing.ID, token); errT != nil {
			err = errors.Wrap(err, errT)
			// fall through
		}
		if errT := svc.mfsdk.DeleteChannel(mfChannelID, token); errT != nil {
			err = errors.Wrap(err, errT)
		}
		return Agent{}, errors.Wrap(ErrCreateAgent, err)
	}

	err = svc.addAgentToAgentGroupChannels(token, a)
	if err != nil {
		// TODO should we roll back?
		svc.logger.Error("failed to add agent to a existing group channel", zap.String("agent_id", a.MFThingID), zap.Error(err))
	}

	return a, nil
}

func (svc fleetService) EditAgent(ctx context.Context, token string, agent Agent) (Agent, error) {
	ownerID, err := svc.identify(token)
	if err != nil {
		return Agent{}, err
	}
	agent.MFOwnerID = ownerID

	err = svc.agentRepo.UpdateAgentByID(ctx, ownerID, agent)
	if err != nil {
		return Agent{}, err
	}

	res, err := svc.agentRepo.RetrieveByID(ctx, ownerID, agent.MFThingID)
	if err != nil {
		return Agent{}, err
	}

	err = svc.addAgentToAgentGroupChannels(token, res)
	if err != nil {
		// TODO should we roll back?
		svc.logger.Error("failed to add agent to a existing group channel", zap.String("agent_id", res.MFThingID), zap.Error(err))
	}

	err = svc.agentComms.NotifyAgentGroupMemberships(res)
	if err != nil {
		svc.logger.Error("failure during agent group membership comms", zap.Error(err))
	}

	return res, nil
}

func (svc fleetService) ValidateAgent(ctx context.Context, token string, a Agent) (Agent, error) {
	mfOwnerID, err := svc.identify(token)
	if err != nil {
		return Agent{}, err
	}

	a.MFOwnerID = mfOwnerID

	return a, nil
}

func (svc fleetService) RemoveAgent(ctx context.Context, token, thingID string) error {
	ownerID, err := svc.identify(token)
	if err != nil {
		return err
	}

	res, err := svc.agentRepo.RetrieveByID(ctx, ownerID, thingID)
	if err != nil {
		return nil
	}

	if errT := svc.mfsdk.DeleteThing(res.MFThingID, token); errT != nil {
		svc.logger.Error("failed to delete thing", zap.Error(errT), zap.String("thing_id", res.MFThingID))
	}

	if errT := svc.mfsdk.DeleteChannel(res.MFChannelID, token); errT != nil {
		svc.logger.Error("failed to delete channel", zap.Error(errT), zap.String("channel_id", res.MFChannelID))
	}

	err = svc.agentRepo.Delete(ctx, ownerID, thingID)
	if err != nil {
		return err
	}

	return nil
}

func (svc fleetService) ListAgentBackends(ctx context.Context, token string) ([]string, error) {
	_, err := svc.identify(token)
	if err != nil {
		return nil, err
	}
	return backend.GetList(), nil
}

func (svc fleetService) ViewAgentBackend(ctx context.Context, token string, name string) (interface{}, error) {
	_, err := svc.identify(token)
	if err != nil {
		return nil, err
	}
	if backend.HaveBackend(name) {
		return backend.GetBackend(name).Metadata(), nil
	}
	return nil, errors.ErrNotFound
}

func (svc fleetService) ViewAgentInfoByChannelIDInternal(ctx context.Context, channelID string) (Agent, error) {
	res, err := svc.agentRepo.RetrieveAgentInfoByChannelID(ctx, channelID)
	if err != nil {
		return Agent{}, err
	}
	return res, nil
}

func (svc fleetService) GetPolicyState(ctx context.Context, agent Agent) (map[string]interface{}, error) {

	jsonHb, err := json.Marshal(agent.LastHBData)
	if err != nil {
		svc.logger.Error("failed to marshal heartbeat data", zap.Error(err))
		return nil, err
	}
	var hb Heartbeat
	if err = json.Unmarshal(jsonHb, &hb); err != nil {
		svc.logger.Error("failed to unmarshal heartbeat data", zap.Error(err))
		return nil, err
	}

	policyState := make(map[string]interface{})
	for policyID, policyInfo := range hb.PolicyState {
		formattedPolicyInfo := policyInfo
		formattedPolicyInfo.Datasets = []string{}

		policyState[policyID] = formattedPolicyInfo
	}

	return policyState, nil
}
