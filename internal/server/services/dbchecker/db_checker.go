package dbchecker

import (
	"context"
	"fmt"
)

type dbChecker struct {
	db conn
}

// New возвращает сервис dbChecker для проверки соединения с базой данных
func New(db conn) *dbChecker {
	return &dbChecker{db}
}

// Check Проверяет соединение с базой данных
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
