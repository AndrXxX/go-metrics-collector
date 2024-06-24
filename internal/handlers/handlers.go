package handlers

import "net/http"

func BadRequest(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
