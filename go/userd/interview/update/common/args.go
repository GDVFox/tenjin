package common

import (
	"net/http"
	"time"

	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/server"
)

// Arguments represents update args
type Arguments struct {
	ID            int64                                  `json:"id"`
	PlannedTime   *time.Time                             `json:"planned_time"`
	Comment       *string                                `json:"comment"`
	TotalScore    *int64                                 `json:"total_score"`
	InterviewerID *int64                                 `json:"interviewer_id"`
	Checks        []*args.RequirementCheckUpdateArgument `json:"checks"`

	vacancyID int64
	status    database.InterviewStatus
}

// Validate represents arguments for interview update
func (a *Arguments) Validate() error {
	if a.PlannedTime == nil && a.Comment == nil &&
		a.TotalScore == nil && a.InterviewerID == nil && len(a.Checks) == 0 {
		return server.EmptyArgumentsHTTPError
	}

	if a.PlannedTime != nil && a.PlannedTime.IsZero() {
		return server.NewHTTPError(http.StatusBadRequest, "planned time can not be zero")
	}

	if a.Comment != nil && *a.Comment == "" {
		return server.NewHTTPError(http.StatusBadRequest, "comment can not be empty")
	}

	if a.TotalScore != nil && (*a.TotalScore < 0 || *a.TotalScore > 100) {
		return server.NewHTTPError(http.StatusBadRequest, "total_score must be in [0; 100]")
	}

	if a.InterviewerID != nil && *a.InterviewerID == 0 {
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
	a.status = interviewInfo.Status

	if a.status != database.InterviewStatusCompleted &&
		(a.Comment != nil || a.TotalScore != nil || a.InterviewerID != nil || len(a.Checks) != 0) {
		return server.NewHTTPError(http.StatusBadRequest,
			"comment, total_score, interviewer_id, check should be updated only for completed interview")
	}

	return nil
}

// IsInterviewUpdated returns true if interview entity needs to be updated
func (a *Arguments) IsInterviewUpdated() bool {
	return a.PlannedTime != nil || a.Comment != nil ||
		a.TotalScore != nil || a.InterviewerID != nil
}

// IsChecksUpdated returns true if requirement_check entities needs to be updated
func (a *Arguments) IsChecksUpdated() bool {
	return len(a.Checks) != 0
}

// NewArguments creates new update arguments
func NewArguments() server.Arguments {
	return &Arguments{}
}
