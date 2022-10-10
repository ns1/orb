package otlpmqttexporter

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
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
func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		CreateDefaultConfig,
		component.WithMetricsExporter(CreateMetricsExporter))
}

func CreateConfig(addr, id, key, channel, pktvisor string) config.Exporter {
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		QueueSettings:    exporterhelper.NewDefaultQueueSettings(),
		RetrySettings:    exporterhelper.NewDefaultRetrySettings(),
		Address:          addr,
		Id:               id,
		Key:              key,
		ChannelID:        channel,
		PktVisorVersion:  pktvisor,
	}
}

func CreateDefaultConfig() config.Exporter {
	base := fmt.Sprintf("channels/%s/messages", defaultMQTTId)
	metricsTopic := fmt.Sprintf("%s/be/%s", base, defaultName)
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		QueueSettings:    exporterhelper.NewDefaultQueueSettings(),
		RetrySettings:    exporterhelper.NewDefaultRetrySettings(),
		Address:          defaultMQTTAddr,
		Id:               defaultMQTTId,
		Key:              defaultMQTTKey,
		ChannelID:        base,
		TLS:              defaultTLS,
		MetricsTopic:     metricsTopic,
	}
}

func CreateConfigClient(client mqtt.Client, metricsTopic, pktvisor string) config.Exporter {
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		QueueSettings:    exporterhelper.NewDefaultQueueSettings(),
		RetrySettings:    exporterhelper.NewDefaultRetrySettings(),
		Client:           client,
		MetricsTopic:     metricsTopic,
		PktVisorVersion:  pktvisor,
	}
}

func CreateDefaultSettings(logger *zap.Logger) component.ExporterCreateSettings {
	return component.ExporterCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  global.MeterProvider(),
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}
}

func createTracesExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.TracesExporter, error) {
	oce, err := newExporter(cfg, set)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)

	return exporterhelper.NewTracesExporter(
		cfg,
		set,
		oce.pushTraces,
		exporterhelper.WithStart(oce.start),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// explicitly disable since we rely on http.Client timeout logic.
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings))
}

func CreateMetricsExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.MetricsExporter, error) {
	oce, err := newExporter(cfg, set)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)

	return exporterhelper.NewMetricsExporter(
		cfg,
		set,
		oce.pushMetrics,
		exporterhelper.WithStart(oce.start),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// explicitly disable since we rely on http.Client timeout logic.
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings))
}

func createLogsExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.LogsExporter, error) {
	oce, err := newExporter(cfg, set)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)

	return exporterhelper.NewLogsExporter(
		cfg,
		set,
		oce.pushLogs,
		exporterhelper.WithStart(oce.start),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		// explicitly disable since we rely on http.Client timeout logic.
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings))
}
