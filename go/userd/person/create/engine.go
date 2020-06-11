package create

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
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
		person := &database.PersonModel{}

		if arg.Person != nil {
			person = &database.PersonModel{
				FirstName: arg.Person.FirstName,
				LastName:  arg.Person.LastName,
			}
			if err := database.CreatePerson(tx, person); err != nil {
				return errors.Wrap(err, "can not create base person")
			}
		} else {
			person.ID = *arg.Employee.PersonID
		}

		if arg.Employee != nil {
			employee := &database.EmployeeModel{
				PersonModel: *person,
				Email:       arg.Employee.Email,
				HiredAt:     arg.Employee.HiredAt,
			}
			if err := database.CreatedEmployee(tx, employee); err != nil {
				return errors.Wrap(err, "can not create employee")
			}

			if arg.Employee.Work != nil {
				wh := &database.WorkHistory{
					EmployeeID:    employee.ID,
					DepartmentID:  arg.Employee.Work.departmentID,
					AppointmentID: arg.Employee.Work.appointmentID,
					DateFrom:      arg.Employee.HiredAt,
				}

				if err := database.CreateWorkHistory(tx, wh); err != nil {
					return errors.Wrap(err, "can not create work history unit")
				}
			}
		}

		resultID = person.ID
		return nil
	})
	if err != nil {
		logger.Errorf("can not create person: %s", err)
		return nil, err
	}

	return map[string]int64{"id": resultID}, nil
}
