package sha256

import (
	"bytes"
	"net/http"
)

type RequestWriter struct {
	HG             hashGenerator
	OriginalWriter http.ResponseWriter
	Buffer         *bytes.Buffer
	Key            string
}

// Header is implementation of http.ResponseWriter
func (w *RequestWriter) Header() http.Header {
	return w.OriginalWriter.Header()
}

// Write is implementation of http.ResponseWriter
func (w *RequestWriter) Write(data []byte) (int, error) {
	w.Buffer.Write(data)
	return w.OriginalWriter.Write(data)
}

// WriteHeader is implementation of http.ResponseWriter
func (w *RequestWriter) WriteHeader(statusCode int) {
	w.Header().Add("HashSHA256", w.HG.Generate(w.Key, w.Buffer.Bytes()))
	w.OriginalWriter.WriteHeader(statusCode)
}
