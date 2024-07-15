package storagesaver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"os"
)

type storageSaver struct {
	path string
	s    storage[*models.Metrics]
}

func (ss *storageSaver) Save(ctx context.Context) error {
	file, err := os.OpenFile(ss.path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Log.Info("failed to close file", zap.Error(err))
		}
	}(file)

	bufWriter := bufio.NewWriter(file)
	encoder := json.NewEncoder(bufWriter)
	for _, value := range ss.s.All(ctx) {
		err := encoder.Encode(&value)
		if err != nil {
			logger.Log.Error("Error on encode value", zap.Error(err))
			continue
		}
	}
	if err := bufWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}
	return nil
}

func (ss *storageSaver) Restore(ctx context.Context) error {
	file, err := os.OpenFile(ss.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Log.Error("Error on close file on restore value", zap.Error(err))
		}
	}(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var m *models.Metrics
		err := json.Unmarshal(scanner.Bytes(), &m)
		if err != nil {
			logger.Log.Error("Error on unmarshall value", zap.Error(err))
			continue
		}
		ss.s.Insert(ctx, m.ID, m)
	}
	return nil
}

func New(path string, s storage[*models.Metrics]) *storageSaver {
	return &storageSaver{path, s}
}
