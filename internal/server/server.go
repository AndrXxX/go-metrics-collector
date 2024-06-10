package server

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/handlers"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories/memstorage"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run(c *config.Config) error {
	storage := memstorage.New()
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/counter/{%v}/{%v}", vars.Metric, vars.Value), handlers.CounterUpdater(&storage))
		r.Post(fmt.Sprintf("/gauge/{%v}/{%v}", vars.Metric, vars.Value), handlers.GaugeUpdater(&storage))
		r.Post(fmt.Sprintf("/{unknownType}/{%v}/{%v}", vars.Metric, vars.Value), handlers.BadRequest)
	})
	r.Route("/value", func(r chi.Router) {
		r.Get(fmt.Sprintf("/counter/{%v}", vars.Metric), handlers.CounterFetcher(&storage))
		r.Get(fmt.Sprintf("/gauge/{%v}", vars.Metric), handlers.GaugeFetcher(&storage))
		r.Get(fmt.Sprintf("/{unknownType}/{%v}/{%v}", vars.Metric, vars.Value), handlers.BadRequest)
	})
	r.Get("/", handlers.MetricsFetcher(&storage))

	return http.ListenAndServe(c.Host, r)
}
