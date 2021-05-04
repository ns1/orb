// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// adapted for Orb project

package main

import (
	"fmt"
	"github.com/ns1labs/orb/pkg/mainflux/consumers/writers/promsink"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/consumers"
	"github.com/mainflux/mainflux/consumers/writers/api"
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/ns1labs/orb/pkg/mainflux/transformers/passthrough"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	svcName = "prom-sink"

	defLogLevel      = "error"
	defNatsURL       = "nats://localhost:4222"
	defPort          = "8180"
	defDBHost        = "localhost"
	defDBPort        = "5432"
	defDBUser        = "mainflux"
	defDBPass        = "mainflux"
	defDB            = "mainflux"
	defDBSSLMode     = "disable"
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""
	defConfigPath    = "/config.toml"
	defContentType   = "application/senml+json"
	defTransformer   = "passthrough"

	envNatsURL       = "MF_NATS_URL"
	envLogLevel      = "MF_POSTGRES_WRITER_LOG_LEVEL"
	envPort          = "MF_POSTGRES_WRITER_PORT"
	envDBHost        = "MF_POSTGRES_WRITER_DB_HOST"
	envDBPort        = "MF_POSTGRES_WRITER_DB_PORT"
	envDBUser        = "MF_POSTGRES_WRITER_DB_USER"
	envDBPass        = "MF_POSTGRES_WRITER_DB_PASS"
	envDB            = "MF_POSTGRES_WRITER_DB"
	envDBSSLMode     = "MF_POSTGRES_WRITER_DB_SSL_MODE"
	envDBSSLCert     = "MF_POSTGRES_WRITER_DB_SSL_CERT"
	envDBSSLKey      = "MF_POSTGRES_WRITER_DB_SSL_KEY"
	envDBSSLRootCert = "MF_POSTGRES_WRITER_DB_SSL_ROOT_CERT"
	envConfigPath    = "MF_POSTGRES_WRITER_CONFIG_PATH"
	envContentType   = "MF_POSTGRES_WRITER_CONTENT_TYPE"
	envTransformer   = "MF_POSTGRES_WRITER_TRANSFORMER"
)

type config struct {
	natsURL     string
	logLevel    string
	port        string
	configPath  string
	contentType string
	transformer string
	//dbConfig    postgres.Config
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	pubSub, err := nats.NewPubSub(cfg.natsURL, "", logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to NATS: %s", err))
		os.Exit(1)
	}
	defer pubSub.Close()

	// prometheus connection: https://github.com/timescale/promscale/blob/master/docs/writing_to_promscale.md
	//db := connectToDB(cfg.dbConfig, logger)
	//defer db.Close()

	repo := newService( /*db, */ logger)
	t := passthrough.New()

	if err = consumers.Start(pubSub, repo, t, cfg.configPath, logger); err != nil {
		logger.Error(fmt.Sprintf("Failed to create promsink writer: %s", err))
	}

	errs := make(chan error, 2)

	go startHTTPServer(cfg.port, errs, logger)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("promsink writer service terminated: %s", err))
}

func loadConfig() config {
	//dbConfig := postgres.Config{
	//	Host:        mainflux.Env(envDBHost, defDBHost),
	//	Port:        mainflux.Env(envDBPort, defDBPort),
	//	User:        mainflux.Env(envDBUser, defDBUser),
	//	Pass:        mainflux.Env(envDBPass, defDBPass),
	//	Name:        mainflux.Env(envDB, defDB),
	//	SSLMode:     mainflux.Env(envDBSSLMode, defDBSSLMode),
	//	SSLCert:     mainflux.Env(envDBSSLCert, defDBSSLCert),
	//	SSLKey:      mainflux.Env(envDBSSLKey, defDBSSLKey),
	//	SSLRootCert: mainflux.Env(envDBSSLRootCert, defDBSSLRootCert),
	//}

	return config{
		natsURL:     mainflux.Env(envNatsURL, defNatsURL),
		logLevel:    mainflux.Env(envLogLevel, defLogLevel),
		port:        mainflux.Env(envPort, defPort),
		configPath:  mainflux.Env(envConfigPath, defConfigPath),
		contentType: mainflux.Env(envContentType, defContentType),
		transformer: mainflux.Env(envTransformer, defTransformer),
		//dbConfig:    dbConfig,
	}
}

//func connectToDB(dbConfig postgres.Config, logger logger.Logger) *sqlx.DB {
//	db, err := postgres.Connect(dbConfig)
//	if err != nil {
//		logger.Error(fmt.Sprintf("Failed to connect to Postgres: %s", err))
//		os.Exit(1)
//	}
//	return db
//}

func newService( /*db *sqlx.DB, */ logger logger.Logger) consumers.Consumer {
	zlog, _ := zap.NewProduction()
	svc := promsink.New(zlog)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "promsink",
			Subsystem: "message_writer",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "promsink",
			Subsystem: "message_writer",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)

	return svc
}

func startHTTPServer(port string, errs chan error, logger logger.Logger) {
	p := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("promsink writer service started, exposed port %s", port))
	errs <- http.ListenAndServe(p, api.MakeHandler(svcName))
}
