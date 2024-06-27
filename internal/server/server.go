package server

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/fetchcounter"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/fetchgauge"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/update_gauge"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/handlers"
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/middlewares"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/conveyor"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/counter"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run(c *config.Config) error {
	gaugeStorage := memory.New[float64]()
	counterStorage := memory.New[int64]()
	counterUpdater := counter.New(&counterStorage)
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/counter/{%v}/{%v}", vars.Metric, vars.Value), conveyor.New(logger.RequestLogger).From([]interfaces.Handler{
			middlewares.HasMetricOr404(),
			updatecounter.New(counterUpdater),
		}).Handler())
		r.Post(fmt.Sprintf("/gauge/{%v}/{%v}", vars.Metric, vars.Value), conveyor.New(logger.RequestLogger).From([]interfaces.Handler{
			middlewares.HasMetricOr404(),
			update_gauge.New(&gaugeStorage),
		}).Handler())
		r.Post(fmt.Sprintf("/{unknownType}/{%v}/{%v}", vars.Metric, vars.Value), logger.RequestLogger(handlers.BadRequest()))
	})
	r.Route("/value", func(r chi.Router) {
		r.Get(fmt.Sprintf("/counter/{%v}", vars.Metric), conveyor.New(logger.RequestLogger).From([]interfaces.Handler{
			middlewares.HasMetricOr404(),
			fetchcounter.New(&counterStorage),
		}).Handler())
		r.Get(fmt.Sprintf("/gauge/{%v}", vars.Metric), conveyor.New(logger.RequestLogger).From([]interfaces.Handler{
			middlewares.HasMetricOr404(),
			fetchgauge.New(&gaugeStorage),
		}).Handler())
		r.Get(fmt.Sprintf("/{unknownType}/{%v}/{%v}", vars.Metric, vars.Value), handlers.BadRequest())
	})
	r.Get("/", logger.RequestLogger(handlers.MetricsFetcher(&gaugeStorage, &counterStorage)))

	return http.ListenAndServe(c.Host, r)
}
