package create

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/pkg/errors"
)

// Process creates appointment record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	if err := database.CreateAppointment(sess, arg.AppointmentName); err != nil {
		logger.Errorf("can not create appointment: %s", err)
		return nil, errors.Wrap(err, "can not create appointment")
	}

	return map[string]string{"result": "OK"}, nil
}
