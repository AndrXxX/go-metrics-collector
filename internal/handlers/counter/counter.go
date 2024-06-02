package counter

import (
	"github.com/AndrXxX/go-metrics-collector/internal/repositories"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func Handler(s repositories.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		params := mux.Vars(r)
		metric := params["metric"]
		if metric == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		value := params["value"]
		if converted, err := strconv.ParseInt(value, 10, 64); err == nil {
			s.Counter(metric, converted)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}
