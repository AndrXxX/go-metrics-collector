package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/contenttypes"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/dbping"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/fetchallmetrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/fetchmetrics"
	apilogger "github.com/AndrXxX/go-metrics-collector/internal/server/api/logger"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/middlewares"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/updatemetrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/conveyor"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricschecker"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsformatter"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsidentifier"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsupdater"
	"github.com/AndrXxX/go-metrics-collector/internal/server/tasks/savestoragetask"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type app struct {
	s  interfaces.MetricsStorage
	db *sql.DB
	c  *config.Config
}

func New(c *config.Config, s interfaces.MetricsStorage) *app {
	return &app{s, getDb(c), c}
}

func (a *app) Run() error {

	cFactory := conveyor.Factory(apilogger.New())
	mc := metricschecker.New()

	ctx, cancel := context.WithCancel(context.Background())

	if ss, ok := a.s.(repositories.StorageSaver); ok {
		sst := savestoragetask.New(time.Duration(a.c.StoreInterval)*time.Second, ss)
		go sst.Execute(ctx)
	}

	r := chi.NewRouter()
	r.Get("/ping", cFactory.From([]interfaces.Handler{
		dbping.New(a.db),
	}).Handler())
	r.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/{%v}/{%v}/{%v}", vars.MetricType, vars.Metric, vars.Value), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.HasMetricOr404(),
			updatemetrics.New(metricsupdater.New(a.s), metricsformatter.MetricsEmptyFormatter{}, metricsidentifier.NewURLIdentifier()),
		}).Handler())

		r.Post("/", cFactory.From([]interfaces.Handler{
			middlewares.CompressGzip(),
			middlewares.SetContentType(contenttypes.ApplicationJSON),
			updatemetrics.New(metricsupdater.New(a.s), metricsformatter.MetricsJSONFormatter{}, metricsidentifier.NewJSONIdentifier()),
		}).Handler())
	})

	r.Route("/value", func(r chi.Router) {
		r.Get(fmt.Sprintf("/{%v}/{%v}", vars.MetricType, vars.Metric), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.HasMetricOr404(),
			fetchmetrics.New(a.s, metricsformatter.MetricsValueFormatter{}, metricsidentifier.NewURLIdentifier(), mc),
		}).Handler())

		r.Post("/", cFactory.From([]interfaces.Handler{
			middlewares.CompressGzip(),
			middlewares.SetContentType(contenttypes.ApplicationJSON),
			fetchmetrics.New(a.s, metricsformatter.MetricsJSONFormatter{}, metricsidentifier.NewJSONIdentifier(), mc),
		}).Handler())
	})

	r.Get("/", cFactory.From([]interfaces.Handler{
		middlewares.CompressGzip(),
		middlewares.SetContentType(contenttypes.TextHTML),
		fetchallmetrics.New(a.s),
	}).Handler())

	srv := &http.Server{Addr: a.c.Host, Handler: r}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		cancel()

		if ss, ok := a.s.(repositories.StorageShutdowner); ok {
			err := ss.Shutdown()
			if err != nil {
				logger.Log.Error("Error on shutdown storage", zap.Error(err))
			}
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

func getDb(c *config.Config) *sql.DB {
	// TODO: Вынести в другое место
	db, err := sql.Open("pgx", c.DatabaseDSN)
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
