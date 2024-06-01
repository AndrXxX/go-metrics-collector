package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"slices"
	"strconv"
)

type Storage interface {
	Gauge(metric string, value float64)
	Counter(metric string, value int64)
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (s *MemStorage) Gauge(metric string, value float64) {
	s.gauge[metric] = value
}

func (s *MemStorage) Counter(metric string, value int64) {
	s.counter[metric] += value
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	memStorage := MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
	muxServe := http.NewServeMux()
	rtr := mux.NewRouter()
	rtr.HandleFunc("/update/{type}/{metric}/{value}", handler(&memStorage))
	rtr.HandleFunc("/*", badRequest)
	muxServe.Handle("/", rtr)
	return http.ListenAndServe(`:8080`, muxServe)
}

func badRequest(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func handler(s Storage) func(http.ResponseWriter, *http.Request) {
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
