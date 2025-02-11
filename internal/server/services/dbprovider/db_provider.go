package dbprovider

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
)

type dbProvider struct {
	c *config.Config
}

// New возвращает сервис dbProvider для предоставления соединения с базой данных
func New(c *config.Config) *dbProvider {
	return &dbProvider{c}
}

// DB возвращает соединение с базой данных
func (p *dbProvider) DB() (*sql.DB, error) {
	if p.c.DatabaseDSN == "" {
		return nil, fmt.Errorf("empty DatabaseDSN")
	}
	db, err := sql.Open("pgx", p.c.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("error opening db %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("error on goose SetDialect %w", err)
	}

	if err := goose.Up(db, "internal/server/migrations/postgresql"); err != nil {
		return nil, fmt.Errorf("error on up migrations %w", err)
	}
	return db, nil
}
