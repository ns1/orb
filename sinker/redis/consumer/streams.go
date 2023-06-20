package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinker"
	"github.com/orb-community/orb/sinker/config"
	"go.uber.org/zap"
)

const (
	stream = "orb.sinks"
	group  = "orb.sinker"

	sinksPrefix = "sinks."
	sinksUpdate = sinksPrefix + "update"
	sinksCreate = sinksPrefix + "create"
	sinksDelete = sinksPrefix + "remove"
	exists      = "BUSYGROUP Consumer Group name already exists"
)

type Subscriber interface {
	Subscribe(context context.Context) error
}

type eventStore struct {
	otelEnabled   bool
	sinkerService sinker.Service
	configRepo    config.ConfigRepo
	client        *redis.Client
	esconsumer    string
	logger        *zap.Logger
}

func (es eventStore) Subscribe(context context.Context) error {
	subGroup := group
	if es.otelEnabled {
		subGroup = group + ".otel"
	}
	err := es.client.XGroupCreateMkStream(context, stream, subGroup, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := es.client.XReadGroup(context, &redis.XReadGroupArgs{
			Group:    subGroup,
			Consumer: es.esconsumer,
			Streams:  []string{stream, ">"},
			Count:    100,
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}

		for _, msg := range streams[0].Messages {
			event := msg.Values

			var err error
			switch event["operation"] {
			case sinksCreate:
				rte, derr := decodeSinksCreate(event)
				if derr != nil {
					err = derr
					break
				}
				err = es.handleSinksCreate(context, rte)
				if err != nil {
					es.logger.Error("Failed to handle event", zap.String("operation", event["operation"].(string)), zap.Error(err))
					break
				}
				es.client.XAck(context, stream, subGroup, msg.ID)
			case sinksUpdate:
				rte, derr := decodeSinksUpdate(event)
				if derr != nil {
					err = derr
					break
				}
				err = es.handleSinksUpdate(context, rte)
			case sinksDelete:
				rte, derr := decodeSinksRemove(event)
				if derr != nil {
					err = derr
					break
				}
				err = es.handleSinksRemove(context, rte)
			}
			if err != nil {
				es.logger.Error("Failed to handle event", zap.String("operation", event["operation"].(string)), zap.Error(err))
				continue
			}
			es.client.XAck(context, stream, subGroup, msg.ID)
		}
	}
}

// NewEventStore returns new event store instance.
func NewEventStore(sinkerService sinker.Service, configRepo config.ConfigRepo, client *redis.Client, esconsumer string, log *zap.Logger) Subscriber {
	return eventStore{
		sinkerService: sinkerService,
		configRepo:    configRepo,
		client:        client,
		esconsumer:    esconsumer,
		logger:        log,
	}
}

func decodeSinksCreate(event map[string]interface{}) (updateSinkEvent, error) {
	val := updateSinkEvent{
		sinkID:    read(event, "sink_id", ""),
		owner:     read(event, "owner", ""),
		timestamp: time.Time{},
	}
	var metadata types.Metadata
	if err := json.Unmarshal([]byte(read(event, "config", "")), &metadata); err != nil {
		return updateSinkEvent{}, err
	}
	val.config = metadata
	return val, nil
}

func decodeSinksUpdate(event map[string]interface{}) (updateSinkEvent, error) {
	val := updateSinkEvent{
		sinkID:    read(event, "sink_id", ""),
		owner:     read(event, "owner", ""),
		timestamp: time.Time{},
	}
	var metadata types.Metadata
	if err := json.Unmarshal([]byte(read(event, "config", "")), &metadata); err != nil {
		return updateSinkEvent{}, err
	}
	val.config = metadata
	return val, nil
}

func decodeSinksRemove(event map[string]interface{}) (updateSinkEvent, error) {
	val := updateSinkEvent{
		sinkID:    read(event, "sink_id", ""),
		owner:     read(event, "owner", ""),
		timestamp: time.Time{},
	}
	return val, nil
}

func (es eventStore) handleSinksRemove(_ context.Context, e updateSinkEvent) error {
	if ok := es.configRepo.Exists(e.owner, e.sinkID); ok {
		err := es.configRepo.Remove(e.owner, e.sinkID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (es eventStore) handleSinksUpdate(_ context.Context, e updateSinkEvent) error {
	data, err := json.Marshal(e.config)
	if err != nil {
		return err
	}
	var cfgParser config.SinkConfigParser
	var cfg config.SinkConfig
	if err := json.Unmarshal(data, &cfgParser); err != nil {
		return err
	}

	if ok := es.configRepo.Exists(e.owner, e.sinkID); ok {
		sinkConfig, err := es.configRepo.Get(e.owner, e.sinkID)
		if err != nil {
			return err
		}
		sinkConfig.AuthType = cfgParser.Authentication.Type
		sinkConfig.User = cfgParser.Authentication.Username
		sinkConfig.Password = cfgParser.Authentication.Password
		sinkConfig.Url = cfgParser.Exporter.RemoteHost
		sinkConfig.Opentelemetry = cfgParser.OpenTelemetry
		if sinkConfig.OwnerID == "" {
			sinkConfig.OwnerID = e.owner
		}
		if sinkConfig.SinkID == "" {
			sinkConfig.SinkID = e.sinkID
		}
		err = es.configRepo.Edit(sinkConfig)
		if err != nil {
			return err
		}
	} else {
		cfg.SinkID = e.sinkID
		cfg.OwnerID = e.owner
		cfg.AuthType = cfgParser.Authentication.Type
		cfg.User = cfgParser.Authentication.Username
		cfg.Password = cfgParser.Authentication.Password
		cfg.Url = cfgParser.Exporter.RemoteHost
		cfg.Opentelemetry = cfgParser.OpenTelemetry
		err = es.configRepo.Add(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (es eventStore) handleSinksCreate(_ context.Context, e updateSinkEvent) error {
	data, err := json.Marshal(e.config)
	if err != nil {
		return err
	}
	var cfgParser config.SinkConfigParser
	var cfg config.SinkConfig
	if err := json.Unmarshal(data, &cfgParser); err != nil {
		return err
	}
	cfg.SinkID = e.sinkID
	cfg.OwnerID = e.owner
	cfg.State = config.Unknown
	cfg.AuthType = cfgParser.Authentication.Type
	cfg.User = cfgParser.Authentication.Username
	cfg.Password = cfgParser.Authentication.Password
	cfg.Url = cfgParser.Exporter.RemoteHost
	cfg.Opentelemetry = cfgParser.OpenTelemetry
	err = es.configRepo.Add(cfg)
	if err != nil {
		return err
	}

	return nil
}

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}
	return val
}
