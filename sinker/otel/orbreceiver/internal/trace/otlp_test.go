// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace

import (
	"context"
	"errors"
	"github.com/ns1labs/orb/sinker/otel/orbreceiver/internal/testdata"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
)

func TestExport(t *testing.T) {
	td := testdata.GenerateTraces(1)
	// Keep trace data to compare the test result against it
	// Clone needed because OTLP proto XXX_ fields are altered in the GRPC downstream
	traceData := td.Clone()
	req := ptraceotlp.NewRequestFromTraces(td)

	traceSink := new(consumertest.TracesSink)
	traceClient := makeTraceServiceClient(t, traceSink)
	resp, err := traceClient.Export(context.Background(), req)
	require.NoError(t, err, "Failed to export trace: %v", err)
	require.NotNil(t, resp, "The response is missing")

	require.Len(t, traceSink.AllTraces(), 1)
	assert.EqualValues(t, traceData, traceSink.AllTraces()[0])
}

func TestExport_EmptyRequest(t *testing.T) {
	traceSink := new(consumertest.TracesSink)
	traceClient := makeTraceServiceClient(t, traceSink)
	resp, err := traceClient.Export(context.Background(), ptraceotlp.NewRequest())
	assert.NoError(t, err, "Failed to export trace: %v", err)
	assert.NotNil(t, resp, "The response is missing")
}

func TestExport_ErrorConsumer(t *testing.T) {
	td := testdata.GenerateTraces(1)
	req := ptraceotlp.NewRequestFromTraces(td)

	traceClient := makeTraceServiceClient(t, consumertest.NewErr(errors.New("my error")))
	resp, err := traceClient.Export(context.Background(), req)
	assert.EqualError(t, err, "rpc error: code = Unknown desc = my error")
	assert.Equal(t, ptraceotlp.Response{}, resp)
}

func makeTraceServiceClient(t *testing.T, tc consumer.Traces) ptraceotlp.Client {
	addr := otlpReceiverOnGRPCServer(t, tc)
	cc, err := grpc.Dial(addr.String(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	require.NoError(t, err, "Failed to create the TraceServiceClient: %v", err)
	t.Cleanup(func() {
		require.NoError(t, cc.Close())
	})

	return ptraceotlp.NewClient(cc)
}

func otlpReceiverOnGRPCServer(t *testing.T, tc consumer.Traces) net.Addr {
	ln, err := net.Listen("tcp", "localhost:")
	require.NoError(t, err, "Failed to find an available address to run the gRPC server: %v", err)

	t.Cleanup(func() {
		require.NoError(t, ln.Close())
	})

	r := New(config.NewComponentIDWithName("otlp", "trace"), tc, componenttest.NewNopReceiverCreateSettings())
	// Now run it as a gRPC server
	srv := grpc.NewServer()
	ptraceotlp.RegisterServer(srv, r)
	go func() {
		_ = srv.Serve(ln)
	}()

	return ln.Addr()
}
