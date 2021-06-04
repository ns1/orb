// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package fleet

import (
	"encoding/json"
	"github.com/mainflux/mainflux/pkg/messaging"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	"go.uber.org/zap"
)

type AgentCommsService interface {
	// Start set up communication with the message bus to communicate with agents
	Start() error
	// Stop end communication with the message bus
	Stop() error
}

var _ AgentCommsService = (*fleetCommsService)(nil)

type fleetCommsService struct {
	logger *zap.Logger
	// agent comms
	agentPubSub mfnats.PubSub
}

func NewFleetCommsService(logger *zap.Logger, agentPubSub mfnats.PubSub) AgentCommsService {
	return &fleetCommsService{
		logger:      logger,
		agentPubSub: agentPubSub,
	}
}

func (svc fleetCommsService) handleMsgFromAgent(msg messaging.Message) error {
	var payload interface{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}
	svc.logger.Info("received message", zap.Any("payload", payload))
	return nil
}

func (svc fleetCommsService) Start() error {
	// TODO make this the agent channel
	if err := svc.agentPubSub.Subscribe(mfnats.SubjectAllChannels, svc.handleMsgFromAgent); err != nil {
		return err
	}
	svc.logger.Info("subscribed to agent info channels")
	return nil
}

func (svc fleetCommsService) Stop() error {
	// TODO make this the agent channel
	if err := svc.agentPubSub.Unsubscribe(mfnats.SubjectAllChannels); err != nil {
		return err
	}
	svc.logger.Info("subscribed to agent info channels")
	return nil
}
