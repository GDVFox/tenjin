package delete

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

// Process deletes task record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	err := db.InTransaction(sess, func(tx *dbr.Tx) error {
		if err := database.DeleteTask(tx, arg.ID); err != nil {
			return errors.Wrap(err, "can not delete task")
		}

		if err := database.DeleteSolution(tx, 0, arg.ID); err != nil && err != db.ErrNoRowsAffected {
			return errors.Wrap(err, "can not delete solutions")
		}

		return nil
	})
	if err != nil {
		logger.Errorf("can not delete task: %s", err)
		return nil, err
	}

	return map[string]int64{"id": arg.ID}, nil
}
