package otlpmqttexporter

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.opentelemetry.io/collector/consumer/consumererror"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	"go.uber.org/zap"
)

type baseExporter struct {
	// Input configuration.
	config   *Config
	client   *http.Client
	logger   *zap.Logger
	settings component.TelemetrySettings
	// Default user-agent header.
	userAgent string
}

func (e *baseExporter) compressBrotli(data []byte) []byte {
	var b bytes.Buffer
	w := brotli.NewWriterLevel(&b, brotli.BestCompression)
	_, err := w.Write(data)
	if err != nil {
		return nil
	}
	err = w.Close()
	if err != nil {
		return nil
	}
	return b.Bytes()
}

// Crete new exporter.
func newExporter(cfg component.Config, set exporter.CreateSettings, ctx context.Context) (*baseExporter, error) {
	oCfg := cfg.(*Config)
	if oCfg.Address != "" {
		_, err := url.Parse(oCfg.Address)
		if err != nil {
			return nil, errors.New("address must be a valid mqtt server")
		}
	}

	userAgent := fmt.Sprintf("%s/%s (%s/%s)",
		set.BuildInfo.Description, set.BuildInfo.Version, runtime.GOOS, runtime.GOARCH)

	// Client construction is deferred to start
	return &baseExporter{
		config:    oCfg,
		logger:    set.Logger,
		userAgent: userAgent,
		settings:  set.TelemetrySettings,
	}, nil
}

// start actually creates the MQTT client.
func (e *baseExporter) start(_ context.Context, _ component.Host) error {
	token := e.config.Client
	if token == nil {
		opts := mqtt.NewClientOptions().AddBroker(e.config.Address).SetClientID(e.config.Id)
		opts.SetUsername(e.config.Id)
		opts.SetPassword(e.config.Key)
		opts.SetKeepAlive(10 * time.Second)
		opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
			e.logger.Info("message on unknown channel, ignoring", zap.String("topic", message.Topic()), zap.ByteString("payload", message.Payload()))
		})
		opts.SetPingTimeout(5 * time.Second)
		opts.SetAutoReconnect(true)

		if e.config.TLS {
			opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		e.config.Client = &client
	}

	return nil
}

// inject attribute on all ScopeMetrics metrics
func (e *baseExporter) injectScopeMetricsAttribute(metricsScope pmetric.ScopeMetrics, attribute string, value string) pmetric.ScopeMetrics {
	metrics := metricsScope.Metrics()
	for i := 0; i < metrics.Len(); i++ {
		metricItem := metrics.At(i)

		switch metricItem.Type() {
		case pmetric.MetricTypeExponentialHistogram:
			for i := 0; i < metricItem.ExponentialHistogram().DataPoints().Len(); i++ {
				metricItem.ExponentialHistogram().DataPoints().At(i).Attributes().PutStr(attribute, value)
			}
		case pmetric.MetricTypeGauge:
			for i := 0; i < metricItem.Gauge().DataPoints().Len(); i++ {
				metricItem.Gauge().DataPoints().At(i).Attributes().PutStr(attribute, value)
			}
		case pmetric.MetricTypeHistogram:
			for i := 0; i < metricItem.Histogram().DataPoints().Len(); i++ {
				metricItem.Histogram().DataPoints().At(i).Attributes().PutStr(attribute, value)
			}
		case pmetric.MetricTypeSum:
			for i := 0; i < metricItem.Sum().DataPoints().Len(); i++ {
				metricItem.Sum().DataPoints().At(i).Attributes().PutStr(attribute, value)
			}
		case pmetric.MetricTypeSummary:
			for i := 0; i < metricItem.Summary().DataPoints().Len(); i++ {
				metricItem.Summary().DataPoints().At(i).Attributes().PutStr(attribute, value)
			}
		default:
			e.logger.Warn("not supported metric type", zap.String("name", metricItem.Name()),
				zap.String("type", metricItem.Type().String()))
			metrics.RemoveIf(func(m pmetric.Metric) bool {
				return m.Name() == metricItem.Name()
			})
		}
	}
	return metricsScope
}

// pushMetrics Exports metrics
func (e *baseExporter) pushMetrics(ctx context.Context, md pmetric.Metrics) error {
	tr := pmetricotlp.NewExportRequest()
	ref := tr.Metrics().ResourceMetrics().AppendEmpty()
	scopes := pmetricotlp.NewExportRequestFromMetrics(md).Metrics().ResourceMetrics().At(0).ScopeMetrics()
	for i := 0; i < scopes.Len(); i++ {
		scope := scopes.At(i)
		policyName, _ := scope.Scope().Attributes().Get("policy_name")
		policyNameStr := policyName.AsString()
		agentData, err := e.config.OrbAgentService.RetrieveAgentInfoByPolicyName(policyNameStr)
		if err != nil {
			e.logger.Warn("Policy is not managed by orb", zap.String("policyName", policyNameStr))
			continue
		}

		// sort datasetIDs to send always on same order
		datasetIDs := strings.Split(agentData.Datasets, ",")
		sort.Strings(datasetIDs)
		datasets := strings.Join(datasetIDs, ",")

		// Insert pivoted agentTags
		for key, value := range agentData.AgentTags {
			scope = e.injectScopeMetricsAttribute(scope, key, value)
		}
		// injecting policyID and datasetIDs attributes
		scope.Scope().Attributes().PutStr("policy_id", agentData.PolicyID)
		scope.Scope().Attributes().PutStr("dataset_ids", datasets)
		scope.CopyTo(ref.ScopeMetrics().AppendEmpty())
		e.logger.Info("scraped metrics for policy", zap.String("policy", policyNameStr), zap.String("policy_id", agentData.PolicyID))
	}

	request, err := tr.MarshalProto()
	if err != nil {
		defer ctx.Done()
		return consumererror.NewPermanent(err)
	}

	err = e.export(ctx, e.config.Topic, request)
	if err != nil {
		ctx.Done()
		return err
	}

	return err
}

