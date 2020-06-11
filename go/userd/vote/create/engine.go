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

	vote := &database.VoteModel{
		EmployeeID: arg.EmployeeID,
		PostID:     db.NewNullInt64(arg.PostID),
		CommentID:  db.NewNullInt64(arg.CommentID),
		Delta:      arg.Delta,
	}
	if err := database.CreateVote(sess, vote); err != nil {
		logger.Errorf("can not create vote: %s", err)
		return nil, errors.Wrap(err, "can not create vote")
	}

	return map[string]string{"result": "OK"}, nil
}
