package database

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

// PersonStatus is alias for status field in table
type PersonStatus int32

// PersonStatus is mapping of person.status column
const (
	PersonStatusActive PersonStatus = 1 + iota
	PersonStatusBlocked
	PersonStatusDeleted
)

var statuses = map[PersonStatus][]byte{
	PersonStatusActive:  []byte("\"active\""),
	PersonStatusBlocked: []byte("\"blocked\""),
	PersonStatusDeleted: []byte("\"deleted\""),
}

// MarshalJSON returns status as JSON
func (s PersonStatus) MarshalJSON() ([]byte, error) {
	b, ok := statuses[s]
	if !ok {
		return nil, &json.MarshalerError{Type: reflect.TypeOf(s), Err: errors.New("unknown status")}
	}

	return b, nil
}

// PersonModel represents person entity in table
type PersonModel struct {
	ID        int64  `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`

	// optional fields
	Status    PersonStatus `db:"status" json:"status,omitempty"`
	CreatedAt time.Time    `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at,omitempty"`
}

// EmployeeModel represents employee entity in table
type EmployeeModel struct {
	PersonModel
	Email    string    `db:"email" json:"email"`
	Password []byte    `db:"password"`
	HiredAt  time.Time `db:"hired_at" json:"hired_at"`
}

// CreatePerson creates person row in database
func CreatePerson(s dbr.SessionRunner, person *PersonModel) error {
	res, err := s.InsertInto("person").
		Columns("first_name", "last_name").
		Record(person).Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// CreatedEmployee creates employee row in database
func CreatedEmployee(s dbr.SessionRunner, employee *EmployeeModel) error {
	res, err := s.InsertInto("employee").
		Pair("person_id", employee.ID).
		Pair("email", employee.Email).
		Pair("password", employee.Password).
		Pair("hired_at", employee.HiredAt).
		Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// ReadPerson reades person row in database
func ReadPerson(s dbr.SessionRunner, id []int64, employeesOnly bool) ([]*PersonModel, error) {
	q := readPerson(s, []string{"p.id", "p.first_name", "p.last_name"}, id, employeesOnly)
	res := make([]*PersonModel, 0, len(id))
	_, err := q.Load(&res)
	return res, err
}

// ReadPersonFull reades person row in database, if exists connected employee, reads employee
func ReadPersonFull(s dbr.SessionRunner, id []int64) ([]*EmployeeModel, error) {
	columns := []string{
		"p.id",
		"p.first_name", "p.last_name",
		"CAST(p.status as unsigned) AS status",
		"p.created_at", "p.updated_at",
		"e.email", "e.hired_at",
	}
	q := readPerson(s, columns, id, true)
	res := make([]*EmployeeModel, 0, len(id))
	_, err := q.Load(&res)
	return res, err
}

func readPerson(s dbr.SessionRunner, columns []string, id []int64, employeesOnly bool) *dbr.SelectStmt {
	q := s.Select(columns...).From(dbr.I("person").As("p"))
	if employeesOnly {
		q.Join(dbr.I("employee").As("e"), "p.id=e.person_id")
	}

	q.Where(dbr.Neq("p.status", PersonStatusDeleted))
	if len(id) != 0 {
		q.Where(dbr.Eq("p.id", id))
	}

	return q
}

// PersonUpdateParams parameters for user update
type PersonUpdateParams struct {
	ID        int64
	Blocked   *bool
	FirstName *string
	LastName  *string
}

// UpdatePerson updates person row in database
func UpdatePerson(s dbr.SessionRunner, p *PersonUpdateParams) error {
	q := s.Update("person")
	if p.Blocked != nil {
		status := PersonStatusActive
		if *p.Blocked {
			status = PersonStatusBlocked
		}
		q.Set("status", status)
	}
	if p.FirstName != nil {
		q.Set("first_name", *p.FirstName)
	}
	if p.LastName != nil {
		q.Set("last_name", *p.LastName)
	}
	// because if employee needs to be updated
	q.Set("updated_at", time.Now())
	q.Where(dbr.Eq("id", p.ID))

	res, err := q.Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// DeletePerson marks person row in database as deleted
func DeletePerson(s dbr.SessionRunner, id int64) error {
	q := s.Update("person")
	q.Set("status", PersonStatusDeleted)
	q.Where(dbr.Neq("status", PersonStatusDeleted))
	q.Where(dbr.Eq("id", id))

	res, err := q.Exec()
	if err != nil {
		return nil
	}

	return database.AssertAffected(res, 1)
}
