package counter

import (
	"github.com/AndrXxX/go-metrics-collector/internal/enums/vars"
	"github.com/AndrXxX/go-metrics-collector/internal/repositories"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func Handler(s repositories.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		params := mux.Vars(r)
		metric := params[vars.METRIC]
		if metric == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		value := params[vars.VALUE]
		if converted, err := strconv.ParseInt(value, 10, 64); err == nil {
			s.Counter(metric, converted)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}
