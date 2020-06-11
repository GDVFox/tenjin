package read

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/pkg/errors"
)

// Process reads appointments record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	result, err := database.ReadAppointment(sess, nil)
	if err != nil {
		logger.Errorf("can not read appointments: %s", err)
		return nil, errors.Wrap(err, "can not read appointments")
	}
	return result, nil
}
