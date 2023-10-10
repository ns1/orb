/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sinker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/go-redis/redis/v8"
	mfnats "github.com/mainflux/mainflux/pkg/messaging/nats"
	fleetpb "github.com/orb-community/orb/fleet/pb"
	policiespb "github.com/orb-community/orb/policies/pb"
	"github.com/orb-community/orb/sinker/backend/pktvisor"
	"github.com/orb-community/orb/sinker/config"
	"github.com/orb-community/orb/sinker/otel"
	"github.com/orb-community/orb/sinker/otel/bridgeservice"
	"github.com/orb-community/orb/sinker/prometheus"
	sinkspb "github.com/orb-community/orb/sinks/pb"
	"go.uber.org/zap"
)

const (
	BackendMetricsTopic = "be.*.m.>"
	OtelMetricsTopic    = "otlp.*.m.>"
	MaxMsgPayloadSize   = 1048 * 1000
)

var (
	ErrPayloadTooBig = errors.New("payload too big")
	ErrNotFound      = errors.New("non-existent entity")
)

type Service interface {
	// Start set up communication with the message bus to communicate with agents
	Start() error
	// Stop end communication with the message bus
	Stop() error
}

type SinkerService struct {
	pubSub                 mfnats.PubSub
	otel                   bool
	otelMetricsCancelFunct context.CancelFunc
	otelLogsCancelFunct    context.CancelFunc
	otelKafkaUrl           string

	sinkerCache             config.ConfigRepo
	inMemoryCacheExpiration time.Duration
	esclient                *redis.Client
	logger                  *zap.Logger

	hbTicker *time.Ticker
	hbDone   chan bool

	promClient prometheus.Client

	policiesClient policiespb.PolicyServiceClient
	fleetClient    fleetpb.FleetServiceClient
	sinksClient    sinkspb.SinkServiceClient

	requestGauge   metrics.Gauge
	requestCounter metrics.Counter

	messageInputCounter metrics.Counter
	cancelAsyncContext  context.CancelFunc
	asyncContext        context.Context
}

func (svc SinkerService) Start() error {
	ctx := context.WithValue(context.Background(), "routine", "async")
	ctx = context.WithValue(ctx, "cache_expiry", svc.inMemoryCacheExpiration)
	svc.asyncContext, svc.cancelAsyncContext = context.WithCancel(ctx)
	if !svc.otel {
		topic := fmt.Sprintf("channels.*.%s", BackendMetricsTopic)
		if err := svc.pubSub.Subscribe(topic, svc.handleMsgFromAgent); err != nil {
			return err
		}
		svc.logger.Info("started metrics consumer", zap.String("topic", topic))
	}

	svc.hbTicker = time.NewTicker(CheckerFreq)
	svc.hbDone = make(chan bool)
	go svc.checkSinker()

	err := svc.startOtel(svc.asyncContext)
	if err != nil {
		svc.logger.Error("error on starting otel, exiting")
		return err
	}

	return nil
}

func (svc SinkerService) startOtel(ctx context.Context) error {
	if svc.otel {
		var err error

		bridgeService := bridgeservice.NewBridgeService(svc.logger, svc.inMemoryCacheExpiration, svc.sinkerCache,
			svc.policiesClient, svc.sinksClient, svc.fleetClient, svc.messageInputCounter)
		svc.otelMetricsCancelFunct, err = otel.StartOtelMetricsComponents(ctx, &bridgeService, svc.logger, svc.otelKafkaUrl, svc.pubSub)

		// starting Otel Logs components
		svc.otelLogsCancelFunct, err = otel.StartOtelLogsComponents(ctx, &bridgeService, svc.logger, svc.otelKafkaUrl, svc.pubSub)

		if err != nil {
			svc.logger.Error("error during StartOtelComponents", zap.Error(err))
			return err
		}
	}
	return nil
}

func (svc SinkerService) Stop() error {
	if svc.otel {
		otelTopic := fmt.Sprintf("channels.*.%s", OtelMetricsTopic)
		if err := svc.pubSub.Unsubscribe(otelTopic); err != nil {
			return err
		}
	} else {
		topic := fmt.Sprintf("channels.*.%s", BackendMetricsTopic)
		if err := svc.pubSub.Unsubscribe(topic); err != nil {
			return err
		}
	}

	svc.logger.Info("unsubscribed from agent metrics")

	svc.hbTicker.Stop()
	svc.hbDone <- true
	svc.cancelAsyncContext()

	return nil
}

// New instantiates the sinker service implementation.
func New(logger *zap.Logger,
	pubSub mfnats.PubSub,
	esclient *redis.Client,
	configRepo config.ConfigRepo,
	policiesClient policiespb.PolicyServiceClient,
	fleetClient fleetpb.FleetServiceClient,
	sinksClient sinkspb.SinkServiceClient,
	otelKafkaUrl string,
	enableOtel bool,
	requestGauge metrics.Gauge,
	requestCounter metrics.Counter,
	inputCounter metrics.Counter,
	defaultCacheExpiration time.Duration,
) Service {

	pktvisor.Register(logger)
	return &SinkerService{
		inMemoryCacheExpiration: defaultCacheExpiration,
		logger:                  logger,
		pubSub:                  pubSub,
		esclient:                esclient,
		sinkerCache:             configRepo,
		policiesClient:          policiesClient,
		fleetClient:             fleetClient,
		sinksClient:             sinksClient,
		requestGauge:            requestGauge,
		requestCounter:          requestCounter,
		messageInputCounter:     inputCounter,
		otel:                    enableOtel,
		otelKafkaUrl:            otelKafkaUrl,
	}
}
