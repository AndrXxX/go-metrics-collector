package server

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/handlers"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
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
		r.Post(fmt.Sprintf("/counter/{%v}/{%v}", vars.Metric, vars.Value), logger.RequestLogger(handlers.CounterUpdater(&counterUpdater)))
		r.Post(fmt.Sprintf("/gauge/{%v}/{%v}", vars.Metric, vars.Value), logger.RequestLogger(handlers.GaugeUpdater(&gaugeStorage)))
		r.Post(fmt.Sprintf("/{unknownType}/{%v}/{%v}", vars.Metric, vars.Value), logger.RequestLogger(handlers.BadRequest()))
	})
	r.Route("/value", func(r chi.Router) {
		r.Get(fmt.Sprintf("/counter/{%v}", vars.Metric), logger.RequestLogger(handlers.CounterFetcher(&counterStorage)))
		r.Get(fmt.Sprintf("/gauge/{%v}", vars.Metric), logger.RequestLogger(handlers.GaugeFetcher(&gaugeStorage)))
		r.Get(fmt.Sprintf("/{unknownType}/{%v}/{%v}", vars.Metric, vars.Value), handlers.BadRequest())
	})
	r.Get("/", logger.RequestLogger(handlers.MetricsFetcher(&gaugeStorage, &counterStorage)))

	return http.ListenAndServe(c.Host, r)
}
