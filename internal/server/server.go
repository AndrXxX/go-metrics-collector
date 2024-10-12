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
	"github.com/go-chi/chi/v5/middleware"
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
	hg := hashgenerator.Factory().SHA256()

	r := chi.NewRouter()
	r.Mount("/debug", middleware.Profiler())
	r.Get("/ping", cFactory.From([]interfaces.Handler{
		dbping.New(dbchecker.New(a.storage.db)),
	}).Handler())

	r.Use(middlewares.CompressGzip().Handler)

	r.Route("/updates", func(r chi.Router) {
		r.Use(middlewares.HasCorrectSHA256HashOr500(hg, a.config.c.Key).Handler)
		r.Use(middlewares.SetContentType(contenttypes.ApplicationJSON).Handler)

		r.Post("/", cFactory.From([]interfaces.Handler{
			middlewares.AddSHA256HashHeader(hg, a.config.c.Key),
			updatemanymetrics.New(metricsupdater.New(a.storage.s)),
		}).Handler())
	})

	r.Route(fmt.Sprintf("/update/{%v}/{%v}/{%v}", vars.MetricType, vars.Metric, vars.Value), func(r chi.Router) {
		r.Use(middlewares.SetContentType(contenttypes.TextPlain).Handler)
		r.Use(middlewares.HasMetricOr404().Handler)

		r.Post("/", cFactory.From([]interfaces.Handler{
			middlewares.AddSHA256HashHeader(hg, a.config.c.Key),
			updatemetrics.New(metricsupdater.New(a.storage.s), metricsformatter.MetricsEmptyFormatter{}, metricsidentifier.NewURLIdentifier()),
		}).Handler())
	})

	r.Route("/update", func(r chi.Router) {
		r.Use(middlewares.SetContentType(contenttypes.ApplicationJSON).Handler)

		r.Post("/", cFactory.From([]interfaces.Handler{
			middlewares.AddSHA256HashHeader(hg, a.config.c.Key),
			updatemetrics.New(metricsupdater.New(a.storage.s), metricsformatter.MetricsJSONFormatter{}, metricsidentifier.NewJSONIdentifier()),
		}).Handler())
	})

	r.Route(fmt.Sprintf("/value/{%v}/{%v}", vars.MetricType, vars.Metric), func(r chi.Router) {
		r.Use(middlewares.SetContentType(contenttypes.TextPlain).Handler)
		r.Use(middlewares.HasMetricOr404().Handler)

		r.Get("/", cFactory.From([]interfaces.Handler{
			middlewares.AddSHA256HashHeader(hg, a.config.c.Key),
			fetchmetrics.New(a.storage.s, metricsformatter.MetricsValueFormatter{}, metricsidentifier.NewURLIdentifier(), mc),
		}).Handler())
	})

	r.Route("/value", func(r chi.Router) {
		r.Use(middlewares.SetContentType(contenttypes.ApplicationJSON).Handler)

		r.Post("/", cFactory.From([]interfaces.Handler{
			middlewares.AddSHA256HashHeader(hg, a.config.c.Key),
			fetchmetrics.New(a.storage.s, metricsformatter.MetricsJSONFormatter{}, metricsidentifier.NewJSONIdentifier(), mc),
		}).Handler())
	})

	r.Route("/", func(r chi.Router) {
		r.Use(middlewares.SetContentType(contenttypes.TextHTML).Handler)

		r.Get("/", cFactory.From([]interfaces.Handler{
			middlewares.AddSHA256HashHeader(hg, a.config.c.Key),
			fetchallmetrics.New(a.storage.s),
		}).Handler())
	})

	srv := &http.Server{Addr: a.config.c.Host, Handler: r}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("HTTP server ListenAndServe", zap.Error(err))
		}
	}()

	logger.Log.Info("listening", zap.String("host", a.config.c.Host))

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
