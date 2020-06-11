package database

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

// VacancyPriority is alias for priority field in table
type VacancyPriority int32

// VacancyPriority is mapping of vacancy.priority column
const (
	VacancyPriorityLow VacancyPriority = 1 + iota
	VacancyPriorityMedium
	VacancyPriorityHigh
)

var vacancyPriorities = map[VacancyPriority][]byte{
	VacancyPriorityLow:    []byte("\"low\""),
	VacancyPriorityMedium: []byte("\"medium\""),
	VacancyPriorityHigh:   []byte("\"high\""),
}

// MarshalJSON returns status as JSON
func (s VacancyPriority) MarshalJSON() ([]byte, error) {
	b, ok := vacancyPriorities[s]
	if !ok {
		return nil, &json.MarshalerError{Type: reflect.TypeOf(s), Err: errors.New("unknown priority")}
	}

	return b, nil
}

// VacancyStatus is alias for status field in table
type VacancyStatus int32

// VacancyStatus is mapping of vacancy.status column
const (
	VacancyStatusActive VacancyStatus = 1 + iota
	VacancyStatusPaused
	VacancyStatusClosed
)

var vacancyStatuses = map[VacancyStatus][]byte{
	VacancyStatusActive: []byte("\"active\""),
	VacancyStatusPaused: []byte("\"paused\""),
	VacancyStatusClosed: []byte("\"closed\""),
}

// MarshalJSON returns status as JSON
func (s VacancyStatus) MarshalJSON() ([]byte, error) {
	b, ok := vacancyStatuses[s]
	if !ok {
		return nil, &json.MarshalerError{Type: reflect.TypeOf(s), Err: errors.New("unknown status")}
	}

	return b, nil
}

// VacancyModel represent vacancy in database
type VacancyModel struct {
	ID              int64            `db:"id"`
	Description     string           `db:"description"`
	Priority        *VacancyPriority `db:"priority"`
	Status          VacancyStatus    `db:"status"`
	CreatedAt       time.Time        `db:"created_at"`
	UpdatedAt       time.Time        `db:"updated_at"`
	DepartmentName  string           `db:"department_name"`
	AppointmentName string           `db:"appointement_name"`

	DepartmentID  int64 `db:"department_id"`
	AppointmentID int64 `db:"appointement_id"`
}

// CreateVacancy creates vacancy row in database
func CreateVacancy(s dbr.SessionRunner, model *VacancyModel) error {
	columns := []string{"description", "department_id", "appointement_id"}
	if model.Priority != nil {
		columns = append(columns, "priority")
	}

	res, err := s.InsertInto("vacancy").
		Columns(columns...).
		Record(model).
		Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// ReadVacancy reads info about vacancy
func ReadVacancy(s dbr.SessionRunner, ids []int64, departments, appointments []string) ([]*VacancyModel, error) {
	q := s.Select(
		"v.id", "v.description",
		"CAST(v.priority AS unsigned) AS priority",
		"CAST(v.status AS unsigned) AS status",
		"v.created_at", "v.updated_at",
		"d.name AS department_name", "a.name AS appointement_name",
	).From(dbr.I("vacancy").As("v"))
	q.Join(dbr.I("department").As("d"), "v.department_id=d.id")
	q.Join(dbr.I("appointement").As("a"), "v.appointement_id=a.id")

	if len(ids) != 0 {
		q.Where(dbr.Eq("v.id", ids))
	}
	if len(departments) != 0 {
		q.Where(dbr.Eq("d.name", departments))
	}
	if len(appointments) != 0 {
		q.Where(dbr.Eq("a.name", appointments))
	}
	q.OrderAsc("status")
	q.OrderDesc("priority")

	res := make([]*VacancyModel, 0, len(ids))
	_, err := q.Load(&res)
	return res, err
}

// VacancyUpdateParams parameters for vacancy update
type VacancyUpdateParams struct {
	ID            int64
	Description   *string
	Priority      *VacancyPriority
	Pause         *bool
	DepartmentID  *int64
	AppointmentID *int64
}

// UpdateVacancy updated vacacny row in database
func UpdateVacancy(s dbr.SessionRunner, params *VacancyUpdateParams) error {
	q := s.Update("vacancy")
	if params.Description != nil {
		q.Set("description", *params.Description)
	}
	if params.Priority != nil {
		q.Set("priority", *params.Priority)
	}
	if params.Pause != nil {
		if *params.Pause {
			q.Set("status", VacancyStatusPaused)
		} else {
			q.Set("status", VacancyStatusActive)
		}
	}
	if params.DepartmentID != nil {
		q.Set("department_id", *params.DepartmentID)
	}
	if params.AppointmentID != nil {
		q.Set("appointement_id", *params.AppointmentID)
	}
	q.Set("updated_at", time.Now())
	q.Where(dbr.Neq("status", VacancyStatusClosed))
	q.Where(dbr.Eq("id", params.ID))

	res, err := q.Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// CloseVacancy closes vacancy in database
func CloseVacancy(s dbr.SessionRunner, id int64) error {
	res, err := s.Update("vacancy").Set("status", VacancyStatusClosed).
		Where(dbr.Neq("status", VacancyStatusClosed)).
		Where(dbr.Eq("id", id)).
		Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}
