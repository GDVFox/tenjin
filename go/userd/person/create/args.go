package create

import (
	"net/http"
	"time"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/server"
	"golang.org/x/sync/errgroup"
)

// PersonArguments represents person create args
type PersonArguments struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Validate checks argument correct
func (a *PersonArguments) Validate() error {
	return nil
}

// WorkArguments represents work create args
type WorkArguments struct {
	DepartmentName  string `json:"department_name"`
	AppointmentName string `json:"appointment_name"`

	departmentID  int64
	appointmentID int64
}

// Validate checks argument correct
func (a *WorkArguments) Validate() error {
	if a.DepartmentName == "" || a.AppointmentName == "" {
		return server.EmptyArgumentsHTTPError
	}

	return nil
}

// Resolve loads additional data
func (a *WorkArguments) Resolve(s *db.Session, r *http.Request) error {
	g, _ := errgroup.WithContext(r.Context())
	g.Go(func() error {
		department, err := database.ReadDepartment(s, []string{a.DepartmentName})
		if err != nil {
			return err
		}
		if len(department) == 0 {
			return server.NewHTTPError(http.StatusBadRequest, "department_name unknown")
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
			return server.NewHTTPError(http.StatusBadRequest, "appointment_name unknown")
		}
		a.appointmentID = appointment[0].ID

		return nil
	})

	return g.Wait()
}

// EmployeeArguments represents employee create args
type EmployeeArguments struct {
	PersonID *int64         `json:"person_id"`
	Email    string         `json:"email"`
	HiredAt  time.Time      `json:"hired_at"`
	Work     *WorkArguments `json:"work"`
}

// Validate checks argument correct
func (a *EmployeeArguments) Validate() error {
	if a.Work != nil {
		return a.Work.Validate()
	}

	return nil
}

// Resolve loads additional data
func (a *EmployeeArguments) Resolve(s *db.Session, r *http.Request) error {
	if a.Work != nil {
		return a.Work.Resolve(s, r)
	}

	return nil
}

// Arguments represents arguments for person creation
type Arguments struct {
	Person   *PersonArguments   `json:"person"`
	Employee *EmployeeArguments `json:"employee"`
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	if a.Person == nil && a.Employee == nil {
		return server.EmptyArgumentsHTTPError
	}

	if a.Person != nil && a.Employee.PersonID != nil {
		return server.NewHTTPError(http.StatusBadRequest, "person_id and person can not be specified together")
	}

	if a.Person != nil {
		if err := a.Person.Validate(); err != nil {
			return err
		}
	}

	if a.Employee != nil {
		if err := a.Employee.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Resolve loads additional data
func (a *Arguments) Resolve(s *db.Session, r *http.Request) error {
	if a.Employee != nil {
		return a.Employee.Resolve(s, r)
	}

	return nil
}

// NewArguments creates new args instanse
func NewArguments() server.Arguments {
	return &Arguments{}
}
