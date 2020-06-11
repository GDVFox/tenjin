package update

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	"github.com/GDVFox/tenjin/userd/skill"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/gocraft/dbr/v2"
)

// Process reads person record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	err := db.InTransaction(sess, func(tx *dbr.Tx) error {
		if arg.IsVacancyUpdated() || arg.IsSkillUpdated() {
			params := &database.VacancyUpdateParams{
				ID:            arg.ID,
				Description:   arg.Description,
				Priority:      arg.Priority.GetValue(),
				Pause:         arg.Pause,
				DepartmentID:  arg.departmentID,
				AppointmentID: arg.appointmentID,
			}
			if err := database.UpdateVacancy(sess, params); err != nil {
				return err
			}
		}

		if arg.IsSkillUpdated() {
			if err := skill.UpdateSkills(tx, skill.VacancyRequirementType, arg.ID, arg.Skills); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Errorf("can not update vacancy: %s", err)
		return nil, err
	}

	return map[string]int64{"id": arg.ID}, nil
}
