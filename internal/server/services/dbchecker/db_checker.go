package dbchecker

import (
	"context"
	"database/sql"
	"fmt"
)

type dbChecker struct {
	db *sql.DB
}

func New(db *sql.DB) *dbChecker {
	return &dbChecker{db}
}

func (c *dbChecker) Check(ctx context.Context) error {
	if c.db == nil {
		return fmt.Errorf("db is not initialized")
	}
	err := c.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("error on ping db")
	}
	return nil
}
