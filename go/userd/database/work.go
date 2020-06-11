package database

import (
	"errors"
	"time"

	"github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

var (
	// ErrNoKeysSpecified no keys in function
	ErrNoKeysSpecified = errors.New("needs at least one key")
)

// WorkHistory job relation representation
type WorkHistory struct {
	EmployeeID      int64        `db:"employee_id"`
	DepartmentName  string       `db:"department_name"`
	AppointmentName string       `db:"appointment_name"`
	DateFrom        time.Time    `db:"date_from"`
	DateTo          dbr.NullTime `db:"date_to"`

	DepartmentID  int64 `db:"-"`
	AppointmentID int64 `db:"-"`
}

// CreateWorkHistory create job relation for user
func CreateWorkHistory(s dbr.SessionRunner, w *WorkHistory) error {
	res, err := s.InsertInto("works_in").
		Pair("employee_id", w.EmployeeID).
		Pair("department_id", w.DepartmentID).
		Pair("appointement_id", w.AppointmentID).
		Pair("date_from", w.DateFrom).
		Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// ReadWorkHistory read work history for person from database
func ReadWorkHistory(s dbr.SessionRunner, id []int64) ([]*WorkHistory, error) {
	q := s.Select(
		"w.employee_id AS employee_id",
		"d.name AS department_name",
		"a.name AS appointment_name",
		"w.date_from AS date_from",
		"w.date_to AS date_to",
		"w.date_to IS NULL AS working_now",
	).From(dbr.I("works_in").As("w"))
	q.Join(dbr.I("department").As("d"), "w.department_id=d.id")
	q.Join(dbr.I("appointement").As("a"), "w.appointement_id=a.id")
	if len(id) != 0 {
		q.Where(dbr.Eq("employee_id", id))
	}
	q.OrderDesc("working_now").OrderDesc("date_from")

	res := make([]*WorkHistory, 0, len(id))
	_, err := q.Load(&res)
	return res, err
}

// StopWork sets date_to in works_in table for some user job
func StopWork(s dbr.SessionRunner, employeeID, departmentID, appointmentID int64, dateTo time.Time) error {
	q := s.Update("works_in").
		Set("date_to", dateTo)
	q.Where(dbr.Eq("employee_id", employeeID))
	if departmentID != 0 {
		q.Where(dbr.Eq("department_id", departmentID))
	}
	if appointmentID != 0 {
		q.Where(dbr.Eq("appointement_id", appointmentID))
	}

	res, err := q.Exec()
	if err != nil {
		return nil
	}

	return database.SomeAffected(res)
}

// Department represents department row
type Department struct {
	ID   int64  `db:"id" json:"-"`
	Name string `db:"name" json:"name"`
}

// CreateDepartment creates department row in database
func CreateDepartment(s dbr.SessionRunner, name string) error {
	res, err := s.InsertInto("department").Pair("name", name).Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// ReadDepartment reads department row from database
func ReadDepartment(s dbr.SessionRunner, names []string) ([]*Department, error) {
	q := s.Select("id", "name").From("department")
	if len(names) != 0 {
		q.Where(dbr.Eq("name", names))
	}

	res := make([]*Department, 0, len(names))
	_, err := q.Load(&res)
	return res, err
}

// Appointment represents appointment row
type Appointment struct {
	ID   int64  `db:"id" json:"-"`
	Name string `db:"name" json:"name"`
}

// CreateAppointment creates appointement row in database
func CreateAppointment(s dbr.SessionRunner, name string) error {
	res, err := s.InsertInto("appointement").Pair("name", name).Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// ReadAppointment reads appointement row from database
func ReadAppointment(s dbr.SessionRunner, names []string) ([]*Appointment, error) {
	q := s.Select("id", "name").From("appointement")
	if len(names) != 0 {
		q.Where(dbr.Eq("name", names))
	}

	res := make([]*Appointment, 0, len(names))
	_, err := q.Load(&res)
	return res, err
}
