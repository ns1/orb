// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"context"
	"fmt"
	sinksgrpc "github.com/ns1labs/orb/sinks/api/grpc"
	"github.com/opentracing/opentracing-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ns1labs/orb/maestro"
	"github.com/ns1labs/orb/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	r "github.com/go-redis/redis/v8"
)

const (
	svcName   = "maestro"
	envPrefix = "orb_maestro"
	httpPort  = "8500"
)

func main() {

	esCfg := config.LoadEsConfig(envPrefix)
	svcCfg := config.LoadBaseServiceConfig(envPrefix, httpPort)
	jCfg := config.LoadJaegerConfig(envPrefix)
	sinksGRPCCfg := config.LoadGRPCConfig("orb", "sinks")

	// logger
	var logger *zap.Logger
	atomicLevel := zap.NewAtomicLevel()
	switch strings.ToLower(svcCfg.LogLevel) {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		os.Stdout,
		atomicLevel,
	)
	logger = zap.New(core, zap.AddCaller())
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	log := logger.Sugar()
	esClient := connectToRedis(esCfg.URL, esCfg.Pass, esCfg.DB, logger)
	defer func(esClient *r.Client) {
		err := esClient.Close()
		if err != nil {
			return
		}
	}(esClient)

	tracer, tracerCloser := initJaeger(svcName, jCfg.URL, logger)
	defer func(tracerCloser io.Closer) {
		err := tracerCloser.Close()
		if err != nil {
			logger.Fatal(err.Error())
		}
	}(tracerCloser)

	sinksGRPCConn := connectToGRPC(sinksGRPCCfg, logger)
	defer func(sinksGRPCConn *grpc.ClientConn) {
		err := sinksGRPCConn.Close()
		if err != nil {
			logger.Fatal(err.Error())
		}
	}(sinksGRPCConn)

	sinksGRPCTimeout, err := time.ParseDuration(sinksGRPCCfg.Timeout)
	if err != nil {
		log.Fatalf("Invalid %s value: %s", sinksGRPCCfg.Timeout, err.Error())
	}
	sinksGRPCClient := sinksgrpc.NewClient(tracer, sinksGRPCConn, sinksGRPCTimeout, logger)

	svc := maestro.NewMaestroService(logger, esClient, sinksGRPCClient, esCfg)
	errs := make(chan error, 2)

	mainContext, mainCancelFunction := context.WithCancel(context.Background())
	err = svc.Start(mainContext, mainCancelFunction)
	if err != nil {
		mainCancelFunction()
		log.Fatalf(fmt.Sprintf("Maestro service terminated: %s", err))
	}

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
		mainCancelFunction()
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("Maestro service terminated: %s", err))
}

func connectToGRPC(cfg config.GRPCConfig, logger *zap.Logger) *grpc.ClientConn {
	var opts []grpc.DialOption
	tls, err := strconv.ParseBool(cfg.ClientTLS)
	if err != nil {
		tls = false
	}
	if tls {
		if cfg.CaCerts != "" {
			tpc, err := credentials.NewClientTLSFromFile(cfg.CaCerts, "")
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to create tls credentials: %s", err))
				os.Exit(1)
			}
			opts = append(opts, grpc.WithTransportCredentials(tpc))
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(cfg.URL, opts...)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to dial to gRPC service %s: %s", cfg.URL, err))
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("Dialed to gRPC service %s at %s, TLS? %t", cfg.Service, cfg.URL, tls))

	return conn
}

func initJaeger(svcName, url string, logger *zap.Logger) (opentracing.Tracer, io.Closer) {
	if url == "" {
		return opentracing.NoopTracer{}, io.NopCloser(nil)
	}

	tracer, closer, err := jconfig.Configuration{
		ServiceName: svcName,
		Sampler: &jconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jconfig.ReporterConfig{
			LocalAgentHostPort: url,
			LogSpans:           true,
		},
	}.NewTracer()
	if err != nil {
		logger.Error("Failed to init Jaeger client", zap.Error(err))
		os.Exit(1)
	}

	return tracer, closer
}

func connectToRedis(redisURL, redisPass, redisDB string, logger *zap.Logger) *r.Client {
	db, err := strconv.Atoi(redisDB)
	if err != nil {
		logger.Error("Failed to connect to redis", zap.Error(err))
		os.Exit(1)
	}

	return r.NewClient(&r.Options{
		Addr:     redisURL,
		Password: redisPass,
		DB:       db,
	})
}
