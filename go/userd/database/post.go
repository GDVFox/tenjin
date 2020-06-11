package database

import (
	"time"

	"github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

// PostStatus is alias for status field in table
type PostStatus int32

// PostStatus is mapping of person.status column
const (
	PostStatusActive PostStatus = 1 + iota
	PostStatusDeleted
)

// PostModel represents abstract post in program
type PostModel struct {
	ID         int64     `db:"id"`
	Text       string    `db:"text"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	EmployeeID int64     `db:"employee_id"`
}

// CreatePost creates user post in database
func CreatePost(s dbr.SessionRunner, model *PostModel) error {
	res, err := s.InsertInto("employee_post").
		Columns("text", "employee_id").
		Record(model).Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// PostUpdateParams parameters for any post update
type PostUpdateParams struct {
	ID   int64
	Text *string
}

// UpdatePost updates person row in database
func UpdatePost(s dbr.SessionRunner, p *PostUpdateParams) error {
	q := s.Update("employee_post")
	if p.Text != nil {
		q.Set("text", *p.Text)
	}
	q.Set("updated_at", time.Now())
	q.Where(dbr.Neq("status", PostStatusDeleted))
	q.Where(dbr.Eq("id", p.ID))

	res, err := q.Exec()
	if err != nil {
		return nil
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return dbr.ErrNotFound
	}

	return nil
}
