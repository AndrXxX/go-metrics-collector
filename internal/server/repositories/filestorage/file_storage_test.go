package filestorage

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
	"github.com/AndrXxX/go-metrics-collector/internal/services/utils"
)

func Test_fileStorage_All(t *testing.T) {
	tests := []struct {
		name  string
		exist map[string]*models.Metrics
		want  map[string]*models.Metrics
	}{
		{
			name: "Test with two exist values",
			exist: map[string]*models.Metrics{
				"test1": {ID: "test1", MType: "type1", Value: utils.Pointer[float64](10.1)},
				"test2": {ID: "test2", MType: "type2", Delta: utils.Pointer[int64](10)},
			},
			want: map[string]*models.Metrics{
				"test1": {ID: "test1", MType: "type1", Value: utils.Pointer[float64](10.1)},
				"test2": {ID: "test2", MType: "type2", Delta: utils.Pointer[int64](10)},
			},
		},
		{
			name:  "Test with empty list",
			exist: map[string]*models.Metrics{},
			want:  map[string]*models.Metrics{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ms := memory.New[*models.Metrics]()
			for _, m := range tt.exist {
				ms.Insert(ctx, m.ID, m)
			}
			s := New(&ms)
			assert.EqualValues(t, tt.want, s.All(ctx))
		})
	}
}

func Test_fileStorage_Delete(t *testing.T) {
	tests := []struct {
		name        string
		exist       map[string]*models.Metrics
		forDeleteID string
		want        map[string]*models.Metrics
	}{
		{
			name: "Test with two exist values",
			exist: map[string]*models.Metrics{
				"test1": {ID: "test1", MType: "type1", Value: utils.Pointer[float64](10.1)},
				"test2": {ID: "test2", MType: "type2", Delta: utils.Pointer[int64](10)},
			},
			forDeleteID: "test1",
			want: map[string]*models.Metrics{
				"test2": {ID: "test2", MType: "type2", Delta: utils.Pointer[int64](10)},
			},
		},
		{
			name:  "Test with empty list",
			exist: map[string]*models.Metrics{},
			want:  map[string]*models.Metrics{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ms := memory.New[*models.Metrics]()
			for _, m := range tt.exist {
				ms.Insert(ctx, m.ID, m)
			}
			s := New(&ms)
			s.Delete(ctx, tt.forDeleteID)
			assert.EqualValues(t, tt.want, s.All(ctx))
		})
	}
}

func Test_fileStorage_Get(t *testing.T) {
	tests := []struct {
		name     string
		exist    map[string]*models.Metrics
		forGetID string
		wantM    *models.Metrics
		wantOk   bool
	}{
		{
			name: "Get item when it exists",
			exist: map[string]*models.Metrics{
				"test1": {ID: "test1", MType: "type1", Value: utils.Pointer[float64](10.1)},
				"test2": {ID: "test2", MType: "type2", Delta: utils.Pointer[int64](10)},
			},
			forGetID: "test2",
			wantM:    &models.Metrics{ID: "test2", MType: "type2", Delta: utils.Pointer[int64](10)},
			wantOk:   true,
		},
		{
			name: "Get item when it doesn't exist",
			exist: map[string]*models.Metrics{
				"test1": {ID: "test1", MType: "type1", Value: utils.Pointer[float64](10.1)},
				"test2": {ID: "test2", MType: "type2", Delta: utils.Pointer[int64](10)},
			},
			forGetID: "test3",
			wantM:    nil,
			wantOk:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ms := memory.New[*models.Metrics]()
			for _, m := range tt.exist {
				ms.Insert(ctx, m.ID, m)
			}
			s := New(&ms)
			m, ok := s.Get(ctx, tt.forGetID)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantM, m)
		})
	}
}

func Test_fileStorage_Insert(t *testing.T) {
	tests := []struct {
		name      string
		exist     map[string]*models.Metrics
		forInsert *models.Metrics
		want      map[string]*models.Metrics
	}{
		{
			name: "Test with two exist values",
			exist: map[string]*models.Metrics{
				"test1": {ID: "test1", MType: "type1", Value: utils.Pointer[float64](10.1)},
				"test2": {ID: "test2", MType: "type2", Delta: utils.Pointer[int64](10)},
			},
			forInsert: &models.Metrics{ID: "test3", MType: "type3", Value: utils.Pointer[float64](10.5)},
			want: map[string]*models.Metrics{
				"test1": {ID: "test1", MType: "type1", Value: utils.Pointer[float64](10.1)},
				"test2": {ID: "test2", MType: "type2", Delta: utils.Pointer[int64](10)},
				"test3": {ID: "test3", MType: "type3", Value: utils.Pointer[float64](10.5)},
			},
		},
		{
			name:      "Test with empty list",
			forInsert: &models.Metrics{ID: "test3", MType: "type3", Delta: utils.Pointer[int64](3)},
			want: map[string]*models.Metrics{
				"test3": {ID: "test3", MType: "type3", Delta: utils.Pointer[int64](3)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ms := memory.New[*models.Metrics]()
			for _, m := range tt.exist {
				ms.Insert(ctx, m.ID, m)
			}
			s := New(&ms)
			s.Insert(ctx, tt.forInsert.ID, tt.forInsert)
			assert.EqualValues(t, tt.want, s.All(ctx))
		})
	}
}

func Test_fileStorage_Save(t *testing.T) {
	tests := []struct {
		name    string
		ss      storageSaver
		wantErr bool
	}{
		{
			name:    "Test when ss is nil",
			ss:      nil,
			wantErr: false,
		},
		{
			name:    "Test with error on save",
			ss:      &testSS{err: fmt.Errorf("some error")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(nil, WithStorageSaver(&config.Config{}, tt.ss))
			require.Equal(t, tt.wantErr, s.Save(context.Background()) != nil)
		})
	}
}

func Test_fileStorage_Shutdown(t *testing.T) {
	tests := []struct {
		name    string
		ss      storageSaver
		wantErr bool
	}{
		{
			name:    "Test when ss is nil",
			ss:      nil,
			wantErr: false,
		},
		{
			name:    "Test with error on save",
			ss:      &testSS{err: fmt.Errorf("some error")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(nil, WithStorageSaver(&config.Config{}, tt.ss))
			require.Equal(t, tt.wantErr, s.Shutdown(context.Background()) != nil)
		})
	}
}
