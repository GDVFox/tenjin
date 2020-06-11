package server

import (
	"context"
	"net/http"

	"github.com/GDVFox/tenjin/utils/logging"
)

type middlwareDataKey int

const (
	loggerKey middlwareDataKey = iota
	argsKey
)

func panicMiddleware(next http.Handler, logger *logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("server: panic in %s: %s", r.URL.Path, err)
				WriteError(w, UnknownHTTPError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler, logger *logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loggerKey, logger)
		logger.Infof("server: %s %s", r.Method, r.URL.Path)

		rw := wrapRW(w)
		next.ServeHTTP(rw, r.WithContext(ctx))

		logger.Infof("server: done %d", rw.statusCode)
	})
}

func argsMiddleware(next http.Handler, factory ArgumentFactory) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(loggerKey).(*logging.Logger)

		args := factory()
		if err := Parse(r, args); err != nil {
			logger.Errorf("server: can not parse args %s", err)
			WriteError(w, NewHTTPError(http.StatusBadRequest, err.Error()))
			return
		}

		if err := args.Validate(); err != nil {
			logger.Errorf("server: can validate args %s", err)
			if herr, ok := err.(*HTTPError); ok {
				WriteError(w, herr)
			} else {
				WriteError(w, BadRequestHTTPError)
			}
			return
		}

		ctx := context.WithValue(r.Context(), argsKey, args)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
