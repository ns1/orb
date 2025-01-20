/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package agent

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/orb-community/orb/agent/backend"
	"github.com/orb-community/orb/agent/config"
	"github.com/orb-community/orb/fleet"
	"go.uber.org/zap"
)

func (a *orbAgent) connect(ctx context.Context, config config.MQTTConfig) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().AddBroker(config.Address).SetClientID(config.Id)
	opts.SetUsername(config.Id)
	opts.SetPassword(config.Key)
	opts.SetKeepAlive(10 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		a.logger.Info("message on unknown channel, ignoring", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		a.logger.Error("connection to mqtt lost", zap.Error(err))
		// If it is a bug on the mqttclient, stop the agent
		if strings.Contains(err.Error(), "BUG") {
			a.Stop(ctx)
			return
		}
		a.logger.Info("reconnecting....")
		a.requestReconnection(ctx, a.client, config)
	})
	opts.SetPingTimeout(5 * time.Second)
	opts.SetAutoReconnect(false)
	opts.SetCleanSession(true)
	opts.SetConnectTimeout(5 * time.Minute)
	opts.SetResumeSubs(true)
	opts.SetReconnectingHandler(func(client mqtt.Client, options *mqtt.ClientOptions) {
		go func() {
			ok := false
			for i := 1; i < 10; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					if len(a.backends) == 0 {
						time.Sleep(time.Duration(i) * time.Second)
						continue
					}
					for name, be := range a.backends {
						backendStatus, s, _ := be.GetRunningStatus()
						a.logger.Debug("backend in status", zap.String("backend", name), zap.String("status", s))
						switch backendStatus {
						case backend.Running:
							ok = true
							a.requestReconnection(ctx, client, config)
							return
						case backend.Waiting:
							ok = true
							a.requestReconnection(ctx, client, config)
							return
						default:
							a.logger.Info("waiting until a backend is in running state", zap.String("backend", name),
								zap.String("current state", s), zap.String("wait time", (time.Duration(i)*time.Second).String()))
							time.Sleep(time.Duration(i) * time.Second)
							continue
						}
					}
				}
			}
			if !ok {
				a.logger.Error("backend wasn't able to change to running, stopping connection")
				ctx.Done()
			}
		}()
	})
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		go func() {
			ok := false
			for i := 1; i < 10; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					if len(a.backends) == 0 {
						time.Sleep(time.Duration(i) * time.Second)
						continue
					}
					for name, be := range a.backends {
						backendStatus, s, _ := be.GetRunningStatus()
						a.logger.Debug("backend in status", zap.String("backend", name), zap.String("status", s))
						switch backendStatus {
						case backend.Running:
							ok = true
							a.requestReconnection(ctx, client, config)
							return
						case backend.Waiting:
							ok = true
							a.requestReconnection(ctx, client, config)
							return
						default:
							a.logger.Info("waiting until a backend is in running state", zap.String("backend", name),
								zap.String("current state", s), zap.String("wait time", (time.Duration(i)*time.Second).String()))
							time.Sleep(time.Duration(i) * time.Second)
							continue
						}
					}
				}
			}
			if !ok {
				a.logger.Error("backend wasn't able to change to running, stopping connection")
				ctx.Done()
			}
		}()
	})

	if !a.config.OrbAgent.TLS.Verify {
		opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return c, nil
}

func (a *orbAgent) requestReconnection(ctx context.Context, client mqtt.Client, config config.MQTTConfig) {
	a.nameAgentRPCTopics(config.ChannelID)
	for name, be := range a.backends {
		be.SetCommsClient(config.Id, &client, fmt.Sprintf("%s/?/%s", a.baseTopic, name))
	}
	a.agent_id = config.Id

	if token := client.Subscribe(a.rpcFromCoreTopic, 1, a.handleRPCFromCore); token.Wait() && token.Error() != nil {
		a.logger.Error("failed to subscribe to agent control plane RPC topic", zap.String("topic", a.rpcFromCoreTopic), zap.Error(token.Error()))
		a.logger.Error("critical failure: unable to subscribe to control plane")
		a.Stop(ctx)
		return
	}

	err := a.sendCapabilities()
	if err != nil {
		a.logger.Error("failed to send agent capabilities", zap.Error(err))
	}

	err = a.sendGroupMembershipReq()
	if err != nil {
		a.logger.Error("failed to send group membership request", zap.Error(err))
	}
}

func (a *orbAgent) nameAgentRPCTopics(channelId string) {

	base := fmt.Sprintf("channels/%s/messages", channelId)
	a.rpcToCoreTopic = fmt.Sprintf("%s/%s", base, fleet.RPCToCoreTopic)
	a.rpcFromCoreTopic = fmt.Sprintf("%s/%s", base, fleet.RPCFromCoreTopic)
	a.capabilitiesTopic = fmt.Sprintf("%s/%s", base, fleet.CapabilitiesTopic)
	a.heartbeatsTopic = fmt.Sprintf("%s/%s", base, fleet.HeartbeatsTopic)
	a.logTopic = fmt.Sprintf("%s/%s", base, fleet.LogTopic)
	a.baseTopic = base

}

