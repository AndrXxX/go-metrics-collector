package server

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/contenttypes"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/metrics"
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
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsidentifier"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricstringifier"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsupdater"
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/metricsvaluesetter"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run(c *config.Config) error {
	modelCounterStorage := memory.New[*models.Metrics]()
	modelGaugeStorage := memory.New[*models.Metrics]()
	cFactory := conveyor.Factory(logger.New())
	mvsFactory := metricsvaluesetter.Factory()
	r := chi.NewRouter()

	r.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/counter/{%v}/{%v}", vars.Metric, vars.Value), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.HasMetricOr404(),
			updatemetrics.New(metricsupdater.New(&modelCounterStorage, mvsFactory.CounterValueSetter(), metrics.Counter)),
		}).Handler())

		r.Post(fmt.Sprintf("/gauge/{%v}/{%v}", vars.Metric, vars.Value), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.HasMetricOr404(),
			updatemetrics.New(metricsupdater.New(&modelGaugeStorage, mvsFactory.GaugeValueSetter(), metrics.Gauge)),
		}).Handler())

		r.Post(fmt.Sprintf("/{unknownType}/{%v}/{%v}", vars.Metric, vars.Value), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.SetHTTPError(http.StatusBadRequest),
		}).Handler())
	})

	r.Route("/value", func(r chi.Router) {
		r.Get(fmt.Sprintf("/counter/{%v}", vars.Metric), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.HasMetricOr404(),
			fetchmetrics.New(&modelCounterStorage, metricstringifier.MetricsValueStringifier{}, metricsidentifier.NewURLIdentifier(metrics.Counter)),
		}).Handler())

		r.Get(fmt.Sprintf("/gauge/{%v}", vars.Metric), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.HasMetricOr404(),
			fetchmetrics.New(&modelGaugeStorage, metricstringifier.MetricsValueStringifier{}, metricsidentifier.NewURLIdentifier(metrics.Gauge)),
		}).Handler())

		r.Get(fmt.Sprintf("/{unknownType}/{%v}/{%v}", vars.Metric, vars.Value), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.SetHTTPError(http.StatusBadRequest),
		}).Handler())
	})

	r.Get("/", cFactory.From([]interfaces.Handler{
		middlewares.SetContentType(contenttypes.TextHTML),
		fetchallmetrics.New(&modelGaugeStorage, &modelCounterStorage),
	}).Handler())

	return http.ListenAndServe(c.Host, r)
}
