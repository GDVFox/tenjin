package update

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/attachment"
	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/gocraft/dbr/v2"
)

// Process updates task record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	err := db.InTransaction(sess, func(tx *dbr.Tx) error {
		if arg.IsPostUpdated() || arg.IsSolutionUpdated() {
			postUpdate := &database.PostUpdateParams{
				ID:   arg.ID,
				Text: arg.Text,
			}
			if err := database.UpdatePost(tx, postUpdate); err != nil {
				return err
			}
		}

		if arg.IsSolutionUpdated() {
			solutionUpdate := &database.SolutionUpdateParams{
				ID:         arg.ID,
				IsApproved: arg.IsApproved,
			}
			if err := database.UpdateSolution(tx, solutionUpdate); err != nil {
				return err
			}
		}

		if arg.IsAttachmentsUpdated() {
			if err := attachment.UpdateAttachments(tx, arg.ID, 0, arg.Attachments); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Errorf("can not update task: %s", err)
		return nil, err
	}

	return map[string]int64{"id": arg.ID}, nil
}
