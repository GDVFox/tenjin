package delete

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
)

// Process deletes task record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	if err := database.DeleteComment(sess, arg.ID); err != nil {
		logger.Errorf("can not delete comment: %s", err)
		return nil, err
	}

	return map[string]int64{"id": arg.ID}, nil
}
