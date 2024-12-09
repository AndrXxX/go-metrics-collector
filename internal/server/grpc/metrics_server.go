package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mp "github.com/AndrXxX/go-metrics-collector/pkg/metricsproto"

	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

// MetricsServer поддерживает все необходимые методы сервера.
type MetricsServer struct {
	mp.UnimplementedMetricsServer
	Updater updater
}

// UpdateMetrics реализует интерфейс обновления метрик.
func (s *MetricsServer) UpdateMetrics(ctx context.Context, in *mp.UpdateMetricsRequest) (*mp.UpdateMetricsResponse, error) {
	var response mp.UpdateMetricsResponse

	var list []models.Metrics
	for _, metric := range in.Metrics {
		list = append(list, models.Metrics{
			ID:    metric.Id,
			MType: metric.Type,
			Delta: &metric.Delta,
			Value: &metric.Value,
		})
	}
	err := s.Updater.UpdateMany(ctx, list)
	if err != nil {
		logger.Log.Error("Ошибка при обновлении метрик", zap.Error(err))
		return nil, status.Error(codes.Internal, "Ошибка при обновлении метрик")
	}
	return &response, nil
}
