package create

import (
	"net/http"

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
		vacancy := &database.VacancyModel{
			Description:   arg.Description,
			Priority:      arg.Priority.GetValue(),
			DepartmentID:  arg.departmentID,
			AppointmentID: arg.appointmentID,
		}
		if err := database.CreateVacancy(sess, vacancy); err != nil {
			return errors.Wrap(err, "can not create vacancy")
		}

		if len(arg.Skills) != 0 {
			if err := skill.CreateSkills(tx, skill.VacancyRequirementType, vacancy.ID, arg.Skills); err != nil {
				return errors.Wrap(err, "can not create skills")
			}
		}

		resultID = vacancy.ID
		return nil
	})
	if err != nil {
		logger.Errorf("can not create vacancy: %s", err)
		return nil, err
	}

	return map[string]int64{"id": resultID}, nil
}
