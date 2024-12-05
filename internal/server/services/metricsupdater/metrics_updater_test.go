package metricsupdater

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
	"github.com/AndrXxX/go-metrics-collector/internal/services/utils"
)

func Test_metricsUpdater_Update(t *testing.T) {
	tests := []struct {
		name     string
		s        storage[*models.Metrics]
		exist    *models.Metrics
		newModel *models.Metrics
		want     *models.Metrics
	}{
		{
			name:     "Test when empty storage with gauge metric",
			newModel: &models.Metrics{ID: "test1", MType: metrics.Gauge, Value: utils.Pointer[float64](1.1)},
			want:     &models.Metrics{ID: "test1", MType: metrics.Gauge, Value: utils.Pointer[float64](1.1)},
		},
		{
			name:     "Test when counter metric exist",
			exist:    &models.Metrics{ID: "test1", MType: metrics.Counter, Delta: utils.Pointer[int64](10)},
			newModel: &models.Metrics{ID: "test1", MType: metrics.Counter, Delta: utils.Pointer[int64](11)},
			want:     &models.Metrics{ID: "test1", MType: metrics.Counter, Delta: utils.Pointer[int64](21)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := memory.New[*models.Metrics]()
			if tt.exist != nil {
				s.Insert(ctx, tt.exist.ID, tt.exist)
			}
			u := New(&s)
			_, err := u.Update(ctx, tt.newModel)
			require.NoError(t, err)
			res, _ := s.Get(ctx, tt.newModel.ID)
			assert.Equal(t, tt.want, res)
		})
	}
}

func Test_metricsUpdater_UpdateMany(t *testing.T) {
	tests := []struct {
		name  string
		exist []models.Metrics
		list  []models.Metrics
		want  []*models.Metrics
	}{
		{
			name:  "Test with empty storage",
			exist: []models.Metrics{},
			list: []models.Metrics{
				{ID: "test1", MType: metrics.Counter, Delta: utils.Pointer[int64](11)},
				{ID: "test2", MType: metrics.Gauge, Value: utils.Pointer[float64](1.1)},
			},
			want: []*models.Metrics{
				{ID: "test1", MType: metrics.Counter, Delta: utils.Pointer[int64](11)},
				{ID: "test2", MType: metrics.Gauge, Value: utils.Pointer[float64](1.1)},
			},
		},
		{
			name: "Test with not empty storage",
			exist: []models.Metrics{
				{ID: "test1", MType: metrics.Counter, Delta: utils.Pointer[int64](11)},
				{ID: "test2", MType: metrics.Gauge, Value: utils.Pointer[float64](1.2)},
			},
			list: []models.Metrics{
				{ID: "test1", MType: metrics.Counter, Delta: utils.Pointer[int64](5)},
				{ID: "test2", MType: metrics.Gauge, Value: utils.Pointer[float64](1.1)},
			},
			want: []*models.Metrics{
				{ID: "test1", MType: metrics.Counter, Delta: utils.Pointer[int64](16)},
				{ID: "test2", MType: metrics.Gauge, Value: utils.Pointer[float64](1.1)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := memory.New[*models.Metrics]()
			for _, m := range tt.exist {
				s.Insert(ctx, m.ID, &m)
			}
			u := New(&s)
			err := u.UpdateMany(ctx, tt.list)
			require.NoError(t, err)
			allExist := s.All(ctx)
			for _, m := range tt.want {
				assert.Equal(t, m, allExist[m.ID])
			}
		})
	}
}
