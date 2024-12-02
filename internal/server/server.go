package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/AndrXxX/go-metrics-collector/internal/enums/contenttypes"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	pb "github.com/AndrXxX/go-metrics-collector/internal/proto"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/dbping"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/fetchallmetrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/fetchmetrics"
	apilogger "github.com/AndrXxX/go-metrics-collector/internal/server/api/logger"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/middlewares"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/updatemanymetrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/updatemetrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	igrpc "github.com/AndrXxX/go-metrics-collector/internal/server/grpc"
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/dbchecker"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricschecker"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsformatter"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsidentifier"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsupdater"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/tlsconfig"
	"github.com/AndrXxX/go-metrics-collector/internal/services/hashgenerator"
	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
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

// New возвращает экземпляр приложения
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

// Run запускает сервер
func (a *app) Run(commonCtx context.Context) error {
	mc := metricschecker.New()
	hg := hashgenerator.Factory().SHA256()

	r := chi.NewRouter()

	r.Mount("/debug", middleware.Profiler())

	r.Route("/ping", func(r chi.Router) {
		r.Use(apilogger.New().Handler)
		r.Get("/", dbping.New(dbchecker.New(a.storage.db)).Handler())
	})

	r.Route("/updates", func(r chi.Router) {
		r.Use(apilogger.New().Handler)
		r.Use(middlewares.HasGrantedXRealIPOr403(a.config.c.TrustedSubnet).Handler)
		r.Use(middlewares.HasCorrectSHA256HashOr500(hg, a.config.c.Key).Handler)
		r.Use(middlewares.CompressGzip().Handler)
		r.Use(middlewares.SetContentType(contenttypes.ApplicationJSON).Handler)
		r.Use(middlewares.AddSHA256HashHeader(hg, a.config.c.Key).Handler)

		r.Post("/", updatemanymetrics.New(metricsupdater.New(a.storage.s)).Handler())
	})

	r.Route(fmt.Sprintf("/update/{%v}/{%v}/{%v}", vars.MetricType, vars.Metric, vars.Value), func(r chi.Router) {
		r.Use(apilogger.New().Handler)
		r.Use(middlewares.HasGrantedXRealIPOr403(a.config.c.TrustedSubnet).Handler)
		r.Use(middlewares.HasCorrectSHA256HashOr500(hg, a.config.c.Key).Handler)
		r.Use(middlewares.SetContentType(contenttypes.TextPlain).Handler)
		r.Use(middlewares.HasMetricOr404().Handler)
		r.Use(middlewares.AddSHA256HashHeader(hg, a.config.c.Key).Handler)

		updater := metricsupdater.New(a.storage.s)
		formatter := metricsformatter.MetricsEmptyFormatter{}
		identifier := metricsidentifier.NewURLIdentifier()
		r.Post("/", updatemetrics.New(updater, formatter, identifier).Handler())
	})

	r.Route("/update", func(r chi.Router) {
		r.Use(apilogger.New().Handler)
		r.Use(middlewares.HasGrantedXRealIPOr403(a.config.c.TrustedSubnet).Handler)
		r.Use(middlewares.HasCorrectSHA256HashOr500(hg, a.config.c.Key).Handler)
		r.Use(middlewares.CompressGzip().Handler)
		r.Use(middlewares.SetContentType(contenttypes.ApplicationJSON).Handler)
		r.Use(middlewares.AddSHA256HashHeader(hg, a.config.c.Key).Handler)

		updater := metricsupdater.New(a.storage.s)
		formatter := metricsformatter.MetricsJSONFormatter{}
		identifier := metricsidentifier.NewJSONIdentifier()
		r.Post("/", updatemetrics.New(updater, formatter, identifier).Handler())
	})

	r.Route(fmt.Sprintf("/value/{%v}/{%v}", vars.MetricType, vars.Metric), func(r chi.Router) {
		r.Use(apilogger.New().Handler)
		r.Use(middlewares.SetContentType(contenttypes.TextPlain).Handler)
		r.Use(middlewares.HasMetricOr404().Handler)
		r.Use(middlewares.AddSHA256HashHeader(hg, a.config.c.Key).Handler)

		formatter := metricsformatter.MetricsValueFormatter{}
		identifier := metricsidentifier.NewURLIdentifier()
		r.Get("/", fetchmetrics.New(a.storage.s, formatter, identifier, mc).Handler())
	})

	r.Route("/value", func(r chi.Router) {
		r.Use(apilogger.New().Handler)
		r.Use(middlewares.CompressGzip().Handler)
		r.Use(middlewares.SetContentType(contenttypes.ApplicationJSON).Handler)
		r.Use(middlewares.AddSHA256HashHeader(hg, a.config.c.Key).Handler)

		formatter := metricsformatter.MetricsJSONFormatter{}
		identifier := metricsidentifier.NewJSONIdentifier()
		r.Post("/", fetchmetrics.New(a.storage.s, formatter, identifier, mc).Handler())
	})

	r.Route("/", func(r chi.Router) {
		r.Use(apilogger.New().Handler)
		r.Use(middlewares.CompressGzip().Handler)
		r.Use(middlewares.SetContentType(contenttypes.TextHTML).Handler)
		r.Use(middlewares.AddSHA256HashHeader(hg, a.config.c.Key).Handler)

		r.Get("/", fetchallmetrics.New(a.storage.s).Handler())
	})

	tlsConfig, err := tlsconfig.Provider{CryptoKeyPath: a.config.c.CryptoKey}.Fetch()
	if err != nil {
		return fmt.Errorf("failed to fetch tls config: %w", err)
	}
	srv := &http.Server{Addr: a.config.c.Host, Handler: r, TLSConfig: tlsConfig}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("HTTP server ListenAndServe", zap.Error(err))
		}
	}()

	logger.Log.Info("listening", zap.String("host", a.config.c.Host))

	var s *grpc.Server
	if a.config.c.GRPCHost != "" {
		listen, err := net.Listen("tcp", a.config.c.GRPCHost)
		if err != nil {
			log.Fatal(err)
		}
		s = grpc.NewServer()

		pb.RegisterMetricsServer(s, &igrpc.MetricsServer{})

		logger.Log.Info("gRPC server starts", zap.String("host", a.config.c.GRPCHost))
		if err := s.Serve(listen); err != nil {
			logger.Log.Error("failed to start gRPC server", zap.Error(err))
		}
	}

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
		if s != nil {
			s.GracefulStop()
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
