package full

import (
	"time"

	"github.com/GDVFox/tenjin/userd/database"
)

// WorkHistoryItem represents reprsents work info
type WorkHistoryItem struct {
	DepartmentName  string     `json:"department_name"`
	AppointmentName string     `json:"appointment_name"`
	DateFrom        time.Time  `json:"date_from"`
	DateTo          *time.Time `json:"date_to,omitempty"`
}

// PersonFullReply represents reply
type PersonFullReply struct {
	ID        int64                 `json:"id"`
	FirstName string                `json:"first_name"`
	LastName  string                `json:"last_name"`
	Status    database.PersonStatus `json:"status"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	Email     string                `json:"email"`
	HiredAt   time.Time             `json:"hired_at"`
	Work      []*WorkHistoryItem    `json:"work"`
}

type replyBuilder struct {
	person *PersonFullReply
}

func newReplyBulder() *replyBuilder {
	return &replyBuilder{
		person: &PersonFullReply{},
	}
}

func (b *replyBuilder) consumePerson(e *database.EmployeeModel) {
	b.person.ID = e.ID
	b.person.FirstName = e.FirstName
	b.person.LastName = e.LastName
	b.person.Status = e.Status
	b.person.CreatedAt = e.CreatedAt
	b.person.UpdatedAt = e.UpdatedAt
	b.person.Email = e.Email
	b.person.HiredAt = e.HiredAt
}

func (b *replyBuilder) consumeWorkHistory(hh []*database.WorkHistory) {
	b.person.Work = make([]*WorkHistoryItem, 0, len(hh))
	for _, h := range hh {
		item := &WorkHistoryItem{
			DepartmentName:  h.DepartmentName,
			AppointmentName: h.AppointmentName,
			DateFrom:        h.DateFrom,
		}

		if h.DateTo.Valid {
			item.DateTo = &h.DateTo.Time
		}

		b.person.Work = append(b.person.Work, item)
	}
}

func (b *replyBuilder) reply() *PersonFullReply {
	return b.person
}
