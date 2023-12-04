package consumer

import (
	"context"
	redis2 "github.com/orb-community/orb/sinks/redis"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/orb-community/orb/sinks"
	"go.uber.org/zap"
)

const (
	stream = "orb.sinker"
	group  = "orb.sinks"

	sinkerPrefix = "sinker."
	sinkerUpdate = sinkerPrefix + "update"

	exists = "BUSYGROUP Consumer Group name already exists"
)

type Subscriber interface {
	Subscribe(context context.Context) error
}

type eventStore struct {
	sinkService sinks.SinkService
	client      *redis.Client
	esconsumer  string
	logger      *zap.Logger
}

func NewEventStore(sinkService sinks.SinkService, client *redis.Client, esconsumer string, logger *zap.Logger) Subscriber {
	return eventStore{
		sinkService: sinkService,
		client:      client,
		esconsumer:  esconsumer,
		logger:      logger,
	}
}

func (es eventStore) Subscribe(context context.Context) error {
	err := es.client.XGroupCreateMkStream(context, stream, group, "$").Err()
	if err != nil && err.Error() != exists {
		return err
	}

	for {
		streams, err := es.client.XReadGroup(context, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: es.esconsumer,
			Streams:  []string{stream, ">"},
			Count:    100,
		}).Result()
		if err != nil || len(streams) == 0 {
			continue
		}

		for _, msg := range streams[0].Messages {
			event := msg.Values
			es.logger.Info("received message in sinker event bus", zap.Any("operation", event["operation"]))
			var err error
			switch event["operation"] {
			case sinkerUpdate:
				rte := es.decodeSinkerStateUpdate(event)
				err = es.handleSinkerStateUpdate(context, rte)
			}
			if err != nil {
				es.logger.Error("Failed to handle event", zap.String("operation", event["operation"].(string)), zap.Error(err))
				break
			}
			es.client.XAck(context, stream, group, msg.ID)
		}
	}
}

func (es eventStore) handleSinkerStateUpdate(ctx context.Context, event redis2.StateUpdateEvent) error {
	state := sinks.NewStateFromString(event.State)
	err := es.sinkService.ChangeSinkStateInternal(ctx, event.SinkID, event.Msg, event.OwnerID, state)
	if err != nil {
		return err
	}
	return nil
}

func (es eventStore) decodeSinkerStateUpdate(event map[string]interface{}) redis2.StateUpdateEvent {
	val := redis2.StateUpdateEvent{
		OwnerID:   read(event, "owner", ""),
		SinkID:    read(event, "sink_id", ""),
		Msg:       read(event, "msg", ""),
		Timestamp: time.Time{},
	}
	val.State = event["state"].(string)
	return val
}

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}

	return val
}