func (a *orbAgent) unsubscribeGroupChannels() {
	a.logger.Debug("calling to unsub group channels")
	for id, groupInfo := range a.groupsInfos {
		base := fmt.Sprintf("channels/%s/messages", groupInfo.ChannelID)
		rpcFromCoreTopic := fmt.Sprintf("%s/%s", base, fleet.RPCFromCoreTopic)
		if token := a.client.Unsubscribe(rpcFromCoreTopic); token.Wait() && token.Error() != nil {
			a.logger.Warn("failed to unsubscribe to group channel", zap.String("group_id", id), zap.String("group_name", groupInfo.Name), zap.String("topic", groupInfo.ChannelID), zap.Error(token.Error()))
		}
		a.logger.Info("completed RPC unsubscription to group", zap.String("group_id", id), zap.String("group_name", groupInfo.Name), zap.String("topic", rpcFromCoreTopic))
	}
	a.groupsInfos = make(map[string]GroupInfo)
}

func (a *orbAgent) unsubscribeGroupChannel(channelID string, agentGroupID string) {
	base := fmt.Sprintf("channels/%s/messages", channelID)
	rpcFromCoreTopic := fmt.Sprintf("%s/%s", base, fleet.RPCFromCoreTopic)
	if token := a.client.Unsubscribe(channelID); token.Wait() && token.Error() != nil {
		a.logger.Warn("failed to unsubscribe to group channel", zap.String("topic", rpcFromCoreTopic), zap.Error(token.Error()))
		return
	}
	a.logger.Info("completed RPC unsubscription to group", zap.String("topic", rpcFromCoreTopic))
	delete(a.groupsInfos, agentGroupID)
}

func (a *orbAgent) removeDatasetFromPolicy(datasetID string, policyID string) {
	for _, be := range a.backends {
		a.policyManager.RemovePolicyDataset(policyID, datasetID, be)
	}
}

func (a *orbAgent) startComms(ctx context.Context, config config.MQTTConfig) error {

	var err error
	a.logger.Debug("starting mqtt connection")
	if a.client == nil || !a.client.IsConnected() {
		a.client, err = a.connect(ctx, config)
		if err != nil {
			a.logger.Error("connection failed", zap.String("channel", config.ChannelID), zap.String("agent_id", config.Id), zap.Error(err))
			return ErrMqttConnection
		}
		// Store the data from connection to cloud config within agent.
		a.config.OrbAgent.Cloud.MQTT.Id = config.Id
		a.config.OrbAgent.Cloud.MQTT.Key = config.Key
		a.config.OrbAgent.Cloud.MQTT.Address = config.Address
		a.config.OrbAgent.Cloud.MQTT.ChannelID = config.ChannelID
	} else {
		a.requestReconnection(ctx, a.client, config)
	}

	return nil
}

func (a *orbAgent) subscribeGroupChannels(groups []fleet.GroupMembershipData) {
	for _, groupData := range groups {

		base := fmt.Sprintf("channels/%s/messages", groupData.ChannelID)
		rpcFromCoreTopic := fmt.Sprintf("%s/%s", base, fleet.RPCFromCoreTopic)

		token := a.client.Subscribe(rpcFromCoreTopic, 1, a.handleGroupRPCFromCore)
		if token.Error() != nil {
			a.logger.Error("failed to subscribe to group channel/topic", zap.String("group_id", groupData.GroupID), zap.String("group_name", groupData.Name), zap.String("topic", rpcFromCoreTopic), zap.Error(token.Error()))
			continue
		}
		ok := token.WaitTimeout(time.Second * 5)
		if ok && token.Error() != nil {
			a.logger.Error("failed to subscribe to group channel/topic", zap.String("group_id", groupData.GroupID), zap.String("group_name", groupData.Name), zap.String("topic", rpcFromCoreTopic), zap.Error(token.Error()))
			continue
		}
		if !ok {
			a.logger.Error("failed to subscribe to group channel/topic: time out", zap.String("group_id", groupData.GroupID), zap.String("group_name", groupData.Name), zap.String("topic", rpcFromCoreTopic))
			continue
		}
		a.logger.Info("completed RPC subscription to group", zap.String("group_id", groupData.GroupID), zap.String("group_name", groupData.Name), zap.String("topic", rpcFromCoreTopic))
		a.groupsInfos[groupData.GroupID] = GroupInfo{
			Name:      groupData.Name,
			ChannelID: groupData.ChannelID,
		}
	}
}
