package common

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	requirement "github.com/GDVFox/tenjin/userd/requirement_check"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/gocraft/dbr/v2"
)

// Process updates task record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	err := db.InTransaction(sess, func(tx *dbr.Tx) error {
		if arg.IsInterviewUpdated() || arg.IsChecksUpdated() {
			interviewUpdate := &database.InterviewUpdateParams{
				ID:            arg.ID,
				PlannedDate:   arg.PlannedTime,
				Comment:       arg.Comment,
				TotalScore:    arg.TotalScore,
				InterviewerID: arg.InterviewerID,
			}
			if err := database.UpdateInterview(tx, interviewUpdate); err != nil {
				return err
			}
		}

		if arg.IsChecksUpdated() {
			if err := requirement.UpdateRequirementChecks(tx, arg.ID, arg.vacancyID, arg.Checks); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Errorf("can not update interview: %s", err)
		return nil, err
	}

	return map[string]int64{"id": arg.ID}, nil
}
