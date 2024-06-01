package main

import (
	"github.com/AndrXxX/go-metrics-collector/internal/repositories"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories/memStorage"
	"github.com/gorilla/mux"
	"net/http"
	"slices"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	storage := memStorage.New()
	muxServe := http.NewServeMux()
	rtr := mux.NewRouter()
	rtr.HandleFunc("/update/{type}/{metric}/{value}", handler(&storage))
	rtr.HandleFunc("/*", badRequest)
	muxServe.Handle("/", rtr)
	return http.ListenAndServe(`:8080`, muxServe)
}

func badRequest(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func handler(s repositories.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		params := mux.Vars(r)
		mType := params["type"]

		if !slices.Contains([]string{"gauge", "counter"}, mType) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		metric := params["metric"]
		if metric == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		value := params["value"]
		if mType == "counter" {
			if converted, err := strconv.ParseInt(value, 10, 64); err == nil {
				s.Counter(metric, converted)
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if converted, err := strconv.ParseFloat(value, 64); err == nil {
			s.Gauge(metric, converted)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}
