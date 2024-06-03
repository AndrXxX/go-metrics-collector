package server

import (
	"fmt"
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
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
	rtr.HandleFunc(fmt.Sprintf("/update/counter/{%v}/{%v}", vars.METRIC, vars.VALUE), counter.Handler(&storage))
	rtr.HandleFunc(fmt.Sprintf("/update/gauge/{%v}/{%v}", vars.METRIC, vars.VALUE), gauge.Handler(&storage))
	rtr.HandleFunc(fmt.Sprintf("/update/{unknownType}/{%v}/{%v}", vars.METRIC, vars.VALUE), handlers.BadRequest)
	rtr.HandleFunc("/", http.NotFound)
	muxServe.Handle("/", rtr)
	return http.ListenAndServe(`:8080`, muxServe)
}
