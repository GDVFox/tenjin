package update

import (
	"net/http"
	"time"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/server"
	"golang.org/x/sync/errgroup"
)

// WorkArguments represents arguments for person work update
type WorkArguments struct {
	DepartmentName  string     `json:"department_name"`
	AppointmentName string     `json:"appointment_name"`
	DateFrom        *time.Time `json:"date_from"`
	DateTo          *time.Time `json:"date_to"`

	departmentID  int64
	appointmentID int64
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

// Validate checks argument correct
func (a *WorkArguments) Validate() error {
	if a.DepartmentName == "" || a.AppointmentName == "" {
		return server.NotFoundHTTPError
	}
	if a.DateFrom == nil {
		tmp := time.Now()
		a.DateTo = &tmp
	}

	// должно быть задано либо DateFrom, либо DateTo
	if a.DateFrom != nil && a.DateTo != nil {
		return server.NewHTTPError(http.StatusBadRequest, "only one in (date_from, date_to) may be specified")
	}

	return nil
}

// Arguments represents arguments for person update
type Arguments struct {
	ID        int64          `json:"id"`
	FirstName *string        `json:"first_name"`
	LastName  *string        `json:"last_name"`
	Blocked   *bool          `json:"blocked"`
	Work      *WorkArguments `json:"work"`
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	if a.FirstName == nil && a.LastName == nil && a.Blocked == nil {
		return server.EmptyArgumentsHTTPError
	}
	if (a.FirstName != nil && *a.FirstName == "") ||
		(a.LastName != nil && *a.LastName == "") {
		return server.NewHTTPError(http.StatusBadRequest, "first_name and last_name can not be empty")
	}

	if a.Work != nil {
		return a.Work.Validate()
	}
	return nil
}

// Resolve loads additional data
func (a *Arguments) Resolve(s *db.Session, r *http.Request) error {
	if a.Work != nil {
		return a.Work.Resolve(s, r)
	}

	return nil
}

// IsCommonUpdated returns true if common fields has to be updated
func (a *Arguments) IsCommonUpdated() bool {
	return a.FirstName != nil || a.LastName != nil || a.Blocked != nil
}

// IsWorkUpdated returns true if work fields has to be updated
func (a *Arguments) IsWorkUpdated() bool {
	return a.Work != nil
}

// NewArguments creates new CreatePersonArguments instanse
func NewArguments() server.Arguments {
	return &Arguments{}
}
