package update

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

// Process reads person record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	err := db.InTransaction(sess, func(tx *dbr.Tx) error {
		if arg.IsCommonUpdated() || arg.IsWorkUpdated() {
			params := &database.PersonUpdateParams{
				ID:        arg.ID,
				Blocked:   arg.Blocked,
				FirstName: arg.FirstName,
				LastName:  arg.LastName,
			}
			if err := database.UpdatePerson(tx, params); err != nil {
				return errors.Wrap(err, "can not update person")
			}
		}

		if arg.IsWorkUpdated() {
			if arg.Work.DateFrom != nil {
				wk := &database.WorkHistory{
					EmployeeID:    arg.ID,
					DepartmentID:  arg.Work.departmentID,
					AppointmentID: arg.Work.appointmentID,
					DateFrom:      *arg.Work.DateFrom,
				}
				if err := database.CreateWorkHistory(tx, wk); err != nil {
					return errors.Wrap(err, "can not create work history")
				}
			} else {
				if err := database.StopWork(tx, arg.ID, arg.Work.departmentID,
					arg.Work.appointmentID, *arg.Work.DateTo); err != nil {
					return errors.Wrap(err, "can not delete work history")
				}
			}
		}

		return nil
	})
	if err != nil {
		logger.Errorf("can not update person: %s", err)
		return nil, err
	}

	return map[string]int64{"id": arg.ID}, nil
}
