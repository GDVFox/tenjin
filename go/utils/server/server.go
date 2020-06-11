package server

import (
	"context"
	"net/http"
	"time"

	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Route is description of api
type Route struct {
	Method     string
	Pattern    string
	Handler    http.Handler
	ArgFactory ArgumentFactory
}

// Server represents http server
type Server struct {
	serv   *http.Server
	logger *logging.Logger
}

func newRouter(routes []*Route, logger *logging.Logger) http.Handler {
	r := mux.NewRouter().PathPrefix("/v1").Subrouter()

	for _, route := range routes {
		handler := route.Handler
		if route.ArgFactory != nil {
			handler = argsMiddleware(handler, route.ArgFactory)
		}

		handler = loggingMiddleware(handler, logger)
		handler = panicMiddleware(handler, logger)

		r.Handle(route.Pattern, handler).Methods(route.Method)
	}

	corsMiddleware := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE"},
			AllowedHeaders:   []string{"Content-Type"},
			AllowCredentials: true,
		},
	)

	return corsMiddleware.Handler(r)
}

// NewServer creates new http server
func NewServer(cfg *Config, routes []*Route, logger *logging.Logger) *Server {
	return &Server{
		serv: &http.Server{
			Addr:    cfg.Address(),
			Handler: newRouter(routes, logger),
		},
		logger: logger,
	}
}

// StartWithCancel creates server
func (s *Server) StartWithCancel(cancel <-chan struct{}) error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		errCh <- s.serv.ListenAndServe()
	}()

	s.logger.Info("server was started!")
	select {
	case <-cancel:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		if err := s.serv.Shutdown(ctx); err != nil {
			s.logger.Errorf("can not shutdown server: %s", err)
		}
		cancel()
	case err := <-errCh:
		return err
	}
	s.logger.Info("server was stoped!")
	return nil
}
