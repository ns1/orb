/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package pktvisor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ns1labs/orb/agent/otel/otlpmqttexporter"
	"github.com/ns1labs/orb/agent/otel/pktvisorreceiver"
	"github.com/ns1labs/orb/fleet"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (p *pktvisorBackend) scrapeMetrics(period uint) (map[string]interface{}, error) {
	var metrics map[string]interface{}
	err := p.request(fmt.Sprintf("policies/__all/metrics/bucket/%d", period), &metrics, http.MethodGet, http.NoBody, "application/json", ScrapeTimeout)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (p *pktvisorBackend) createOtlpMqttExporter(ctx context.Context) (component.MetricsExporter, error) {

	if p.mqttClient != nil {
		cfg := otlpmqttexporter.CreateConfigClient(p.mqttClient, p.otlpMetricsTopic, p.pktvisorVersion)
		set := otlpmqttexporter.CreateDefaultSettings(p.logger)
		// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
		exporter, err := otlpmqttexporter.CreateMetricsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	} else {
		cfg := otlpmqttexporter.CreateConfig(p.mqttConfig.Address, p.mqttConfig.Id, p.mqttConfig.Key,
			p.mqttConfig.ChannelID, p.pktvisorVersion, p.otlpMetricsTopic)
		set := otlpmqttexporter.CreateDefaultSettings(p.logger)
		// Create the OTLP metrics exporter that'll receive and verify the metrics produced.
		exporter, err := otlpmqttexporter.CreateMetricsExporter(ctx, set, cfg)
		if err != nil {
			return nil, err
		}
		return exporter, nil
	}

}

func createReceiver(ctx context.Context, exporter component.MetricsExporter, logger *zap.Logger) (component.MetricsReceiver, error) {
	set := pktvisorreceiver.CreateDefaultSettings(logger)
	cfg := pktvisorreceiver.CreateDefaultConfig()
	// Create the Prometheus receiver and pass in the previously created Prometheus exporter.
	receiver, err := pktvisorreceiver.CreateMetricsReceiver(ctx, set, cfg, exporter)
	if err != nil {
		return nil, err
	}
	return receiver, nil
}

func (p *pktvisorBackend) scrapeDefault() error {
	// scrape all policy json output with one call every minute.
	// TODO support policies with custom bucket times
	job, err := p.scraper.Every(1).Minute().WaitForSchedule().Do(func() {
		metrics, err := p.scrapeMetrics(1)
		if err != nil {
			p.logger.Error("scrape failed", zap.Error(err))
			return
		}
		if len(metrics) == 0 {
			p.logger.Warn("scrape: no policies found, skipping")
			return
		}

		var batchPayload []fleet.AgentMetricsRPCPayload
		totalSize := 0
		for pName, pMetrics := range metrics {
			policyData, err := p.policyRepo.GetByName(pName)
			if err != nil {
				p.logger.Warn("skipping pktvisor policy not managed by orb", zap.String("policy", pName), zap.Error(err))
				continue
			}
			payloadData, err := json.Marshal(pMetrics)
			if err != nil {
				p.logger.Error("error marshalling scraped metric json", zap.String("policy", pName), zap.Error(err))
				continue
			}
			metricPayload := fleet.AgentMetricsRPCPayload{
				PolicyID:   policyData.ID,
				PolicyName: policyData.Name,
				Datasets:   policyData.GetDatasetIDs(),
				Format:     "json",
				BEVersion:  p.pktvisorVersion,
				Data:       payloadData,
			}
			batchPayload = append(batchPayload, metricPayload)
			totalSize += len(payloadData)
			policyData.LastScrapeBytes = int64(totalSize)
			policyData.LastScrapeTS = time.Now()
			err = p.policyRepo.Update(policyData)
			if err != nil {
				p.logger.Error("unable to update policy repo during scrape", zap.Error(err))
			}
			p.logger.Info("scraped metrics for policy", zap.String("policy", pName), zap.String("policy_id", policyData.ID), zap.Int("payload_size_b", len(payloadData)))
		}

		rpc := fleet.AgentMetricsRPC{
			SchemaVersion: fleet.CurrentRPCSchemaVersion,
			Func:          fleet.AgentMetricsRPCFunc,
			Payload:       batchPayload,
		}

		body, err := json.Marshal(rpc)
		if err != nil {
			p.logger.Error("error marshalling metric rpc payload", zap.Error(err))
			return
		}

		if token := p.mqttClient.Publish(p.metricsTopic, 1, false, body); token.Wait() && token.Error() != nil {
			p.logger.Error("error sending metrics RPC", zap.String("topic", p.metricsTopic), zap.Error(token.Error()))
			return
		}

		p.logger.Info("scraped and published metrics", zap.String("topic", p.metricsTopic), zap.Int("payload_size_b", totalSize), zap.Int("batch_count", len(batchPayload)))

	})

	if err != nil {
		return err
	}

	job.SingletonMode()
	return nil
}

func (p *pktvisorBackend) scrapeOpenTelemetry(ctx context.Context) {
	go func() {
		startExpCtx, cancelFunc := context.WithCancel(ctx)
		var ok bool
		var err error
		for i := 1; i < 10; i++ {
			select {
			case <-startExpCtx.Done():
				return
			default:
				if p.mqttClient != nil {
					var errStartExp error
					p.exporter, errStartExp = p.createOtlpMqttExporter(ctx)
					if errStartExp != nil {
						p.logger.Error("failed to create a exporter", zap.Error(err))
						return
					}

					p.receiver, err = createReceiver(ctx, p.exporter, p.logger)
					if err != nil {
						p.logger.Error("failed to create a receiver", zap.Error(err))
						return
					}

					err = p.exporter.Start(ctx, nil)
					if err != nil {
						p.logger.Error("otel mqtt exporter startup error", zap.Error(err))
						return
					}

					err = p.receiver.Start(ctx, nil)
					if err != nil {
						p.logger.Error("otel receiver startup error", zap.Error(err))
						return
					}

					ok = true
					return
				} else {
					p.logger.Info("waiting until mqtt client is connected", zap.String("wait time", (time.Duration(i)*time.Second).String()))
					time.Sleep(time.Duration(i) * time.Second)
					continue
				}
			}
		}
		if !ok {
			p.logger.Error("mqtt did not established a connection, stopping agent")
			p.Stop(startExpCtx)
		}
		cancelFunc()
		return
	}()
}
