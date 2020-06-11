package update

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/server"
	"golang.org/x/sync/errgroup"
)

// Arguments represents arguments for person update
type Arguments struct {
	ID              int64                 `json:"id"`
	Description     *string               `json:"description"`
	Priority        *args.VacacnyPriority `json:"priority"`
	Pause           *bool                 `json:"pause"`
	DepartmentName  *string               `json:"department_name"`
	AppointmentName *string               `json:"appointment_name"`
	Skills          []*args.SkillUpdate   `json:"skills"`

	departmentID  *int64
	appointmentID *int64
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	if a.Description == nil && a.Priority == nil && a.Pause == nil &&
		a.DepartmentName == nil && a.AppointmentName == nil && len(a.Skills) == 0 {
		return server.EmptyArgumentsHTTPError
	}

	if a.Description != nil && *a.Description == "" {
		return server.NewHTTPError(http.StatusBadRequest, "description can not be empty")
	}

	if a.DepartmentName != nil && *a.DepartmentName == "" {
		return server.NewHTTPError(http.StatusBadRequest, "department_name can not be empty")
	}

	if a.DepartmentName != nil && *a.AppointmentName == "" {
		return server.NewHTTPError(http.StatusBadRequest, "appointment_name can not be empty")
	}

	for _, s := range a.Skills {
		if err := s.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Resolve loads additional data
func (a *Arguments) Resolve(s *db.Session, r *http.Request) error {
	g, _ := errgroup.WithContext(r.Context())
	if a.DepartmentName != nil {
		g.Go(func() error {
			department, err := database.ReadDepartment(s, []string{*a.DepartmentName})
			if err != nil {
				return err
			}
			if len(department) == 0 {
				return server.NewHTTPError(http.StatusBadRequest, "department_name not found")
			}
			a.departmentID = &department[0].ID

			return nil
		})
	}

	if a.AppointmentName != nil {
		g.Go(func() error {
			appointment, err := database.ReadAppointment(s, []string{*a.AppointmentName})
			if err != nil {
				return err
			}
			if len(appointment) == 0 {
				return server.NewHTTPError(http.StatusBadRequest, "appointment_name not found")
			}
			a.appointmentID = &appointment[0].ID

			return nil
		})
	}

	return g.Wait()
}

// IsVacancyUpdated returns true if vacancy needs to be updated
func (a *Arguments) IsVacancyUpdated() bool {
	return a.Description != nil || a.Priority != nil || a.Pause != nil ||
		a.DepartmentName != nil || a.AppointmentName != nil
}

// IsSkillUpdated return true if some skills has to be updated
func (a *Arguments) IsSkillUpdated() bool {
	return len(a.Skills) != 0
}

// NewArguments creates new CreatePersonArguments instanse
func NewArguments() server.Arguments {
	return &Arguments{}
}
