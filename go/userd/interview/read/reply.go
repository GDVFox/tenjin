package read

import (
	"time"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
)

// InterviewPerson represents person object
type InterviewPerson struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Interviewer represents interviewer entity
type Interviewer struct {
	InterviewPerson
	Email string `json:"email"`
}

// Vacancy represents vacanct entity
type Vacancy struct {
	ID              int64  `json:"id"`
	DepartmentName  string `json:"department_name"`
	AppointmentName string `json:"appointement_name"`
}

// CheckItem represents vacancy skill check
type CheckItem struct {
	TaskID    int64  `json:"task_id"`
	SkillName string `json:"skill_name"`
	Comment   string `json:"comment"`
	Score     int64  `json:"score"`
}

// InterviewItem represents interview object
type InterviewItem struct {
	ID          int64                    `json:"id"`
	Comment     *string                  `json:"comment,omitempty"`
	PlanedDate  time.Time                `json:"planned_date"`
	Status      database.InterviewStatus `json:"status"`
	TotalScore  *int64                   `json:"total_score,omitempty"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	Person      InterviewPerson          `json:"person"`
	Vacancy     Vacancy                  `json:"vacancy"`
	Interviewer *Interviewer             `json:"interviewer,omitempty"`
	Checks      []*CheckItem             `json:"checks,omitempty"`
}

type replyBuilder struct {
	interview *InterviewItem
	checks    []*CheckItem
}

func newReplyBuilder() *replyBuilder {
	return &replyBuilder{}
}

func (b *replyBuilder) consumeInterview(i *database.InterviewModel) {
	b.interview = &InterviewItem{
		ID:         i.ID,
		Comment:    db.NullStringToPtr(i.Comment),
		PlanedDate: i.PlanedDate,
		Status:     i.Status,
		TotalScore: db.NullInt64ToPtr(i.TotalScore),
		CreatedAt:  i.CreatedAt,
		UpdatedAt:  i.UpdatedAt,
		Person: InterviewPerson{
			ID:        i.PersonID,
			FirstName: i.Person.FirstName,
			LastName:  i.Person.LastName,
		},
		Vacancy: Vacancy{
			ID:              i.VacancyID,
			DepartmentName:  i.Vacancy.DepartmentName,
			AppointmentName: i.Vacancy.AppointmentName,
		},
	}

	if i.InterviewerID.Valid {
		b.interview.Interviewer = &Interviewer{
			InterviewPerson: InterviewPerson{
				ID:        i.InterviewerID.Int64,
				FirstName: i.Interviewer.FirstName.String,
				LastName:  i.Interviewer.LastName.String,
			},
			Email: i.Interviewer.Email.String,
		}
	}
}

func (b *replyBuilder) consumeChecks(cc []*database.RequirementCheckModel) {
	b.checks = make([]*CheckItem, 0, len(cc))
	for _, c := range cc {
		ch := &CheckItem{
			TaskID:    c.TaskID,
			SkillName: c.SkillName,
			Comment:   c.Comment.String,
			Score:     c.Score,
		}

		b.checks = append(b.checks, ch)
	}
}

func (b *replyBuilder) reply() interface{} {
	b.interview.Checks = b.checks
	return b.interview
}
