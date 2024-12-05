package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	mp "github.com/AndrXxX/go-metrics-collector/pkg/metricsproto"
)

type testGRPCServer struct {
	mp.UnimplementedMetricsServer
	err error
	srv *grpc.Server
	req *mp.UpdateMetricsRequest
}

func (s *testGRPCServer) UpdateMetrics(_ context.Context, req *mp.UpdateMetricsRequest) (*mp.UpdateMetricsResponse, error) {
	s.req = req
	return nil, s.err
}

func (s *testGRPCServer) start(host string) {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	s.srv = grpc.NewServer()
	mp.RegisterMetricsServer(s.srv, s)
	go func() {
		err = s.srv.Serve(listen)
	}()
}

func (s *testGRPCServer) stop() {
	s.srv.Stop()
}

func Test_metricsUpdater_Update(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		opts    []grpc.DialOption
		list    []*mp.Metric
		srv     *testGRPCServer
		wantErr bool
	}{
		{
			name:    "Test with error on create client",
			wantErr: true,
		},
		{
			name:    "Test with code error on try to connect",
			host:    ":23423",
			opts:    []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
			wantErr: true,
		},
		{
			name:    "Test with error on try to connect",
			host:    ":23423",
			opts:    []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
			srv:     &testGRPCServer{err: fmt.Errorf("test")},
			wantErr: true,
		},
		{
			name:    "Test with succeed connection",
			host:    ":23423",
			opts:    []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
			srv:     &testGRPCServer{},
			wantErr: false,
		},
	}
	sortFunc := func(e, e2 *mp.Metric) int {
		return strings.Compare(strings.ToLower(e.Id), strings.ToLower(e2.Id))
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewGRPCMetricsUpdater(tt.host, tt.opts)
			if tt.srv != nil {
				tt.srv.start(tt.host)
				time.Sleep(50 * time.Millisecond)
			}
			err := u.Update(context.Background(), tt.list)
			if tt.srv != nil {
				tt.srv.stop()
				slices.SortFunc(tt.srv.req.Metrics, sortFunc)
				slices.SortFunc(tt.list, sortFunc)
				assert.EqualExportedValues(t, tt.list, tt.srv.req.Metrics)
			}
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
