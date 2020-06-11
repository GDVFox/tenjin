package server

import (
	"encoding/json"
	"net/http"
)

// WriteApplicationJSON write answer in json
func WriteApplicationJSON(w http.ResponseWriter, code int, v interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, err := json.Marshal(v)
	if err != nil {
		w.Header().Set("Content-Type", "plain/text")

		code = http.StatusInternalServerError
		resp = []byte("internal server error")
	}

	w.WriteHeader(code)
	_, wErr := w.Write(resp)
	// err has more priority
	if err == nil {
		err = wErr
	}
	return
}
