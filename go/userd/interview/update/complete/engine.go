package complete

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	requirement "github.com/GDVFox/tenjin/userd/requirement_check"
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
		interviewComplete := &database.InterviewCompleteParams{
			ID:            arg.ID,
			Comment:       arg.Comment,
			TotalScore:    arg.TotalScore,
			InterviewerID: arg.InterviewerID,
		}
		if err := database.CompleteInterview(tx, interviewComplete); err != nil {
			return errors.Wrap(err, "can not complete interview")
		}

		if len(arg.Checks) != 0 {
			if err := requirement.CreateRequirementChecks(tx, arg.ID, arg.vacancyID, arg.Checks); err != nil {
				return errors.Wrap(err, "can not create checks")
			}
		}

		return nil
	})
	if err != nil {
		logger.Errorf("can not complete interview: %s", err)
		return nil, err
	}

	return map[string]int64{"id": arg.ID}, nil
}
