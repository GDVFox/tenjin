package create

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/attachment"
	"github.com/GDVFox/tenjin/userd/database"
	"github.com/GDVFox/tenjin/userd/skill"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

// Process creates person record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	var resultID int64
	err := db.InTransaction(sess, func(tx *dbr.Tx) error {
		post := database.PostModel{
			Text:       arg.Text,
			EmployeeID: arg.AuthorID,
		}

		if err := database.CreatePost(tx, &post); err != nil {
			return errors.Wrap(err, "can not create post")
		}

		task := database.TaskModel{
			PostModel: post,
			Title:     arg.Title,
		}
		if err := database.CreateTask(tx, &task); err != nil {
			return errors.Wrap(err, "can not create task")
		}

		if len(arg.Skills) != 0 {
			if err := skill.CreateSkills(tx, skill.TaskRequirementType, post.ID, arg.Skills); err != nil {
				return errors.Wrap(err, "can not create skills")
			}
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
		logger.Errorf("can not create task: %s", err)
		return nil, err
	}

	return map[string]int64{"id": resultID}, nil
}
