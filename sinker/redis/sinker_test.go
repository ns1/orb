package redis_test

import (
	"context"
	"fmt"

	"testing"
	"time"

	"github.com/orb-community/orb/sinker/redis/producer"

	"github.com/stretchr/testify/require"
)

func TestSinkActivityStoreAndMessage(t *testing.T) {
	// Create SinkActivityService
	sinkTTLSvc := producer.NewSinkerKeyService(logger, redisClient)
	sinkActivitySvc := producer.NewSinkActivityProducer(logger, redisClient, sinkTTLSvc)
	args := []struct {
		testCase string
		event    producer.SinkActivityEvent
	}{
		{
			testCase: "sink activity for new sink",
			event: producer.SinkActivityEvent{
				OwnerID:   "1",
				SinkID:    "1",
				State:     "active",
				Size:      "40",
				Timestamp: time.Now(),
			},
		},
		{
			testCase: "sink activity for existing sink",
			event: producer.SinkActivityEvent{
				OwnerID:   "1",
				SinkID:    "1",
				State:     "active",
				Size:      "55",
				Timestamp: time.Now(),
			},
		},
		{
			testCase: "sink activity for another new sink",
			event: producer.SinkActivityEvent{
				OwnerID:   "2",
				SinkID:    "1",
				State:     "active",
				Size:      "37",
				Timestamp: time.Now(),
			},
		},
	}
	for _, tt := range args {
		ctx := context.WithValue(context.Background(), "test_case", tt.testCase)
		err := sinkActivitySvc.PublishSinkActivity(ctx, tt.event)
		require.NoError(t, err, fmt.Sprintf("%s: unexpected error: %s", tt.testCase, err))
	}
	logger.Debug("debugging breakpoint")
}
