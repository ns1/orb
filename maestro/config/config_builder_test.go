package config

import (
	"context"
	"fmt"
	"testing"
)

func TestReturnConfigYamlFromSink(t *testing.T) {
	type args struct {
		in0            context.Context
		kafkaUrlConfig string
		sinkId         string
		sinkUrl        string
		sinkUsername   string
		sinkPassword   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "simple test", args: args{
			in0:            context.Background(),
			kafkaUrlConfig: "kafka:9092",
			sinkId:         "sink-id-222",
			sinkUrl:        "https://mysinkurl:9922",
			sinkUsername:   "1234123",
			sinkPassword:   "CarnivorousVulgaris",
		}, want: `---\nreceivers:\n  kafka:\n    brokers:\n    - kafka:9092\n    topic: otlp_metrics-sink-id-222\n    protocol_version: 2.0.0\nextensions:\n  pprof:\n    endpoint: 0.0.0.0:1888\n  basicauth/exporter:\n    client_auth:\n      username: 1234123\n      password: CarnivorousVulgaris\nexporters:\n  prometheusremotewrite:\n    endpoint: https://mysinkurl:9922\n    auth:\n      authenticator: basicauth/exporter\n  logging:\n    verbosity: detailed\n    sampling_initial: 5\n    sampling_thereafter: 50\nservice:\n  extensions:\n  - pprof\n  - basicauth/exporter\n  pipelines:\n    metrics:\n      receivers:\n      - kafka\n      exporters:\n      - prometheusremotewrite\n`,
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReturnConfigYamlFromSink(tt.args.in0, tt.args.kafkaUrlConfig, tt.args.sinkId, tt.args.sinkUrl, tt.args.sinkUsername, tt.args.sinkPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReturnConfigYamlFromSink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("%s\n", got)
			if got != tt.want {
				t.Errorf("ReturnConfigYamlFromSink() got = \n%v\n, want \n%v", got, tt.want)
			}
		})
	}
}
