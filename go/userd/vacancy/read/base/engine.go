package base

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/pkg/errors"
)

// Process reads person record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	b := newReplyBuilder()
	vacancies, err := database.ReadVacancy(sess, arg.ID, arg.DepartmentNames, arg.AppointmentNames)
	if err != nil {
		logger.Errorf("can not load vacancies: %s", err)
		return nil, errors.Wrap(err, "can not load vacancies")
	}
	b.consumeVacancies(vacancies)

	return b.reply(), nil
}
