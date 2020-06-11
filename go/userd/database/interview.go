package database

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

// InterviewStatus is alias for status field in table
type InterviewStatus int32

// InterviewStatus is mapping of interview.status column
const (
	InterviewStatusWaiting InterviewStatus = 1 + iota
	InterviewStatusCanceled
	InterviewStatusCompleted
)

var interviewStatuses = map[InterviewStatus][]byte{
	InterviewStatusWaiting:   []byte("\"waiting\""),
	InterviewStatusCanceled:  []byte("\"canceled\""),
	InterviewStatusCompleted: []byte("\"completed\""),
}

// MarshalJSON returns status as JSON
func (s InterviewStatus) MarshalJSON() ([]byte, error) {
	b, ok := interviewStatuses[s]
	if !ok {
		return nil, &json.MarshalerError{Type: reflect.TypeOf(s), Err: errors.New("unknown status")}
	}

	return b, nil
}

// InterviewPerson represents person_id entity
type InterviewPerson struct {
	FirstName string `db:"person_first_name"`
	LastName  string `db:"person_last_name"`
}

// Interviewer represents interviewer_id entity
type Interviewer struct {
	Email     dbr.NullString `db:"interviewer_email"`
	FirstName dbr.NullString `db:"interviewer_first_name"`
	LastName  dbr.NullString `db:"interviewer_last_name"`
}

// InterviewVacancy represents vacancy_id entity
type InterviewVacancy struct {
	DepartmentName  string `db:"vacancy_department_name"`
	AppointmentName string `db:"vacancy_appointement_name"`
}

// InterviewModel represents interview row in database
type InterviewModel struct {
	ID            int64           `db:"id"`
	Comment       dbr.NullString  `db:"comment"`
	PlanedDate    time.Time       `db:"planned_date"`
	Status        InterviewStatus `db:"status"`
	TotalScore    dbr.NullInt64   `db:"total_score"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
	PersonID      int64           `db:"person_id"`
	InterviewerID dbr.NullInt64   `db:"interviewer_id"`
	VacancyID     int64           `db:"vacancy_id"`

	Person      InterviewPerson
	Interviewer Interviewer
	Vacancy     InterviewVacancy
}

// CreateInterview creates interview row in database
func CreateInterview(s dbr.SessionRunner, interview *InterviewModel) error {
	res, err := s.InsertInto("interview").
		Columns("planned_date", "person_id", "vacancy_id").
		Record(interview).
		Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// ReadInterviewByVacancy reads short info about interview
func ReadInterviewByVacancy(s dbr.SessionRunner, vacancyID int64) ([]*InterviewModel, error) {
	q := s.Select(
		"i.id", "i.planned_date", "CAST(i.status AS unsigned) AS status",
		"p.id AS person_id", "p.first_name AS person_first_name", "p.last_name AS person_last_name",
		"pe.id AS interviewer_id", "e.email AS interviewer_email",
		"pe.first_name AS interviewer_first_name", "pe.last_name AS interviewer_last_name",
	).From(dbr.I("interview").As("i"))
	q.Join(dbr.I("person").As("p"), "i.person_id=p.id")
	q.LeftJoin(dbr.I("employee").As("e"), "i.interviewer_id=e.person_id")
	q.LeftJoin(dbr.I("person").As("pe"), "e.person_id=pe.id")
	q.Where(dbr.Eq("i.vacancy_id", vacancyID))

	res := make([]*InterviewModel, 0)
	_, err := q.Load(&res)
	return res, err
}

// ReadInterview reads short info about interview
func ReadInterview(s dbr.SessionRunner, id int64) ([]*InterviewModel, error) {
	q := s.Select(
		"i.id", "i.planned_date", "CAST(i.status AS unsigned) AS status", "i.comment",
		"i.total_score", "i.created_at", "i.updated_at", "i.person_id", "i.interviewer_id", "i.vacancy_id",
		"p.id AS person_id", "p.first_name AS person_first_name", "p.last_name AS person_last_name",
		"pe.id AS interviewer_id", "e.email AS interviewer_email",
		"pe.first_name AS interviewer_first_name", "pe.last_name AS interviewer_last_name",
		"v.id AS vacancy_id", "d.name AS vacancy_department_name", "a.name AS vacancy_appointement_name",
	).From(dbr.I("interview").As("i"))
	q.Join(dbr.I("person").As("p"), "i.person_id=p.id")
	q.LeftJoin(dbr.I("employee").As("e"), "i.interviewer_id=e.person_id")
	q.LeftJoin(dbr.I("person").As("pe"), "e.person_id=pe.id")
	q.Join(dbr.I("vacancy").As("v"), "i.vacancy_id=v.id")
	q.Join(dbr.I("department").As("d"), "v.department_id=d.id")
	q.Join(dbr.I("appointement").As("a"), "v.appointement_id=a.id")
	q.Where(dbr.Eq("i.id", id))

	res := make([]*InterviewModel, 0)
	_, err := q.Load(&res)
	return res, err
}

// InterviewShort short info about interview
type InterviewShort struct {
	ID        int64           `db:"id"`
	VacancyID int64           `db:"vacancy_id"`
	Status    InterviewStatus `db:"status"`
}

// ReadInterviewVacancy reads vacancy_id for interview with id
func ReadInterviewVacancy(s dbr.SessionRunner, id int64) (*InterviewShort, error) {
	q := s.Select("id", "vacancy_id", "CAST(status AS unsigned) AS status").
		From("interview").
		Where(dbr.Eq("id", id))

	interview := &InterviewShort{}
	err := q.LoadOne(&interview)
	return interview, err
}

// InterviewUpdateParams parameters for interview update
type InterviewUpdateParams struct {
	ID            int64
	PlannedDate   *time.Time
	Comment       *string
	TotalScore    *int64
	InterviewerID *int64
}

// UpdateInterview updates common interview params
func UpdateInterview(s dbr.SessionRunner, params *InterviewUpdateParams) error {
	q := s.Update("interview")
	if params.PlannedDate != nil {
		q.Set("planned_date", *params.PlannedDate)
	}
	if params.Comment != nil {
		q.Set("comment", *params.Comment)
	}
	if params.TotalScore != nil {
		q.Set("total_score", *params.TotalScore)
	}
	if params.InterviewerID != nil {
		q.Set("interviewer_id", *params.InterviewerID)
	}
	q.Set("updated_at", time.Now())
	q.Where(dbr.Neq("status", InterviewStatusCanceled))
	q.Where(dbr.Eq("id", params.ID))

	res, err := q.Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// InterviewCompleteParams parameters for interview complete
type InterviewCompleteParams struct {
	ID            int64
	Comment       string
	TotalScore    int64
	InterviewerID int64
}

// CompleteInterview sets interview result
func CompleteInterview(s dbr.SessionRunner, params *InterviewCompleteParams) error {
	res, err := s.Update("interview").
		Set("comment", params.Comment).
		Set("total_score", params.TotalScore).
		Set("interviewer_id", params.InterviewerID).
		Set("status", InterviewStatusCompleted).
		Where(dbr.Eq("id", params.ID)).
		Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// CancelInterview sets interview status canceled
func CancelInterview(s dbr.SessionRunner, id int64) error {
	res, err := s.Update("interview").
		Set("status", InterviewStatusCanceled).
		Where(dbr.Eq("status", InterviewStatusWaiting)).
		Where(dbr.Eq("id", id)).
		Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}
