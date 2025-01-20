package otlpmqttexporter

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/orb-community/orb/agent/otel"
	"go.uber.org/zap"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr         = "otlpmqtt"
	defaultMQTTAddr = "localhost:1883"
	defaultMQTTId   = "uuid1"
	defaultMQTTKey  = "uuid2"
	defaultName     = "pktvisor"
	// For testing will disable  TLS
	defaultTLS = false
)

// NewFactory creates a factory for OTLP exporter.
// Reducing the scope to just Metrics since it is our use-case
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		CreateDefaultConfig,
		exporter.WithMetrics(CreateMetricsExporter, component.StabilityLevelStable))
}

func CreateConfig(addr, id, key, channel, pktvisor, topic string, bridgeService otel.AgentBridgeService) component.Config {
	return &Config{
		TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
		QueueSettings:   exporterhelper.NewDefaultQueueSettings(),
		RetrySettings:   exporterhelper.NewDefaultRetrySettings(),
		Topic:           topic,
		Address:         addr,
		Id:              id,
		Key:             key,
		ChannelID:       channel,
		PktVisorVersion: pktvisor,
		OrbAgentService: bridgeService,
	}
}

func CreateDefaultSettings(logger *zap.Logger) exporter.CreateSettings {
	return exporter.CreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: noop.NewTracerProvider(),
			MeterProvider:  metric.NewMeterProvider(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
}

func CreateDefaultConfig() component.Config {
	base := fmt.Sprintf("channels/%s/messages", defaultMQTTId)
	topic := fmt.Sprintf("%s/otlp/%s", base, defaultName)
	return &Config{
		TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
		QueueSettings:   exporterhelper.NewDefaultQueueSettings(),
		RetrySettings:   exporterhelper.NewDefaultRetrySettings(),
		Address:         defaultMQTTAddr,
		Id:              defaultMQTTId,
		Key:             defaultMQTTKey,
		ChannelID:       base,
		TLS:             defaultTLS,
		Topic:           topic,
	}
}

func CreateConfigClient(client *mqtt.Client, topic, pktvisor string, bridgeService otel.AgentBridgeService) component.Config {
	return &Config{
		TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
		QueueSettings:   exporterhelper.NewDefaultQueueSettings(),
		RetrySettings:   exporterhelper.NewDefaultRetrySettings(),
		Client:          client,
		Topic:           topic,
		PktVisorVersion: pktvisor,
		OrbAgentService: bridgeService,
	}
}

func CreateTracesExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Traces, error) {
	oce, err := newExporter(cfg, set, ctx)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)

	return exporterhelper.NewTracesExporter(
		ctx,
		set,
		cfg,
		oce.pushTraces,
		exporterhelper.WithStart(oce.start),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// explicitly disable since we rely on http.Client timeout logic.
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings))
}

func CreateMetricsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Metrics, error) {
	oce, err := newExporter(cfg, set, ctx)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)
	return exporterhelper.NewMetricsExporter(
		ctx,
		set,
		cfg,
		oce.pushMetrics,
		exporterhelper.WithStart(oce.start),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// explicitly disable since we rely on http.Client timeout logic.
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings),
		exporterhelper.WithShutdown(oce.shutdown))
}

func CreateLogsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Logs, error) {
	oce, err := newExporter(cfg, set, ctx)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)

	return exporterhelper.NewLogsExporter(
		ctx,
		set,
		cfg,
		oce.pushLogs,
		exporterhelper.WithStart(oce.start),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// explicitly disable since we rely on http.Client timeout logic.
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings))
}
