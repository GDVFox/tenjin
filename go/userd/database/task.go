package database

import (
	"github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

// TaskModel represents task in program
type TaskModel struct {
	PostModel
	Title string `db:"title"`
}

// CreateTask creates user task in database
func CreateTask(s dbr.SessionRunner, model *TaskModel) error {
	res, err := s.InsertInto("task").
		Pair("title", model.Title).
		Pair("post_id", model.ID).
		Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, 1)
}

// ReadTask reads tasks with taskID with employeeID as author
func ReadTask(s dbr.SessionRunner, taskID, employeeID []int64) ([]*TaskModel, error) {
	columns := []string{
		"p.id AS id", "t.title AS title", "p.rating AS rating",
		"p.text AS text", "p.employee_id AS employee_id",
	}
	q := readTaks(s, columns, taskID, employeeID)
	res := make([]*TaskModel, 0)
	_, err := q.Load(&res)
	return res, err
}

// ReadTaskFull reads FULL task with taskID with employeeID as author
func ReadTaskFull(s dbr.SessionRunner, taskID int64) (*TaskModel, error) {
	columns := []string{
		"p.id AS id", "t.title AS title", "p.rating AS rating",
		"p.text AS text", "CAST(p.status as unsigned) AS status",
		"p.created_at AS created_at", "p.updated_at AS updated_at",
		"p.employee_id AS employee_id",
	}
	q := readTaks(s, columns, []int64{taskID}, nil)
	res := &TaskModel{}
	err := q.LoadOne(&res)
	return res, err
}

func readTaks(s dbr.SessionRunner, columns []string, taskID, employeeID []int64) *dbr.SelectStmt {
	q := s.Select(columns...).From(dbr.I("employee_post").As("p"))
	q.Join(dbr.I("task").As("t"), "p.id=t.post_id")
	q.Where(dbr.Neq("p.status", PostStatusDeleted))
	if len(taskID) != 0 {
		q.Where(dbr.Eq("p.id", taskID))
	}
	if len(employeeID) != 0 {
		q.Where(dbr.Eq("p.employee_id", employeeID))
	}

	return q
}

// TaskUpdateParams parameters for task update
type TaskUpdateParams struct {
	ID    int64
	Title *string
}

// UpdateTask updates person row in database
func UpdateTask(s dbr.SessionRunner, p *TaskUpdateParams) error {
	q := s.Update("task")
	if p.Title != nil {
		q.Set("title", *p.Title)
	}
	q.Where(dbr.Eq("post_id", p.ID))

	_, err := q.Exec()
	return err
}

// DeleteTask marks post row, connected with task, in database as deleted
func DeleteTask(s dbr.SessionRunner, id int64) error {
	q := s.DeleteBySql("UPDATE employee_post AS p "+
		"INNER JOIN task AS t "+
		"ON p.id=t.post_id SET p.status=? "+
		"WHERE p.id=? AND p.status != ?",
		PostStatusDeleted, id, PostStatusDeleted,
	)

	res, err := q.Exec()
	if err != nil {
		return nil
	}

	return database.AssertAffected(res, 1)
}
