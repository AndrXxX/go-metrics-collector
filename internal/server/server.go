package server

import (
	"github.com/AndrXxX/go-metrics-collector/internal/handlers"
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
	rtr.HandleFunc("/*", handlers.BadRequest)
	muxServe.Handle("/", rtr)
	return http.ListenAndServe(`:8080`, muxServe)
}
