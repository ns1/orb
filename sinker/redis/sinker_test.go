package redis_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/mainflux/mainflux/pkg/uuid"
	"github.com/orb-community/orb/pkg/errors"
	config2 "github.com/orb-community/orb/sinker/config"
	"github.com/orb-community/orb/sinker/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var idProvider = uuid.New()

func TestSinkerConfigSave(t *testing.T) {
	sinkerCache := redis.NewSinkerCache(redisClient, logger)
	var config config2.SinkConfig
	config.SinkID = "123"
	config.OwnerID = "test"
	config.Authentication.Type = "basic_auth"
	config.Authentication.Username = "user"
	config.Authentication.Password = "password"
	config.Exporter.RemoteHost = "localhost"
	config.State = 0
	config.Msg = ""
	config.LastRemoteWrite = time.Time{}

	err := sinkerCache.Add(config)
	require.Nil(t, err, fmt.Sprintf("save sinker config to cache: expected nil got %s", err))

	cases := map[string]struct {
		config config2.SinkConfig
		err    error
	}{
		"Save sinker to cache": {
			config: config2.SinkConfig{
				SinkID:          "124",
				OwnerID:         "test",
				Exporter:        config.Exporter,
				Authentication:  config.Authentication,
				State:           0,
				Msg:             "",
				LastRemoteWrite: time.Time{},
			},
			err: nil,
		},
		"Save already cached sinker config to cache": {
			config: config,
			err:    nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			err := sinkerCache.Add(tc.config)
			assert.Nil(t, err, fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestGetSinkerConfig(t *testing.T) {
	sinkerCache := redis.NewSinkerCache(redisClient, logger)
	var config config2.SinkConfig
	config.SinkID = "123"
	config.OwnerID = "test"
	config.Authentication.Type = "basic_auth"
	config.Authentication.Username = "user"
	config.Authentication.Password = "password"
	config.Exporter.RemoteHost = "localhost"
	config.State = 0
	config.Msg = ""
	config.LastRemoteWrite = time.Time{}

	err := sinkerCache.Add(config)
	require.Nil(t, err, fmt.Sprintf("save sinker config to cache: expected nil got %s", err))

	cases := map[string]struct {
		sinkID string
		config config2.SinkConfig
		err    error
	}{
		"Get Config by existing sinker-key": {
			sinkID: "123",
			config: config,
			err:    nil,
		},
		"Get Config by non-existing sinker-key": {
			sinkID: "000",
			config: config2.SinkConfig{},
			err:    errors.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			sinkConfig, err := sinkerCache.Get(tc.config.OwnerID, tc.sinkID)
			assert.True(t, reflect.DeepEqual(tc.config, sinkConfig), fmt.Sprintf("%s: expected %v got %v", desc, tc.config, sinkConfig))
			assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s", desc, tc.err, err))
		})
	}
}

func TestGetAllSinkerConfig(t *testing.T) {
	sinkerCache := redis.NewSinkerCache(redisClient, logger)
	var config config2.SinkConfig
	config.SinkID = "123"
	config.OwnerID = "test"
	config.Authentication.Type = "basic_auth"
	config.Authentication.Username = "user"
	config.Authentication.Password = "password"
	config.Exporter.RemoteHost = "localhost"
	config.State = 0
	config.Msg = ""
	config.LastRemoteWrite = time.Time{}
	sinksConfig := map[string]struct {
		config config2.SinkConfig
	}{
		"config 1": {
			config: config2.SinkConfig{
				SinkID:          "123",
				OwnerID:         "test",
				Exporter:        config.Exporter,
				Authentication:  config.Authentication,
				State:           0,
				Msg:             "",
				LastRemoteWrite: time.Time{},
			},
		},
		"config 2": {
			config: config2.SinkConfig{
				SinkID:          "134",
				OwnerID:         "test",
				Exporter:        config.Exporter,
				Authentication:  config.Authentication,
				State:           0,
				Msg:             "",
				LastRemoteWrite: time.Time{},
			},
		},
	}

	for _, val := range sinksConfig {
		err := sinkerCache.Add(val.config)
		require.Nil(t, err, fmt.Sprintf("save sinker config to cache: expected nil got %s", err))
	}

	cases := map[string]struct {
		size    int
		ownerID string
		err     error
	}{
		"Get Config by existing sinker-key": {
			size:    2,
			ownerID: "test",
			err:     nil,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			sinksConfig, err := sinkerCache.GetAll(tc.ownerID)
			assert.Nil(t, err, fmt.Sprintf("%s: unexpected error: %s", desc, err))
			assert.GreaterOrEqual(t, len(sinksConfig), tc.size, fmt.Sprintf("%s: expected %d got %d", desc, tc.size, len(sinksConfig)))
		})
	}
}
