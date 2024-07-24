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
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/updatemanymetrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/updatemetrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/conveyor"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/dbchecker"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricschecker"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsformatter"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsidentifier"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsupdater"
	"github.com/AndrXxX/go-metrics-collector/internal/services/hashgenerator"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

const shutdownTimeout = 5 * time.Second

type app struct {
	config struct {
		c *config.Config
	}
	storage struct {
		s  interfaces.MetricsStorage
		db *sql.DB
	}
}

func New(c *config.Config, s interfaces.MetricsStorage, db *sql.DB) *app {
	return &app{
		config: struct {
			c *config.Config
		}{c: c},
		storage: struct {
			s  interfaces.MetricsStorage
			db *sql.DB
		}{s: s, db: db},
	}
}

func (a *app) Run(commonCtx context.Context) error {

	cFactory := conveyor.Factory(apilogger.New())
	mc := metricschecker.New()
	hg, err := hashgenerator.New(a.config.c.Key)
	if err != nil {
		logger.Log.Error("Error on create hash generator", zap.Error(err))
	}

	r := chi.NewRouter()
	r.Get("/ping", cFactory.From([]interfaces.Handler{
		dbping.New(dbchecker.New(a.storage.db)),
	}).Handler())

	r.Route("/updates", func(r chi.Router) {
		r.Post("/", cFactory.From([]interfaces.Handler{
			middlewares.HasCorrectSHA256HashOr500(hg),
			middlewares.CompressGzip(),
			middlewares.SetContentType(contenttypes.ApplicationJSON),
			middlewares.AddSHA256HashHeader(hg),
			updatemanymetrics.New(metricsupdater.New(a.storage.s)),
		}).Handler())
	})

	r.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/{%v}/{%v}/{%v}", vars.MetricType, vars.Metric, vars.Value), cFactory.From([]interfaces.Handler{
			middlewares.HasCorrectSHA256HashOr500(hg),
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.HasMetricOr404(),
			middlewares.AddSHA256HashHeader(hg),
			updatemetrics.New(metricsupdater.New(a.storage.s), metricsformatter.MetricsEmptyFormatter{}, metricsidentifier.NewURLIdentifier()),
		}).Handler())

		r.Post("/", cFactory.From([]interfaces.Handler{
			middlewares.HasCorrectSHA256HashOr500(hg),
			middlewares.CompressGzip(),
			middlewares.SetContentType(contenttypes.ApplicationJSON),
			middlewares.AddSHA256HashHeader(hg),
			updatemetrics.New(metricsupdater.New(a.storage.s), metricsformatter.MetricsJSONFormatter{}, metricsidentifier.NewJSONIdentifier()),
		}).Handler())
	})

	r.Route("/value", func(r chi.Router) {
		r.Get(fmt.Sprintf("/{%v}/{%v}", vars.MetricType, vars.Metric), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.HasMetricOr404(),
			middlewares.AddSHA256HashHeader(hg),
			fetchmetrics.New(a.storage.s, metricsformatter.MetricsValueFormatter{}, metricsidentifier.NewURLIdentifier(), mc),
		}).Handler())

		r.Post("/", cFactory.From([]interfaces.Handler{
			middlewares.CompressGzip(),
			middlewares.SetContentType(contenttypes.ApplicationJSON),
			middlewares.AddSHA256HashHeader(hg),
			fetchmetrics.New(a.storage.s, metricsformatter.MetricsJSONFormatter{}, metricsidentifier.NewJSONIdentifier(), mc),
		}).Handler())
	})

	r.Get("/", cFactory.From([]interfaces.Handler{
		middlewares.CompressGzip(),
		middlewares.SetContentType(contenttypes.TextHTML),
		middlewares.AddSHA256HashHeader(hg),
		fetchallmetrics.New(a.storage.s),
	}).Handler())

	srv := &http.Server{Addr: a.config.c.Host, Handler: r}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("HTTP server ListenAndServe", zap.Error(err))
		}
	}()

	logger.Log.Info(fmt.Sprintf("listening on %s", a.config.c.Host))

	<-commonCtx.Done()
	logger.Log.Info("shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	shutdown := make(chan struct{}, 1)
	go func() {
		if ss, ok := a.storage.s.(repositories.StorageShutdowner); ok {
			err := ss.Shutdown(shutdownCtx)
			if err != nil {
				logger.Log.Error("Error on shutdown storage", zap.Error(err))
			}
		}
		if a.storage.db != nil {
			_ = a.storage.db.Close()
		}
		shutdown <- struct{}{}
	}()

	select {
	case <-shutdownCtx.Done():
		return fmt.Errorf("server shutdown: %w", shutdownCtx.Err())
	case <-shutdown:
		log.Println("finished")
	}

	return nil
}
