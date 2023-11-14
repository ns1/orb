// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package maestro

import (
	"context"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/orb-community/orb/maestro/deployment"
	"github.com/orb-community/orb/maestro/kubecontrol"
	"github.com/orb-community/orb/maestro/monitor"
	rediscons1 "github.com/orb-community/orb/maestro/redis/consumer"
	"github.com/orb-community/orb/maestro/redis/producer"
	"github.com/orb-community/orb/maestro/service"
	"github.com/orb-community/orb/pkg/config"
	sinkspb "github.com/orb-community/orb/sinks/pb"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var _ Service = (*maestroService)(nil)

type maestroService struct {
	serviceContext    context.Context
	serviceCancelFunc context.CancelFunc

	deploymentService   deployment.Service
	sinkListenerService rediscons1.SinksListener
	activityListener    rediscons1.SinkerActivityListener

	kubecontrol       kubecontrol.Service
	monitor           monitor.Service
	logger            *zap.Logger
	streamRedisClient *redis.Client
	sinkerRedisClient *redis.Client
	sinksClient       sinkspb.SinkServiceClient
	eventService      service.EventService
	esCfg             config.EsConfig
	kafkaUrl          string
}

func NewMaestroService(logger *zap.Logger, streamRedisClient *redis.Client, sinkerRedisClient *redis.Client,
	sinksGrpcClient sinkspb.SinkServiceClient, otelCfg config.OtelConfig, db *sqlx.DB, svcCfg config.BaseSvcConfig) Service {
	kubectr := kubecontrol.NewService(logger)
	repo := deployment.NewRepositoryService(db, logger)
	maestroProducer := producer.NewMaestroProducer(logger, streamRedisClient)
	deploymentService := deployment.NewDeploymentService(logger, repo, otelCfg.KafkaUrl, svcCfg.EncryptionKey, maestroProducer, kubectr)
	ps := producer.NewMaestroProducer(logger, streamRedisClient)
	monitorService := monitor.NewMonitorService(logger, &sinksGrpcClient, ps, &kubectr, deploymentService)
	eventService := service.NewEventService(logger, deploymentService, &sinksGrpcClient)
	eventService = service.NewTracingService(logger, eventService,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "maestro",
			Subsystem: "comms",
			Name:      "message_count",
			Help:      "Number of messages received.",
		}, []string{"method", "sink_id", "owner_id"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "maestro",
			Subsystem: "comms",
			Name:      "message_latency_microseconds",
			Help:      "Total duration of messages processed in microseconds.",
		}, []string{"method", "sink_id", "owner_id"}))
	sinkListenerService := rediscons1.NewSinksListenerController(logger, eventService, streamRedisClient, sinksGrpcClient)
	activityListener := rediscons1.NewSinkerActivityListener(logger, eventService, streamRedisClient)
	return &maestroService{
		logger:              logger,
		deploymentService:   deploymentService,
		streamRedisClient:   streamRedisClient,
		sinkerRedisClient:   sinkerRedisClient,
		sinksClient:         sinksGrpcClient,
		sinkListenerService: sinkListenerService,
		activityListener:    activityListener,
		kubecontrol:         kubectr,
		monitor:             monitorService,
		kafkaUrl:            otelCfg.KafkaUrl,
	}
}

// Start will load all sinks from DB using SinksGRPC,
//
//	then for each sink, will create DeploymentEntry in Redis
//	And for each sink with active state, deploy OtelCollector
func (svc *maestroService) Start(ctx context.Context, cancelFunction context.CancelFunc) error {

	svc.serviceContext = ctx
	svc.serviceCancelFunc = cancelFunction

	go svc.subscribeToSinksEvents(ctx)
	go svc.subscribeToSinkerIdleEvents(ctx)
	go svc.subscribeToSinkerActivityEvents(ctx)

	monitorCtx := context.WithValue(ctx, "routine", "monitor")
	err := svc.monitor.Start(monitorCtx, cancelFunction)
	if err != nil {
		svc.logger.Error("error during monitor routine start", zap.Error(err))
		cancelFunction()
		return err
	}
	svc.logger.Info("Maestro service started")

	return nil
}

func (svc *maestroService) Stop() {
	svc.serviceCancelFunc()
	svc.logger.Info("Maestro service stopped")
}

func (svc *maestroService) subscribeToSinksEvents(ctx context.Context) {
	if err := svc.sinkListenerService.SubscribeSinksEvents(ctx); err != nil {
		svc.logger.Error("Bootstrap service failed to subscribe to event sourcing", zap.Error(err))
	}
	svc.logger.Info("finished reading sinks events")
	ctx.Done()
}

func (svc *maestroService) subscribeToSinkerIdleEvents(ctx context.Context) {
	if err := svc.activityListener.SubscribeSinkerIdleEvents(ctx); err != nil {
		svc.logger.Error("Bootstrap service failed to subscribe to event sourcing", zap.Error(err))
	}
	svc.logger.Info("finished reading sinker_idle events")
}

func (svc *maestroService) subscribeToSinkerActivityEvents(ctx context.Context) {
	if err := svc.activityListener.SubscribeSinkerActivityEvents(ctx); err != nil {
		svc.logger.Error("Bootstrap service failed to subscribe to event sourcing", zap.Error(err))
	}
	svc.logger.Info("finished reading sinker_activity events")
}
