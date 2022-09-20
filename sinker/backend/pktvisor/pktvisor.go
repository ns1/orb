/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/ns1labs/orb/fleet"
	"github.com/ns1labs/orb/fleet/pb"
	"github.com/ns1labs/orb/pkg/errors"
	"github.com/ns1labs/orb/sinker/backend"
	"github.com/ns1labs/orb/sinker/prometheus"
	"go.uber.org/zap"
)

var _ backend.Backend = (*pktvisorBackend)(nil)

type pktvisorBackend struct {
	logger *zap.Logger
}

type context struct {
	agent        *pb.AgentInfoRes
	agentID      string
	policyID     string
	policyName   string
	deviceID     string
	handlerLabel string
	tags         map[string]string
	logger       *zap.Logger
}

func (p pktvisorBackend) ProcessMetrics(agent *pb.AgentInfoRes, agentID string, data fleet.AgentMetricsRPCPayload) ([]prometheus.TimeSeries, error) {
	// TODO check pktvisor version in data.BEVersion against PktvisorVersion
	if data.Format != "json" {
		p.logger.Warn("ignoring non-json pktvisor payload", zap.String("format", data.Format))
		return nil, nil
	}
	// unmarshal pktvisor metrics
	var metrics map[string]map[string]interface{}
	err := json.Unmarshal(data.Data, &metrics)
	if err != nil {
		p.logger.Warn("unable to unmarshal pktvisor metric payload", zap.Any("payload", data.Data))
		return nil, err
	}

	tags := make(map[string]string)
	for k, v := range agent.AgentTags {
		tags[k] = v
	}
	for k, v := range agent.OrbTags {
		tags[k] = v
	}

	context := context{
		agent:        agent,
		agentID:      agentID,
		policyID:     data.PolicyID,
		policyName:   data.PolicyName,
		deviceID:     "",
		handlerLabel: "",
		tags:         tags,
		logger:       p.logger,
	}
	stats := make(map[string]StatSnapshot)
	for handlerLabel, handlerData := range metrics {
		if data, ok := handlerData["pcap"]; ok {
			sTmp := StatSnapshot{}
			err := mapstructure.Decode(data, &sTmp.Pcap)
			if err != nil {
				p.logger.Error("error decoding pcap handler", zap.Error(err))
				continue
			}
			stats[handlerLabel] = sTmp
		} else if data, ok := handlerData["dns"]; ok {
			sTmp := StatSnapshot{}
			err := mapstructure.Decode(data, &sTmp.DNS)
			if err != nil {
				p.logger.Error("error decoding dns handler", zap.Error(err))
				continue
			}
			stats[handlerLabel] = sTmp
		} else if data, ok := handlerData["packets"]; ok {
			sTmp := StatSnapshot{}
			err := mapstructure.Decode(data, &sTmp.Packets)
			if err != nil {
				p.logger.Error("error decoding packets handler", zap.Error(err))
				continue
			}
			stats[handlerLabel] = sTmp
		} else if data, ok := handlerData["dhcp"]; ok {
			sTmp := StatSnapshot{}
			err := mapstructure.Decode(data, &sTmp.DHCP)
			if err != nil {
				p.logger.Error("error decoding dhcp handler", zap.Error(err))
				continue
			}
			stats[handlerLabel] = sTmp
		} else if data, ok := handlerData["flow"]; ok {
			sTmp := StatSnapshot{}
			err := mapstructure.Decode(data, &sTmp.Flow)
			if err != nil {
				p.logger.Error("error decoding dhcp handler", zap.Error(err))
				continue
			}
			stats[handlerLabel] = sTmp
		}
	}
	return parseToProm(&context, stats), nil
}

func parseToProm(ctxt *context, statsMap map[string]StatSnapshot) prometheus.TSList {
	var finalTs = prometheus.TSList{}
	for handlerLabel, stats := range statsMap {
		var tsList = prometheus.TSList{}
		statsMap := structs.Map(stats)
		ctxt.handlerLabel = handlerLabel
		if stats.Flow != nil {
			convertFlowToPromParticle(ctxt, statsMap, "", &tsList)
		} else {
			convertToPromParticle(ctxt, statsMap, "", &tsList)
		}
		finalTs = append(finalTs, tsList...)
	}
	return finalTs
}

