package server

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/fetchallmetrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/fetchmetrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/logger"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/middlewares"
	"github.com/AndrXxX/go-metrics-collector/internal/server/api/updatemetrics"
	"github.com/AndrXxX/go-metrics-collector/internal/server/config"
	"github.com/AndrXxX/go-metrics-collector/internal/server/interfaces"
	"github.com/AndrXxX/go-metrics-collector/internal/server/models"
	"github.com/AndrXxX/go-metrics-collector/internal/server/repositories/memory"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/conveyor"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricstringifier"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsupdater"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run(c *config.Config) error {
	modelCounterStorage := memory.New[*models.Metrics]()
	modelGaugeStorage := memory.New[*models.Metrics]()
	cFactory := conveyor.Factory(logger.New())
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/counter/{%v}/{%v}", vars.Metric, vars.Value), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType("text/plain"),
			middlewares.HasMetricOr404(),
			updatemetrics.New(metricsupdater.NewCounterUpdater(&modelCounterStorage)),
		}).Handler())
		r.Post(fmt.Sprintf("/gauge/{%v}/{%v}", vars.Metric, vars.Value), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType("text/plain"),
			middlewares.HasMetricOr404(),
			updatemetrics.New(metricsupdater.NewGaugeUpdater(&modelGaugeStorage)),
		}).Handler())
		r.Post(fmt.Sprintf("/{unknownType}/{%v}/{%v}", vars.Metric, vars.Value), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType("text/plain"),
			middlewares.SetHTTPError(http.StatusBadRequest),
		}).Handler())
	})
	r.Route("/value", func(r chi.Router) {
		r.Get(fmt.Sprintf("/counter/{%v}", vars.Metric), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType("text/plain"),
			middlewares.HasMetricOr404(),
			fetchmetrics.New(&modelCounterStorage, metricstringifier.MetricsValueStringifier{}),
		}).Handler())
		r.Get(fmt.Sprintf("/gauge/{%v}", vars.Metric), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType("text/plain"),
			middlewares.HasMetricOr404(),
			fetchmetrics.New(&modelGaugeStorage, metricstringifier.MetricsValueStringifier{}),
		}).Handler())
		r.Get(fmt.Sprintf("/{unknownType}/{%v}/{%v}", vars.Metric, vars.Value), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType("text/plain"),
			middlewares.SetHTTPError(http.StatusBadRequest),
		}).Handler())
	})
	r.Get("/", cFactory.From([]interfaces.Handler{
		middlewares.SetContentType("text/html; charset=utf-8"),
		fetchallmetrics.New(&modelGaugeStorage, &modelCounterStorage),
	}).Handler())

	return http.ListenAndServe(c.Host, r)
}
