// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"context"
	"github.com/mainflux/mainflux"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	"github.com/ns1labs/orb/pkg/errors"
	"go.uber.org/zap"
)

var (
	ErrCreateAgent = errors.New("failed to create agent")

	// ErrThings indicates failure to communicate with Mainflux Things service.
	// It can be due to networking error or invalid/unauthorized request.
	ErrThings = errors.New("failed to receive response from Things service")

	errCreateThing   = errors.New("failed to create thing")
	errThingNotFound = errors.New("thing not found")
)

func (svc fleetService) ViewAgentByID(ctx context.Context, token string, thingID string) (Agent, error) {
	ownerID, err := svc.identify(token)
	if err != nil {
		return Agent{}, err
	}
	return svc.agentRepo.RetrieveByID(ctx, ownerID, thingID)
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

	err = svc.agentComms.NotifyAgentGroupMembership(res)
	if err != nil {
		svc.logger.Error("failure during agent group membership comms", zap.Error(err))
	}

	return res, nil
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
		err = errors.Wrap(err, errT)
	}

	if errT := svc.mfsdk.DeleteChannel(res.MFChannelID, token); errT != nil {
		err = errors.Wrap(err, errT)
	}

	err = svc.agentRepo.Delete(ctx, ownerID, thingID)
	if err != nil {
		return err
	}

	return nil
}
