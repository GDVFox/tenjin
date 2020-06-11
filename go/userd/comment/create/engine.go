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

// Process creates appointment record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	var resultID int64
	err := db.InTransaction(sess, func(tx *dbr.Tx) error {
		comment := &database.CommentModel{
			Text:       arg.Text,
			EmployeeID: arg.AuthorID,
			PostID:     arg.PostID,
			Parent:     db.NewNullInt64(arg.Parent),
		}

		if err := database.CreateComment(tx, comment); err != nil {
			return errors.Wrap(err, "can not create comment")
		}

		if len(arg.Attachments) != 0 {
			if err := attachment.CreateAttachments(tx, 0, comment.ID, arg.Attachments); err != nil {
				return errors.Wrap(err, "can not create attachments")
			}
		}

		resultID = comment.ID
		return nil
	})
	if err != nil {
		logger.Errorf("can not create comment: %s", err)
		return nil, err
	}

	return map[string]int64{"id": resultID}, nil
}
