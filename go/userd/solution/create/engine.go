package create

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/attachment"
	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

// Process creates solution record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	var resultID int64
	err := db.InTransaction(sess, func(tx *dbr.Tx) error {
		post := database.PostModel{
			Text:       arg.Text,
			EmployeeID: arg.AuthorID,
		}

		if err := database.CreatePost(tx, &post); err != nil {
			return errors.Wrap(err, "can not create base post")
		}

		solution := database.SolutionModel{
			PostModel: post,
			TaskID:    arg.TaskID,
		}
		if err := database.CreateSolution(tx, &solution); err != nil {
			return errors.Wrap(err, "can not create solution")
		}

		if len(arg.Attachments) != 0 {
			if err := attachment.CreateAttachments(tx, post.ID, 0, arg.Attachments); err != nil {
				return errors.Wrap(err, "can not create attachments")
			}
		}

		resultID = post.ID
		return nil
	})
	if err != nil {
		logger.Errorf("can not create solution: %s", err)
		return nil, err
	}

	return map[string]int64{"id": resultID}, nil
}
