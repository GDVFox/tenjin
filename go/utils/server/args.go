package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/GDVFox/tenjin/utils/database"
	"github.com/go-playground/form"
	"github.com/gorilla/mux"
)

var (
	decoder *form.Decoder
)

func init() {
	decoder = form.NewDecoder()
	decoder.SetTagName("json")
}

// Parse creates struct from arguments
func Parse(r *http.Request, v interface{}) error {
	defer r.Body.Close()

	switch strings.Split(r.Header.Get("Content-Type"), ";")[0] {
	case "application/json":
		return parseJSON(r, v)
	default:
		return parseForm(r, v)
	}
}

func parseJSON(r *http.Request, to interface{}) error {
	if err := parseForm(r, to); err != nil {
		return err
	}
	if err := json.NewDecoder(r.Body).Decode(to); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func parseForm(r *http.Request, to interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	for k, v := range mux.Vars(r) {
		r.Form.Add(k, v)
	}

	if err := decoder.Decode(to, r.Form); err != nil {
		return err
	}

	return nil
}

// Arguments interface for arguments
type Arguments interface {
	Validate() error
}

// ArgumentFactory creates new argument instance
type ArgumentFactory func() Arguments

// Resolver interface for structures that needs database for loading additional data
type Resolver interface {
	Resolve(s *database.Session, r *http.Request) error
}
