package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	mp "github.com/AndrXxX/go-metrics-collector/pkg/metricsproto"
)

type metricsUpdater struct {
	host string
	opts []grpc.DialOption
}

func NewGRPCMetricsUpdater(host string, opts []grpc.DialOption) *metricsUpdater {
	return &metricsUpdater{host: host, opts: opts}
}

func (u metricsUpdater) Update(ctx context.Context, list []*mp.Metric) error {
	conn, err := grpc.NewClient(u.host, u.opts...)
	if err != nil {
		return fmt.Errorf("grpc connection error: %w", err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	c := mp.NewMetricsClient(conn)
	req := &mp.UpdateMetricsRequest{Metrics: list}
	_, err = c.UpdateMetrics(ctx, req)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			return fmt.Errorf("grpc response error: %s %w", e.Code(), e.Err())
		}
	}
	return nil
}
