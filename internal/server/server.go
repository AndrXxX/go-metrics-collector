package server

import (
	"github.com/AndrXxX/go-metrics-collector/internal/handlers"
	"github.com/AndrXxX/go-metrics-collector/internal/handlers/counter"
	"github.com/AndrXxX/go-metrics-collector/internal/handlers/gauge"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories/memstorage"
	"github.com/gorilla/mux"
	"net/http"
)

func Run() error {
	storage := memstorage.New()
	muxServe := http.NewServeMux()
	rtr := mux.NewRouter()
	rtr.HandleFunc("/update/counter/{metric}/{value}", counter.Handler(&storage))
	rtr.HandleFunc("/update/gauge/{metric}/{value}", gauge.Handler(&storage))
	rtr.HandleFunc("/update/{unknownType}/{metric}/{value}", handlers.BadRequest)
	rtr.HandleFunc("/", http.NotFound)
	muxServe.Handle("/", rtr)
	return http.ListenAndServe(`:8080`, muxServe)
}
