package update

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

// Process updates task record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	err := db.InTransaction(sess, func(tx *dbr.Tx) error {
		if arg.IsPostUpdated() {
			params := &database.CommentUpdateParams{
				ID:   arg.ID,
				Text: arg.Text,
			}
			if err := database.UpdateComment(tx, params); err != nil {
				return errors.Wrap(err, "can not update comment")
			}
		}

		if arg.IsAttachmentsUpdated() {
			if err := attachment.UpdateAttachments(tx, 0, arg.ID, arg.Attachments); err != nil {
				return errors.Wrap(err, "can not update attachments")
			}
		}

		return nil
	})
	if err != nil {
		logger.Errorf("can not update comment: %s", err)
		return nil, err
	}

	return map[string]int64{"id": arg.ID}, nil
}
