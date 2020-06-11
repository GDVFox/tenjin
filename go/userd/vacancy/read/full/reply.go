package full

import (
	"time"

	"github.com/GDVFox/tenjin/userd/database"
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

// InterviewItem represents interview object
type InterviewItem struct {
	ID          int64                    `json:"id"`
	PlanedDate  time.Time                `json:"planned_date"`
	Status      database.InterviewStatus `json:"status"`
	Person      InterviewPerson          `json:"person"`
	Interviewer *Interviewer             `json:"interviewer,omitempty"`
}

// SkillItem represents skill requirement
type SkillItem struct {
	Name       string                   `json:"name"`
	Difficulty database.SkillDifficulty `json:"difficulty"`
}

// VacancyItem represents vacancy in reply
type VacancyItem struct {
	ID              int64                    `json:"id"`
	Description     string                   `json:"description"`
	Priority        database.VacancyPriority `json:"priority"`
	Status          database.VacancyStatus   `json:"status"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	DepartmentName  string                   `json:"department_name"`
	AppointmentName string                   `json:"appointement_name"`
	Skills          []*SkillItem             `json:"require_skills,omitempty"`
	Interviews      []*InterviewItem         `json:"interviews,omitempty"`
}

type replyBuilder struct {
	vacancy    *VacancyItem
	skills     []*SkillItem
	interviews []*InterviewItem
}

func newReplyBuilder() *replyBuilder {
	return &replyBuilder{}
}

func (b *replyBuilder) consumeVacancy(vv *database.VacancyModel) {
	b.vacancy = &VacancyItem{
		ID:              vv.ID,
		Description:     vv.Description,
		Priority:        *vv.Priority,
		Status:          vv.Status,
		CreatedAt:       vv.CreatedAt,
		UpdatedAt:       vv.UpdatedAt,
		DepartmentName:  vv.DepartmentName,
		AppointmentName: vv.AppointmentName,
	}
}

func (b *replyBuilder) consumeSkillRequirements(ss []*database.SkillRequirement) {
	b.skills = make([]*SkillItem, 0, len(ss))
	for _, s := range ss {
		b.skills = append(b.skills, &SkillItem{
			Name:       s.SkillName,
			Difficulty: *s.Difficulty,
		})
	}
}

func (b *replyBuilder) consumeInterviews(ii []*database.InterviewModel) {
	b.interviews = make([]*InterviewItem, 0, len(ii))
	for _, i := range ii {
		item := &InterviewItem{
			ID:         i.ID,
			PlanedDate: i.PlanedDate,
			Status:     i.Status,
			Person: InterviewPerson{
				ID:        i.PersonID,
				FirstName: i.Person.FirstName,
				LastName:  i.Person.LastName,
			},
		}

		if i.InterviewerID.Valid {
			item.Interviewer = &Interviewer{
				InterviewPerson: InterviewPerson{
					ID:        i.InterviewerID.Int64,
					FirstName: i.Interviewer.FirstName.String,
					LastName:  i.Interviewer.LastName.String,
				},
				Email: i.Interviewer.Email.String,
			}
		}

		b.interviews = append(b.interviews, item)
	}
}

func (b *replyBuilder) reply() *VacancyItem {
	b.vacancy.Interviews = b.interviews
	b.vacancy.Skills = b.skills
	return b.vacancy
}
