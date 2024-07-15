package dbprovider

import (
	"database/sql"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type dbProvider struct {
	c *config.Config
}

func New(c *config.Config) *dbProvider {
	return &dbProvider{c}
}

func (p *dbProvider) Db() *sql.DB {
	if p.c.DatabaseDSN == "" {
		return nil
	}
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

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Log.Error("Error on goose SetDialect", zap.Error(err))
		return nil
	}

	if err := goose.Up(db, "internal/server/migrations/postgresql"); err != nil {
		logger.Log.Error("Error on up migrations", zap.Error(err))
		return nil
	}
	return db
}
