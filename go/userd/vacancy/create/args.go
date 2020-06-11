package create

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/server"
	"golang.org/x/sync/errgroup"
)

// Arguments represents arguments for person creation
type Arguments struct {
	Description     string                `json:"description"`
	Priority        *args.VacacnyPriority `json:"priority"`
	DepartmentName  string                `json:"department_name"`
	AppointmentName string                `json:"appointment_name"`
	Skills          []*args.SkillCreate   `json:"skills"`

	departmentID  int64
	appointmentID int64
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	if a.Description == "" {
		return server.NewHTTPError(http.StatusBadRequest, "description can not be empty")
	}

	if a.DepartmentName == "" {
		return server.NewHTTPError(http.StatusBadRequest, "department_name can not be empty")
	}

	if a.AppointmentName == "" {
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
	g.Go(func() error {
		department, err := database.ReadDepartment(s, []string{a.DepartmentName})
		if err != nil {
			return err
		}
		if len(department) == 0 {
			return server.NewHTTPError(http.StatusBadRequest, "Wrong department name")
		}
		a.departmentID = department[0].ID

		return nil
	})

	g.Go(func() error {
		appointment, err := database.ReadAppointment(s, []string{a.AppointmentName})
		if err != nil {
			return err
		}
		if len(appointment) == 0 {
			return server.NewHTTPError(http.StatusBadRequest, "Wrong appointment name")
		}
		a.appointmentID = appointment[0].ID

		return nil
	})

	return g.Wait()
}

// NewArguments creates new CreatePersonArguments instanse
func NewArguments() server.Arguments {
	return &Arguments{}
}