func convertToPromParticle(ctxt *context, statsMap map[string]interface{}, label string, tsList *prometheus.TSList) {
	for key, value := range statsMap {
		switch statistic := value.(type) {
		case map[string]interface{}:
			// Call convertToPromParticle recursively until the last interface of the StatSnapshot struct
			// The prom particle label it's been formed during the recursive call (concatenation)
			convertToPromParticle(ctxt, statistic, label+key, tsList)
		// The StatSnapshot has two ways to record metrics (i.e. Live int64 `mapstructure:"live"`)
		// It's why we check if the type is int64
		case int64:
			{
				// Use this regex to identify if the value it's a quantile
				var matchFirstQuantile = regexp.MustCompile("^([Pp])+[0-9]")
				if ok := matchFirstQuantile.MatchString(key); ok {
					// If it's quantile, needs to be parsed to prom quantile format
					tsList = makePromParticle(ctxt, label, key, value, tsList, ok, "")
				} else {
					tsList = makePromParticle(ctxt, label+key, "", value, tsList, false, "")
				}
			}
		// The StatSnapshot has two ways to record metrics (i.e. P50 float64 `mapstructure:"p50"`)
		// It's why we check if the type is float64
		case float64:
			{
				// Use this regex to identify if the value it's a quantile
				var matchFirstQuantile = regexp.MustCompile("^[Pp]+[0-9]")
				if ok := matchFirstQuantile.MatchString(key); ok {
					// If it's quantile, needs to be parsed to prom quantile format
					tsList = makePromParticle(ctxt, label, key, value, tsList, ok, "")
				} else {
					tsList = makePromParticle(ctxt, label+key, "", value, tsList, false, "")
				}
			}
		// The StatSnapshot has two ways to record metrics (i.e. TopIpv4   []NameCount   `mapstructure:"top_ipv4"`)
		// It's why we check if the type is []interface
		// Here we extract the value for Name and Estimate
		case []interface{}:
			{
				for _, value := range statistic {
					m, ok := value.(map[string]interface{})
					if !ok {
						return
					}
					var promLabel string
					var promDataPoint interface{}
					for k, v := range m {
						switch k {
						case "Name":
							{
								promLabel = fmt.Sprintf("%v", v)
							}
						case "Estimate":
							{
								promDataPoint = v
							}
						}
					}
					tsList = makePromParticle(ctxt, label+key, promLabel, promDataPoint, tsList, false, key)
				}
			}
		}
	}
}

