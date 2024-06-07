package server

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/handlers"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories/memstorage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run() error {
	storage := memstorage.New()
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/counter/{%v}/{%v}", vars.METRIC, vars.VALUE), handlers.CounterUpdater(&storage))
		r.Post(fmt.Sprintf("/gauge/{%v}/{%v}", vars.METRIC, vars.VALUE), handlers.GaugeUpdater(&storage))
		r.Post(fmt.Sprintf("/{unknownType}/{%v}/{%v}", vars.METRIC, vars.VALUE), handlers.BadRequest)
	})
	r.Route("/value", func(r chi.Router) {
		r.Get(fmt.Sprintf("/counter/{%v}", vars.METRIC), handlers.CounterFetcher(&storage))
		r.Get(fmt.Sprintf("/gauge/{%v}", vars.METRIC), handlers.GaugeFetcher(&storage))
		r.Get(fmt.Sprintf("/{unknownType}/{%v}/{%v}", vars.METRIC, vars.VALUE), handlers.BadRequest)
	})
	r.Get("/", handlers.MetricsFetcher(&storage))

	return http.ListenAndServe(":8080", r)
}
