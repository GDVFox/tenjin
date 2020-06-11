package server

import "net/http"

type responseWriterWithCode struct {
	http.ResponseWriter
	statusCode int
}

func wrapRW(w http.ResponseWriter) *responseWriterWithCode {
	return &responseWriterWithCode{w, http.StatusOK}
}

func (r *responseWriterWithCode) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}
