package database

import (
	"time"

	"github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

// CommentStatus is alias for status field in table
type CommentStatus int32

// CommentStatus is mapping of comment.status column
const (
	CommentStatusActive CommentStatus = 1 + iota
	CommentStatusDeleted
)

// CommentModel represents comment
type CommentModel struct {
	ID         int64         `db:"id"`
	Text       string        `db:"text"`
	Rating     int64         `db:"rating"`
	Status     CommentStatus `db:"status"`
	CreatedAt  time.Time     `db:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at"`
	EmployeeID int64         `db:"employee_id"`
	PostID     int64         `db:"post_id"`
	Parent     dbr.NullInt64 `db:"parent"`
}

// CreateComment creates new comment object
func CreateComment(s dbr.SessionRunner, comment *CommentModel) error {
	res, err := s.InsertInto("comment").
		Columns("post_id", "employee_id", "text", "parent").
		Record(comment).
		Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// ReadComments loads comment from database
func ReadComments(s dbr.SessionRunner, employeeID int64, postIDs []int64) ([]*CommentModel, error) {
	q := s.Select(
		"id", "text", "rating", "CAST(status AS unsigned) AS status",
		"created_at", "updated_at",
		"employee_id", "post_id", "parent",
	).From("comment")
	if employeeID != 0 {
		q.Where(dbr.Eq("employee_id", employeeID))
	}
	if len(postIDs) != 0 {
		q.Where(dbr.Eq("post_id", postIDs))
	}

	res := []*CommentModel{}
	_, err := q.Load(&res)
	return res, err
}

// CommentUpdateParams parameters for comment update
type CommentUpdateParams struct {
	ID   int64
	Text *string
}

// UpdateComment changes solution with p.id if needed.
func UpdateComment(s dbr.SessionRunner, p *CommentUpdateParams) error {
	q := s.Update("comment")
	if p.Text != nil {
		q.Set("text", *p.Text)
	}
	q.Set("updated_at", time.Now())
	q.Where(dbr.Neq("status", CommentStatusDeleted))
	q.Where(dbr.Eq("id", p.ID))

	res, err := q.Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// DeleteComment remove comment
func DeleteComment(s dbr.SessionRunner, commentID int64) error {
	q := s.Update("comment").
		Set("status", CommentStatusDeleted).
		Where(dbr.Neq("status", CommentStatusDeleted)).
		Where(dbr.Eq("id", commentID))
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
