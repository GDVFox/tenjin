package base

import (
	"time"

	"github.com/GDVFox/tenjin/userd/database"
)

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
}

type replyBuilder struct {
	vacancies []*VacancyItem
}

func newReplyBuilder() *replyBuilder {
	return &replyBuilder{}
}

func (b *replyBuilder) consumeVacancies(vv []*database.VacancyModel) {
	b.vacancies = make([]*VacancyItem, 0, len(vv))
	for _, v := range vv {
		iv := &VacancyItem{
			ID:              v.ID,
			Description:     v.Description,
			Priority:        *v.Priority,
			Status:          v.Status,
			CreatedAt:       v.CreatedAt,
			UpdatedAt:       v.UpdatedAt,
			DepartmentName:  v.DepartmentName,
			AppointmentName: v.AppointmentName,
		}

		b.vacancies = append(b.vacancies, iv)
	}
}

func (b *replyBuilder) reply() []*VacancyItem {
	return b.vacancies
}