func convertFlowToPromParticle(ctxt *context, statsMap map[string]interface{}, label string, tsList *prometheus.TSList) {
	for key, value := range statsMap {
		switch statistic := value.(type) {
		case map[string]interface{}:
			// Call convertToPromParticle recursively until the last interface of the StatSnapshot struct
			// The prom particle label it's been formed during the recursive call (concatenation)
			ipv6_regex := `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
			ipv4_regex := `^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`

			if ok, _ := regexp.MatchString(ipv4_regex+`|`+ipv6_regex, key); ok {
				if ok = strings.Contains(label, "Devices"); !ok {
					return
				}
				label = strings.ReplaceAll(label, "Devices", "")
				ctxt.deviceID = key
				convertFlowToPromParticle(ctxt, statistic, label, tsList)
			} else {
				convertFlowToPromParticle(ctxt, statistic, label+key, tsList)
			}

		// The StatSnapshot has two ways to record metrics (i.e. Live int64 `mapstructure:"live"`)
		// It's why we check if the type is int64
		case int64:
			{
				// Use this regex to identify if the value it's a quantile
				var matchFirstQuantile = regexp.MustCompile("^([Pp])+[0-9]")
				if ok := matchFirstQuantile.MatchString(key); ok {
					// If it's quantile, needs to be parsed to prom quantile format
					tsList = makePromParticle(ctxt, label, key, value, tsList, ok, "")
				} else {
					tsList = makePromParticle(ctxt, label+key, "", value, tsList, false, "")
				}
			}
		// The StatSnapshot has two ways to record metrics (i.e. TopIpv4   []NameCount   `mapstructure:"top_ipv4"`)
		// It's why we check if the type is []interface
		// Here we extract the value for Name and Estimate
		case []interface{}:
			{
				for _, value := range statistic {
					m, ok := value.(map[string]interface{})
					if !ok {
						return
					}
					var promLabel string
					var promDataPoint interface{}
					for k, v := range m {
						switch k {
						case "Name":
							{
								promLabel = fmt.Sprintf("%v", v)
							}
						case "Estimate":
							{
								promDataPoint = v
							}
						}
					}
					tsList = makePromParticle(ctxt, label+key, promLabel, promDataPoint, tsList, false, key)
				}
			}
		}
	}
}

func makePromParticle(ctxt *context, label string, k string, v interface{}, tsList *prometheus.TSList, quantile bool, name string) *prometheus.TSList {
	mapQuantiles := make(map[string]string)
	mapQuantiles["P50"] = "0.5"
	mapQuantiles["P90"] = "0.9"
	mapQuantiles["P95"] = "0.95"
	mapQuantiles["P99"] = "0.99"

	var dpFlag dp
	var labelsListFlag labelList
	if err := labelsListFlag.Set(fmt.Sprintf("__name__;%s", camelToSnake(label))); err != nil {
		handleParticleError(ctxt, err)
		return tsList
	}
	if err := labelsListFlag.Set("instance;" + ctxt.agent.AgentName); err != nil {
		handleParticleError(ctxt, err)
		return tsList
	}
	if err := labelsListFlag.Set("agent_id;" + ctxt.agentID); err != nil {
		handleParticleError(ctxt, err)
		return tsList
	}
	if err := labelsListFlag.Set("agent;" + ctxt.agent.AgentName); err != nil {
		handleParticleError(ctxt, err)
		return tsList
	}
	if err := labelsListFlag.Set("policy_id;" + ctxt.policyID); err != nil {
		handleParticleError(ctxt, err)
		return tsList
	}
	if err := labelsListFlag.Set("policy;" + ctxt.policyName); err != nil {
		handleParticleError(ctxt, err)
		return tsList
	}
	if err := labelsListFlag.Set("handler;" + ctxt.handlerLabel); err != nil {
		handleParticleError(ctxt, err)
		return tsList
	}
	if ctxt.deviceID != "" {
		if err := labelsListFlag.Set("device;" + ctxt.deviceID); err != nil {
			handleParticleError(ctxt, err)
			ctxt.deviceID = ""
			return tsList
		}
	}

	for k, v := range ctxt.tags {
		if err := labelsListFlag.Set(k + ";" + v); err != nil {
			handleParticleError(ctxt, err)
			return tsList
		}
	}

	if k != "" {
		if quantile {
			if value, ok := mapQuantiles[k]; ok {
				if err := labelsListFlag.Set(fmt.Sprintf("quantile;%s", value)); err != nil {
					handleParticleError(ctxt, err)
					return tsList
				}
			}
		} else {
			parsedName, err := topNMetricsParser(name)
			if err != nil {
				ctxt.logger.Error("failed to parse Top N metric, default value it'll be used", zap.Error(err))
				parsedName = "name"
			}
			if err := labelsListFlag.Set(fmt.Sprintf("%s;%s", parsedName, k)); err != nil {
				handleParticleError(ctxt, err)
				return tsList
			}
		}
	}
	if err := dpFlag.Set(fmt.Sprintf("now,%d", v)); err != nil {
		if err := dpFlag.Set(fmt.Sprintf("now,%v", v)); err != nil {
			handleParticleError(ctxt, err)
			return tsList
		}
	}
	timeSeries := prometheus.TimeSeries{
		Labels:    labelsListFlag,
		Datapoint: prometheus.Datapoint(dpFlag),
	}
	*tsList = append(*tsList, timeSeries)
	return tsList
}

func handleParticleError(ctxt *context, err error) {
	ctxt.logger.Error("failed to set prometheus element", zap.Error(err))
}

func camelToSnake(s string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	// Approach to avoid change the values to TopGeoLoc and TopASN
	// Should continue camel case or upper case
	var matchExcept = regexp.MustCompile(`(oLoc$|pASN$)`)
	sub := matchExcept.Split(s, 2)
	var strExcept = ""
	if len(sub) > 1 {
		strExcept = matchExcept.FindAllString(s, 1)[0]
		if strExcept == "pASN" {
			strExcept = "p_ASN"
		}
		s = sub[0]
	}

	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	lower := strings.ToLower(snake)
	return lower + strExcept
}

func topNMetricsParser(label string) (string, error) {
	mapNMetrics := make(map[string]string)
	mapNMetrics["TopGeoLocECS"] = "geo_loc"
	mapNMetrics["TopGeoLoc"] = "geo_loc"
	mapNMetrics["TopAsnECS"] = "asn"
	mapNMetrics["TopASN"] = "asn"
	mapNMetrics["TopQueryECS"] = "ecs"
	mapNMetrics["TopIpv6"] = "ipv6"
	mapNMetrics["TopIpv4"] = "ipv4"
	mapNMetrics["TopQname2"] = "qname"
	mapNMetrics["TopQname3"] = "qname"
	mapNMetrics["TopQnameByRespBytes"] = "qname"
	mapNMetrics["TopNxdomain"] = "qname"
	mapNMetrics["TopQtype"] = "qtype"
	mapNMetrics["TopRcode"] = "rcode"
	mapNMetrics["TopREFUSED"] = "qname"
	mapNMetrics["TopNODATA"] = "qname"
	mapNMetrics["TopSRVFAIL"] = "qname"
	mapNMetrics["TopUDPPorts"] = "port"
	mapNMetrics["TopSlow"] = "qname"
	mapNMetrics["TopGeoLocBytes"] = "geo_loc"
	mapNMetrics["TopGeoLocPackes"] = "geo_loc"
	mapNMetrics["TopAsnBytes"] = "asn"
	mapNMetrics["TopAsnPackets"] = "asn"
	mapNMetrics["TopDstIpsBytes"] = "ip"
	mapNMetrics["TopDstIpsPackets"] = "ip"
	mapNMetrics["TopSrcIpsBytes"] = "ip"
	mapNMetrics["TopSrcIpsPackets"] = "ip"
	mapNMetrics["TopDstPortsBytes"] = "port"
	mapNMetrics["TopDstPortsPackets"] = "port"
	mapNMetrics["TopSrcPortsBytes"] = "port"
	mapNMetrics["TopSrcPortsPackets"] = "port"
	mapNMetrics["TopDstIpsAndPortBytes"] = "ip_port"
	mapNMetrics["TopDstIpsAndPortPackets"] = "ip_port"
	mapNMetrics["TopSrcIpsAndPortBytes"] = "ip_port"
	mapNMetrics["TopSrcIpsAndPortPackets"] = "ip_port"
	mapNMetrics["TopConversationsBytes"] = "conversations"
	mapNMetrics["TopConversationsPackets"] = "conversations"
	mapNMetrics["TopInIfIndexBytes"] = "index"
	mapNMetrics["TopInIfIndexPackets"] = "index"
	mapNMetrics["TopOutIfIndexBytes"] = "index"
	mapNMetrics["TopOutIfIndexPackets"] = "index"
	if value, ok := mapNMetrics[label]; ok {
		return value, nil
	} else {
		return "", errors.New(fmt.Sprintf("top N metric not mapped for parse:  %s", label))
	}
}

func Register(logger *zap.Logger) bool {
	backend.Register("pktvisor", &pktvisorBackend{logger: logger})
	return true
}
