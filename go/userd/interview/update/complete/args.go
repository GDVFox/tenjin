package complete

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/server"
)

// Arguments represents update args
type Arguments struct {
	ID            int64                                  `json:"id"`
	InterviewerID int64                                  `json:"interviewer_id"`
	Comment       string                                 `json:"comment"`
	TotalScore    int64                                  `json:"total_score"`
	Checks        []*args.RequirementCheckCreateArgument `json:"checks"`

	vacancyID int64
}

// Validate represents arguments for task update
func (a *Arguments) Validate() error {
	if a.Comment == "" {
		return server.NewHTTPError(http.StatusBadRequest, "comment can not be empty")
	}

	if a.TotalScore < 0 || a.TotalScore > 100 {
		return server.NewHTTPError(http.StatusBadRequest, "total_score must be in [0; 100]")
	}

	if a.InterviewerID == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "interviewer_id can not be 0")
	}

	for _, c := range a.Checks {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Resolve loads additional data
func (a *Arguments) Resolve(s *db.Session, r *http.Request) error {
	interviewInfo, err := database.ReadInterviewVacancy(s, a.ID)
	if err != nil {
		return err
	}
	a.vacancyID = interviewInfo.VacancyID
	return nil
}

// NewArguments creates new update arguments
func NewArguments() server.Arguments {
	return &Arguments{}
}