// inject attribute on all ScopeLogs records
func (e *baseExporter) injectScopeLogsAttribute(logsScope plog.ScopeLogs, attribute string, value string) plog.ScopeLogs {
	logs := logsScope.LogRecords()
	for i := 0; i < logs.Len(); i++ {
		logItem := logs.At(i)
		logItem.Attributes().PutStr(attribute, value)
	}
	return logsScope
}

func (e *baseExporter) shutdown(_ context.Context) error {
	if e.config.Client != nil && (*e.config.Client).IsConnected() {
		(*e.config.Client).Disconnect(0)
	}
	return nil
}

func (e *baseExporter) pushLogs(ctx context.Context, ld plog.Logs) error {
	tr := plogotlp.NewExportRequest()
	ref := tr.Logs().ResourceLogs().AppendEmpty()
	scopes := plogotlp.NewExportRequestFromLogs(ld).Logs().ResourceLogs().At(0).ScopeLogs()
	for i := 0; i < scopes.Len(); i++ {
		scope := scopes.At(i)
		policyName := scope.Scope().Name()
		agentData, err := e.config.OrbAgentService.RetrieveAgentInfoByPolicyName(policyName)
		if err != nil {
			e.logger.Warn("Policy is not managed by orb", zap.String("policyName", policyName))
			continue
		}

		// sort datasetIDs to send always on same order
		datasetIDs := strings.Split(agentData.Datasets, ",")
		sort.Strings(datasetIDs)
		datasets := strings.Join(datasetIDs, ",")

		// Insert pivoted agentTags
		for key, value := range agentData.AgentTags {
			scope = e.injectScopeLogsAttribute(scope, key, value)
		}
		// injecting policyID and datasetIDs attributes
		scope.Scope().Attributes().PutStr("policy_id", agentData.PolicyID)
		scope.Scope().Attributes().PutStr("dataset_ids", datasets)
		scope.CopyTo(ref.ScopeLogs().AppendEmpty())
		e.logger.Info("scraped logs for policy", zap.String("policy", policyName), zap.String("policy_id", agentData.PolicyID))
	}

	request, err := tr.MarshalProto()
	if err != nil {
		defer ctx.Done()
		return consumererror.NewPermanent(err)
	}

	err = e.export(ctx, e.config.Topic, request)
	if err != nil {
		ctx.Done()
		return err
	}

	return err
}

// inject attribute on all ScopeSpans spans
func (e *baseExporter) injectScopeSpansAttribute(spanScope ptrace.ScopeSpans, attribute string, value string) ptrace.ScopeSpans {
	spans := spanScope.Spans()
	for i := 0; i < spans.Len(); i++ {
		spanItem := spans.At(i)
		spanItem.Attributes().PutStr(attribute, value)
	}
	return spanScope
}

func (e *baseExporter) pushTraces(ctx context.Context, td ptrace.Traces) error {
	tr := ptraceotlp.NewExportRequest()
	ref := tr.Traces().ResourceSpans().AppendEmpty()
	scopes := ptraceotlp.NewExportRequestFromTraces(td).Traces().ResourceSpans().At(0).ScopeSpans()
	for i := 0; i < scopes.Len(); i++ {
		scope := scopes.At(i)
		policyName := scope.Scope().Name()
		agentData, err := e.config.OrbAgentService.RetrieveAgentInfoByPolicyName(policyName)
		if err != nil {
			e.logger.Warn("Policy is not managed by orb", zap.String("policyName", policyName))
			continue
		}

		// sort datasetIDs to send always on same order
		datasetIDs := strings.Split(agentData.Datasets, ",")
		sort.Strings(datasetIDs)
		datasets := strings.Join(datasetIDs, ",")

		// Insert pivoted agentTags
		for key, value := range agentData.AgentTags {
			scope = e.injectScopeSpansAttribute(scope, key, value)
		}
		// injecting policyID and datasetIDs attributes
		scope.Scope().Attributes().PutStr("policy_id", agentData.PolicyID)
		scope.Scope().Attributes().PutStr("dataset_ids", datasets)
		scope.CopyTo(ref.ScopeSpans().AppendEmpty())
		e.logger.Info("scraped traces for policy", zap.String("policy", policyName), zap.String("policy_id", agentData.PolicyID))
	}

	request, err := tr.MarshalProto()
	if err != nil {
		defer ctx.Done()
		return consumererror.NewPermanent(err)
	}

	err = e.export(ctx, e.config.Topic, request)
	if err != nil {
		ctx.Done()
		return err
	}

	return err
}

func (e *baseExporter) export(ctx context.Context, topic string, request []byte) error {
	compressedPayload := e.compressBrotli(request)
	c := *e.config.Client
	if token := c.Publish(topic, 1, false, compressedPayload); token.Wait() && token.Error() != nil {
		e.logger.Error("error sending metrics RPC", zap.String("topic", topic), zap.Error(token.Error()))
		e.config.OrbAgentService.NotifyAgentDisconnection(ctx, token.Error())
		return token.Error()
	}
	e.logger.Debug("scraped and published telemetry", zap.String("topic", topic),
		zap.Int("payload_size_b", len(request)),
		zap.Int("compressed_payload_size_b", len(compressedPayload)))

	return nil
}
