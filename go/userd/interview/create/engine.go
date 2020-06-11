package create

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/pkg/errors"
)

// Process creates person record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	model := &database.InterviewModel{
		PlanedDate: arg.PlannedDate,
		PersonID:   arg.PersonID,
		VacancyID:  arg.VacancyID,
	}
	if err := database.CreateInterview(sess, model); err != nil {
		logger.Errorf("can not create interview: %s", err)
		return nil, errors.Wrap(err, "can not create interview")
	}

	return map[string]int64{"id": model.ID}, nil
}
