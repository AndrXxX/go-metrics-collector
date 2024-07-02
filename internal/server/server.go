package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/contenttypes"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/fetchallmetrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/fetchmetrics"
	apilogger "github.com/AndrXxX/go-metrics-collector/internal/server/api/logger"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/middlewares"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/updatemetrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/conveyor"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsidentifier"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricstringifier"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsupdater"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/storagesaver"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
)

func Run(c *config.Config) error {
	storage := memory.New[*models.Metrics]()
	ss := storagesaver.New(c.FileStoragePath)
	if c.Restore {
		err := ss.Restore(&storage)
		if err != nil {
			logger.Log.Error("Error restoring storage", zap.Error(err))
		}
	}
	cFactory := conveyor.Factory(apilogger.New())

	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/{%v}/{%v}/{%v}", vars.MetricType, vars.Metric, vars.Value), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.HasMetricOr404(),
			updatemetrics.New(metricsupdater.New(&storage), metricstringifier.MetricsEmptyStringifier{}, metricsidentifier.NewURLIdentifier()),
		}).Handler())

		r.Post("/", cFactory.From([]interfaces.Handler{
			middlewares.CompressGzip(),
			middlewares.SetContentType(contenttypes.ApplicationJSON),
			updatemetrics.New(metricsupdater.New(&storage), metricstringifier.MetricsJSONStringifier{}, metricsidentifier.NewJSONIdentifier()),
		}).Handler())
	})

	r.Route("/value", func(r chi.Router) {
		r.Get(fmt.Sprintf("/{%v}/{%v}", vars.MetricType, vars.Metric), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.HasMetricOr404(),
			fetchmetrics.New(&storage, metricstringifier.MetricsValueStringifier{}, metricsidentifier.NewURLIdentifier()),
		}).Handler())

		r.Post("/", cFactory.From([]interfaces.Handler{
			middlewares.CompressGzip(),
			middlewares.SetContentType(contenttypes.ApplicationJSON),
			fetchmetrics.New(&storage, metricstringifier.MetricsJSONStringifier{}, metricsidentifier.NewJSONIdentifier()),
		}).Handler())
	})

	r.Get("/", cFactory.From([]interfaces.Handler{
		middlewares.CompressGzip(),
		middlewares.SetContentType(contenttypes.TextHTML),
		fetchallmetrics.New(&storage),
	}).Handler())

	srv := &http.Server{Addr: c.Host, Handler: r}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		err := ss.Save(&storage)
		if err != nil {
			logger.Log.Error("Error on save storage", zap.Error(err))
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Log.Info("HTTP server Shutdown", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Info("HTTP server ListenAndServe", zap.Error(err))
	}

	<-idleConnsClosed

	return nil
}
