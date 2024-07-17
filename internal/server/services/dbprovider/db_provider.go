package dbprovider

import (
	"database/sql"
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type dbProvider struct {
	c *config.Config
}

func New(c *config.Config) *dbProvider {
	return &dbProvider{c}
}

func (p *dbProvider) DB() (error, *sql.DB) {
	if p.c.DatabaseDSN == "" {
		return fmt.Errorf("empty DatabaseDSN"), nil
	}
	db, err := sql.Open("pgx", p.c.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("error opening db %w", err), nil
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("error on goose SetDialect %w", err), nil
	}

	if err := goose.Up(db, "internal/server/migrations/postgresql"); err != nil {
		return fmt.Errorf("error on up migrations %w", err), nil
	}
	return nil, db
}
