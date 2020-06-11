package server

import (
	"net/http"

	"github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
)

// RequestProcessor standart http handler
type RequestProcessor func(r *http.Request, sess *database.Session, args Arguments, log *logging.Logger) (interface{}, error)

func (p RequestProcessor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(loggerKey).(*logging.Logger)
	dbsess := database.NewSession(logger)
	a, ok := r.Context().Value(argsKey).(Arguments)
	if !ok {
		a = nil
	}

	if resolver, ok := a.(Resolver); ok {
		if err := resolver.Resolve(dbsess, r); err != nil {
			logger.Errorf("can not resolve args: %s", err)
			if err := WriteError(w, err); err != nil {
				logger.Errorf("can not send error: %s", err)
			}
			return
		}
	}

	reply, err := p(r, dbsess, a, logger)
	if err != nil {
		logger.Errorf("error in processor: %s", err)
		if err := WriteError(w, err); err != nil {
			logger.Errorf("can not send error: %s", err)
		}
		return
	}

	if err := WriteApplicationJSON(w, http.StatusOK, reply); err != nil {
		logger.Errorf("can not send reply: %s", err)
	}
}
