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
	"github.com/AndrXxX/go-metrics-collector/internal/server/services/storageprovider"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run(c *config.Config) error {
	modelCounterStorage := memory.New[*models.Metrics]()
	modelGaugeStorage := memory.New[*models.Metrics]()
	sp := storageprovider.New[interfaces.MetricsStorage]()
	sp.RegisterStorage(metrics.Counter, &modelCounterStorage)
	sp.RegisterStorage(metrics.Gauge, &modelGaugeStorage)
	cFactory := conveyor.Factory(logger.New())
	r := chi.NewRouter()

	r.Route("/update", func(r chi.Router) {
		r.Post(fmt.Sprintf("/{%v}/{%v}/{%v}", vars.MetricType, vars.Metric, vars.Value), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.HasMetricOr404(),
			updatemetrics.New(metricsupdater.New(sp), metricsidentifier.NewURLIdentifier()),
		}).Handler())
	})

	r.Route("/value", func(r chi.Router) {
		r.Get(fmt.Sprintf("/{%v}/{%v}", vars.MetricType, vars.Metric), cFactory.From([]interfaces.Handler{
			middlewares.SetContentType(contenttypes.TextPlain),
			middlewares.HasMetricOr404(),
			fetchmetrics.New(sp, metricstringifier.MetricsValueStringifier{}, metricsidentifier.NewURLIdentifier()),
		}).Handler())
	})

	r.Get("/", cFactory.From([]interfaces.Handler{
		middlewares.SetContentType(contenttypes.TextHTML),
		fetchallmetrics.New(&modelGaugeStorage, &modelCounterStorage),
	}).Handler())

	return http.ListenAndServe(c.Host, r)
}
