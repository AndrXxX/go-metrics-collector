package agent

import (
	"github.com/AndrXxX/go-metrics-collector/internal/agent/config"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/executors"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/metrics"
	"github.com/AndrXxX/go-metrics-collector/internal/agent/utils"
	"net/http"
	"time"
)

func Run(config *config.Config) error {
	m := metrics.NewMetrics()
	scheduler := executors.NewIntervalScheduler(config.Intervals.SleepInterval)
	scheduler.Add(executors.NewCollector(&config.Metrics), time.Duration(config.Intervals.PollInterval)*time.Second)
	// TODO: Perederey Добавь логирование важных событий, например, успешный запуск агент и периодические задачи.

	rs := utils.NewRequestSender(utils.NewMetricURLBuilder(config.Common.Host), http.DefaultClient)
	scheduler.Add(executors.NewUploader(rs), time.Duration(config.Intervals.ReportInterval)*time.Second)

	err := scheduler.Run(*m)
	return err
}
