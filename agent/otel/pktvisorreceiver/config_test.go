package pktvisorreceiver_test

import (
	"github.com/ns1labs/orb/agent/otel/pktvisorreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configtest"
	"path"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("load config", func(t *testing.T) {
		factories, err := componenttest.NopFactories()
		assert.NoError(t, err)

		factories.Receivers[typeStr] = pktvisorreceiver.NewFactory()
		factories.Exporters["prometheus"] = prometheusexporter.NewFactory()
		cfgPath := path.Join(".", "testdata", "config.yaml")
		cfg, err := configtest.LoadConfig(cfgPath, factories)

		require.NoError(t, err)
		require.NotNil(t, cfg)
	})
}
