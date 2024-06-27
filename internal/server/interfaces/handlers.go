package interfaces

import "net/http"

type HandlerFunc func(http.ResponseWriter, *http.Request) (ok bool)
