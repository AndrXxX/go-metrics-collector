package dbprovider

import (
	"database/sql"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"go.uber.org/zap"
)

type dbProvider struct {
	c *config.Config
}

func New(c *config.Config) *dbProvider {
	return &dbProvider{c}
}

func (p *dbProvider) Db() *sql.DB {
	db, err := sql.Open("pgx", p.c.DatabaseDSN)
	if err != nil {
		logger.Log.Error("Error opening db", zap.Error(err))
		return nil
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Log.Error("Error closing db", zap.Error(err))
		}
	}(db)
	return db
}
