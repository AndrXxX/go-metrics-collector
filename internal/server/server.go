package server

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/handlers"
	"github.com/AndrXxX/go-metrics-collector/internal/handlers/counter"
	"github.com/AndrXxX/go-metrics-collector/internal/handlers/gauge"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories/memstorage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run() error {
	storage := memstorage.New()
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/counter/{%v}/{%v}", vars.METRIC, vars.VALUE), counter.UpdateHandler(&storage))
		r.Post(fmt.Sprintf("/gauge/{%v}/{%v}", vars.METRIC, vars.VALUE), gauge.UpdateHandler(&storage))
		r.Post("/{unknownType}", handlers.BadRequest)
	})
	r.Route("/value", func(r chi.Router) {
		r.Get(fmt.Sprintf("/counter/{%v}", vars.METRIC), counter.FetchHandler(&storage))
		r.Get("/{unknownType}", handlers.BadRequest)
	})
	r.Handle("/", http.NotFoundHandler())

	return http.ListenAndServe(":8080", r)
}
