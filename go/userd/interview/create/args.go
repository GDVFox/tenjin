package create

import (
	"net/http"
	"time"

	"github.com/GDVFox/tenjin/utils/server"
)

// Arguments represents arguments for interview creation
type Arguments struct {
	VacancyID   int64     `json:"vacancy_id"`
	PersonID    int64     `json:"person_id"`
	PlannedDate time.Time `json:"planned_date"`
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	if a.VacancyID == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "vacancy_id can not be 0")
	}

	if a.PersonID == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "person_id can not be 0")
	}

	if a.PlannedDate.IsZero() {
		return server.NewHTTPError(http.StatusBadRequest, "planned_date can not be zero")
	}

	return nil
}

// NewArguments creates new CreatePersonArguments instanse
func NewArguments() server.Arguments {
	return &Arguments{}
}
