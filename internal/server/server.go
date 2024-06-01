package server

import (
	"github.com/AndrXxX/go-metrics-collector/internal/handlers/counter"
	"github.com/AndrXxX/go-metrics-collector/internal/handlers/gauge"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories/memStorage"
	"github.com/gorilla/mux"
	"net/http"
)

func Run() error {
	storage := memStorage.New()
	muxServe := http.NewServeMux()
	rtr := mux.NewRouter()
	rtr.HandleFunc("/update/counter/{metric}/{value}", counter.Handler(&storage))
	rtr.HandleFunc("/update/gauge/{metric}/{value}", gauge.Handler(&storage))
	rtr.HandleFunc("/*", badRequest)
	muxServe.Handle("/", rtr)
	return http.ListenAndServe(`:8080`, muxServe)
}

func badRequest(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
