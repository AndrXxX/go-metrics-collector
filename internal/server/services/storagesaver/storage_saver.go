package storagesaver

import (
	"bufio"
	"encoding/json"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"go.uber.org/zap"
	"os"
)

type storageSaver struct {
	path string
}

func (ss *storageSaver) Save(s storage[*models.Metrics]) error {
	file, err := os.OpenFile(ss.path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Log.Error("Error on close file on save value", zap.Error(err))
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
		s.Insert(m.ID, &m)
	}
	encoder := json.NewEncoder(file)
	for _, value := range s.All() {
		err := encoder.Encode(&value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ss *storageSaver) Restore(s storage[*models.Metrics]) error {
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
	encoder := json.NewEncoder(file)
	for _, value := range s.All() {
		err := encoder.Encode(&value)
		if err != nil {
			logger.Log.Error("Error on encode value", zap.Error(err))
			continue
		}
	}
	return nil
}

func New(path string) *storageSaver {
	return &storageSaver{path}
}
