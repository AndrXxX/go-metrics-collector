package metricsuploader

import (
	"context"
	"fmt"

	"github.com/AndrXxX/go-metrics-collector/internal/agent/dto"
	pb "github.com/AndrXxX/go-metrics-collector/internal/proto"
)

type grpcMetricsUploader struct {
	u grpcMetricsUpdater
}

func (c *grpcMetricsUploader) convert(result dto.MetricsDto) []*pb.Metric {
	var list []*pb.Metric
	for _, metric := range result.All() {
		var d int64
		var v float64
		if metric.Delta != nil {
			d = *metric.Delta
		}
		if metric.Value != nil {
			v = *metric.Value
		}
		list = append(list, &pb.Metric{
			Id:    metric.ID,
			Type:  metric.MType,
			Delta: d,
			Value: v,
		})
	}
	return list
}

// Process выполняет загрузку метрик
func (c *grpcMetricsUploader) Process(results <-chan dto.MetricsDto) error {
	for result := range results {
		list := c.convert(result)
		if len(list) == 0 {
			continue
		}
		err := c.u.Update(context.Background(), list)
		if err != nil {
			return fmt.Errorf("error on update: %w", err)
		}
	}
	return nil
}

// NewGRPCUploader возвращает сервис grpcMetricsUploader для загрузки метрик с помощью gRPC
func NewGRPCUploader(u grpcMetricsUpdater) *grpcMetricsUploader {
	return &grpcMetricsUploader{u}
}
