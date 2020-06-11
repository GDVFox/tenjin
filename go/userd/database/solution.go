package database

import (
	"github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

// SolutionModel represents task in program
type SolutionModel struct {
	PostModel
	IsApproved bool  `db:"is_approved"`
	TaskID     int64 `db:"task_id"`
}

// CreateSolution creates user task in database
func CreateSolution(s dbr.SessionRunner, model *SolutionModel) error {
	res, err := s.InsertInto("solution").
		Pair("post_id", model.ID).
		Pair("task_id", model.TaskID).
		Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// SolutionOrder enum for ReadSolution ordering switch
type SolutionOrder uint8

const (
	// NewestSolutionOrder order by created_at desc
	NewestSolutionOrder SolutionOrder = iota
	// ApprovedNewestSolutionOrder order by created_at desc, but approved first
	ApprovedNewestSolutionOrder
)

// ReadSolution loads solutions to taskIDs with author authorIDs.
func ReadSolution(s dbr.SessionRunner, taskIDs, authorIDs []int64, approvedOnly bool, order SolutionOrder) ([]*SolutionModel, error) {
	q := s.Select(
		"p.id AS id", "p.text AS text", "p.rating AS rating",
		"p.created_at AS created_at ", "p.updated_at AS updated_at",
		"p.employee_id AS employee_id",
		"s.task_id AS task_id", "s.is_approved AS is_approved",
	).From(dbr.I("employee_post").As("p"))
	q.Join(dbr.I("solution").As("s"), "p.id=s.post_id")
	q.Where(dbr.Neq("p.status", PostStatusDeleted))
	if len(taskIDs) != 0 {
		q.Where(dbr.Eq("s.task_id", taskIDs))
	}
	if len(authorIDs) != 0 {
		q.Where(dbr.Eq("p.employee_id", authorIDs))
	}
	if approvedOnly {
		q.Where(dbr.Eq("s.is_approved", true))
	}

	if order == ApprovedNewestSolutionOrder {
		q.OrderDesc("s.is_approved")
	}
	q.OrderDesc("p.created_at")

	res := make([]*SolutionModel, 0)
	_, err := q.Load(&res)
	return res, err
}

// SolutionUpdateParams parameters for solution update
type SolutionUpdateParams struct {
	ID         int64
	IsApproved *bool
}

// UpdateSolution changes solution with p.id if needed.
func UpdateSolution(s dbr.SessionRunner, p *SolutionUpdateParams) error {
	q := s.Update("solution")
	if p.IsApproved != nil {
		q.Set("is_approved", *p.IsApproved)
	}
	q.Where(dbr.Eq("post_id", p.ID))

	_, err := q.Exec()
	return err
}

// DeleteSolution marks post row, connected with solution, in database as deleted
func DeleteSolution(s dbr.SessionRunner, id int64, taskID int64) error {
	if id == 0 && taskID == 0 {
		return ErrNoKeysSpecified
	}

	sqlQuery := "UPDATE employee_post AS p " +
		"INNER JOIN solution AS t " +
		"ON p.id=t.post_id SET p.status=? " +
		"WHERE p.status != ? AND "
	values := []interface{}{PostStatusDeleted, PostStatusDeleted}
	if id != 0 {
		sqlQuery += "p.id = ?"
		values = append(values, id)
	} else if taskID != 0 {
		sqlQuery += "t.task_id = ?"
		values = append(values, taskID)
	}

	q := s.UpdateBySql(sqlQuery, values...)
	res, err := q.Exec()
	if err != nil {
		return nil
	}

	return database.SomeAffected(res)
}
