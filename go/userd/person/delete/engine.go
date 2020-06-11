package delete

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/pkg/errors"
)

// Process deletes person record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	if err := database.DeletePerson(sess, arg.ID); err != nil {
		logger.Errorf("can not delete person: %s", err)
		return nil, errors.Wrap(err, "can not delete person")
	}

	return map[string]int64{"id": arg.ID}, nil
}
