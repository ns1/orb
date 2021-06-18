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
	"fmt"
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

const CapabilitiesChannel = "agent"
const HeartbeatsChannel = "hb"
const RPCToCoreChannel = "tocore"
const RPCFromCoreChannel = "fromcore"
const LogChannel = "log"

type fleetCommsService struct {
	logger    *zap.Logger
	agentRepo AgentRepository

	// agent comms
	agentPubSub mfnats.PubSub
}

func NewFleetCommsService(logger *zap.Logger, agentRepo AgentRepository, agentPubSub mfnats.PubSub) AgentCommsService {
	return &fleetCommsService{
		logger:      logger,
		agentRepo:   agentRepo,
		agentPubSub: agentPubSub,
	}
}

func (svc fleetCommsService) handleCapabilities(thingID string, channelID string, payload map[string]interface{}) error {
	agent := Agent{MFThingID: thingID, MFChannelID: channelID}
	agent.AgentMetadata = payload
	err := svc.agentRepo.UpdateDataByIDWithChannel(context.Background(), agent)
	if err != nil {
		return err
	}
	return nil
}

func (svc fleetCommsService) handleHeartbeat(thingID string, channelID string, payload map[string]interface{}) error {
	agent := Agent{MFThingID: thingID, MFChannelID: channelID}
	agent.LastHBData = payload
	err := svc.agentRepo.UpdateHeartbeatByIDWithChannel(context.Background(), agent)
	if err != nil {
		return err
	}
	return nil
}

func (svc fleetCommsService) handleMsgFromAgent(msg messaging.Message) error {

	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}

	svc.logger.Debug("received agent message",
		zap.Any("payload", payload),
		zap.String("subtopic", msg.Subtopic),
		zap.String("channel", msg.Channel),
		zap.String("protocol", msg.Protocol),
		zap.Int64("created", msg.Created),
		zap.String("publisher", msg.Publisher))

	// dispatch
	switch msg.Subtopic {
	case CapabilitiesChannel:
		if err := svc.handleCapabilities(msg.Publisher, msg.Channel, payload); err != nil {
			svc.logger.Error("parse capabilities failure", zap.Error(err))
			return nil
		}
	case HeartbeatsChannel:
		if err := svc.handleHeartbeat(msg.Publisher, msg.Channel, payload); err != nil {
			svc.logger.Error("parse heartbeat failure", zap.Error(err))
		}
	case RPCToCoreChannel:
		svc.logger.Error("implement me: RPCToCoreChannel")
	case RPCFromCoreChannel:
		svc.logger.Error("implement me: RPCFromCoreChannel")
	case LogChannel:
		svc.logger.Error("implement me: LogChannel")
	default:
		svc.logger.Warn("unsupported/unhandled agent subtopic, ignoring",
			zap.String("subtopic", msg.Subtopic),
			zap.String("thing_id", msg.Publisher),
			zap.String("channel_id", msg.Channel))
	}

	return nil
}

func (svc fleetCommsService) Start() error {
	if err := svc.agentPubSub.Subscribe(fmt.Sprintf("channels.*.%s", CapabilitiesChannel), svc.handleMsgFromAgent); err != nil {
		return err
	}
	if err := svc.agentPubSub.Subscribe(fmt.Sprintf("channels.*.%s", HeartbeatsChannel), svc.handleMsgFromAgent); err != nil {
		return err
	}
	if err := svc.agentPubSub.Subscribe(fmt.Sprintf("channels.*.%s", RPCToCoreChannel), svc.handleMsgFromAgent); err != nil {
		return err
	}
	if err := svc.agentPubSub.Subscribe(fmt.Sprintf("channels.*.%s", LogChannel), svc.handleMsgFromAgent); err != nil {
		return err
	}
	svc.logger.Info("subscribed to agent channels")
	return nil
}

func (svc fleetCommsService) Stop() error {
	if err := svc.agentPubSub.Unsubscribe(fmt.Sprintf("channels.*.%s", CapabilitiesChannel)); err != nil {
		return err
	}
	if err := svc.agentPubSub.Unsubscribe(fmt.Sprintf("channels.*.%s", HeartbeatsChannel)); err != nil {
		return err
	}
	if err := svc.agentPubSub.Unsubscribe(fmt.Sprintf("channels.*.%s", RPCToCoreChannel)); err != nil {
		return err
	}
	if err := svc.agentPubSub.Unsubscribe(fmt.Sprintf("channels.*.%s", LogChannel)); err != nil {
		return err
	}
	svc.logger.Info("unsubscribed from agent channels")
	return nil
}
