package base

import (
	"github.com/GDVFox/tenjin/userd/database"
)

// PersonReply represents reply
type PersonReply struct {
	ID              int64  `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	DepartmentName  string `json:"department_name,omitempty"`
	AppointmentName string `json:"appointment_name,omitempty"`
}

type replyBuilder struct {
	persons         []*PersonReply
	lastCorrectWork map[int64]*database.WorkHistory
}

func newReplyBulder(size int) *replyBuilder {
	return &replyBuilder{
		persons: make([]*PersonReply, 0, size),
	}
}

func (b *replyBuilder) consumePersons(pp []*database.PersonModel) {
	for _, p := range pp {
		b.persons = append(b.persons, &PersonReply{
			ID:        p.ID,
			FirstName: p.FirstName,
			LastName:  p.LastName,
		})
	}
}

func (b *replyBuilder) consumeWorkHistory(hh []*database.WorkHistory) {
	b.lastCorrectWork = make(map[int64]*database.WorkHistory, len(hh))
	for _, h := range hh {
		if _, ok := b.lastCorrectWork[h.EmployeeID]; !ok {
			b.lastCorrectWork[h.EmployeeID] = h
		}
	}
}

func (b *replyBuilder) reply() []*PersonReply {
	if b.lastCorrectWork != nil {
		for _, p := range b.persons {
			if work, ok := b.lastCorrectWork[p.ID]; ok {
				p.DepartmentName = work.DepartmentName
				p.AppointmentName = work.AppointmentName
			}
		}
	}

	return b.persons
}
