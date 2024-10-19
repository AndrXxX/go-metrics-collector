package storageprovider

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/dbstorage"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/filestorage"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
)

func TestNew(t *testing.T) {
	type args struct {
		c  *config.Config
		db *sql.DB
	}
	tests := []struct {
		name string
		args args
		want *storageProvider
	}{
		{
			name: "Test OK",
			args: args{c: &config.Config{}, db: nil},
			want: &storageProvider{&config.Config{}, nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, New(tt.args.c, tt.args.db))
		})
	}
}

func Test_storageProvider_Storage(t *testing.T) {
	dbStorage := dbstorage.New(nil, nil)
	mStorage := memory.New[*models.Metrics]()
	fsStorage := filestorage.New(&config.Config{FileStoragePath: "test"}, &mStorage)
	type fields struct {
		c  *config.Config
		db *sql.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   interfaces.MetricsStorage
	}{
		{
			name:   "DB Storage",
			fields: fields{&config.Config{DatabaseDSN: "test"}, nil},
			want:   &dbStorage,
		},
		{
			name:   "File Storage",
			fields: fields{&config.Config{FileStoragePath: "test"}, nil},
			want:   &fsStorage,
		},
		{
			name:   "Memory Storage",
			fields: fields{&config.Config{}, nil},
			want:   &mStorage,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := New(tt.fields.c, tt.fields.db)
			assert.Equal(t, tt.want, sp.Storage(context.Background()))
			time.Sleep(100 * time.Millisecond)
		})
	}
}
