package consumer

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/orb-community/orb/sinker/redis/producer"
	"go.uber.org/zap"
)

type SinkerKeyExpirationListener interface {
	// SubscribeToKeyExpiration Listen to the sinker key expiration
	SubscribeToKeyExpiration(ctx context.Context) error
	// ReceiveMessage to be used to receive the message from the sinker key expiration, async
	ReceiveMessage(ctx context.Context, message string) error
}

type sinkerKeyExpirationListener struct {
	logger           *zap.Logger
	cacheRedisClient *redis.Client
	idleProducer     producer.SinkIdleProducer
}

func NewSinkerKeyExpirationListener(l *zap.Logger, cacheRedisClient *redis.Client, idleProducer producer.SinkIdleProducer) SinkerKeyExpirationListener {
	logger := l.Named("sinker_key_expiration_listener")
	return &sinkerKeyExpirationListener{logger: logger, cacheRedisClient: cacheRedisClient, idleProducer: idleProducer}
}

// SubscribeToKeyExpiration to be used to subscribe to the sinker key expiration
func (s *sinkerKeyExpirationListener) SubscribeToKeyExpiration(ctx context.Context) error {
	go func() {
		pubsub := s.cacheRedisClient.Subscribe(ctx, "__key*__:*")
		defer func(pubsub *redis.PubSub) {
			_ = pubsub.Close()
		}(pubsub)
		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				s.logger.Info("key expired", zap.String("key", msg.Payload))
				subCtx := context.WithValue(ctx, "msg", msg.Payload)
				err := s.ReceiveMessage(subCtx, msg.Payload)
				if err != nil {
					s.logger.Error("error receiving message", zap.Error(err))
					return
				}
			}
		}
	}()
	return nil
}

// ReceiveMessage to be used to receive the message from the sinker key expiration
func (s *sinkerKeyExpirationListener) ReceiveMessage(ctx context.Context, message string) error {
	// goroutine
	go func(msg string) {
		ownerID := message[16:52]
		sinkID := message[53:]
		event := producer.SinkIdleEvent{
			OwnerID: ownerID,
			SinkID:  sinkID,
			State:   "idle",
			Size:    "0",
		}
		_ = s.idleProducer.PublishSinkIdle(ctx, event)
	}(message)
	return nil
}
