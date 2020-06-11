package database

import (
	"github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

// VoteModel represents vote in database
type VoteModel struct {
	ID         int64         `db:"id"`
	EmployeeID int64         `db:"employee_id"`
	PostID     dbr.NullInt64 `db:"post_id"`
	CommentID  dbr.NullInt64 `db:"comment_id"`
	Delta      int           `db:"delta"`
}

// CreateVote creates vote row in database
func CreateVote(s dbr.SessionRunner, model *VoteModel) error {
	sqlQuery := "INSERT INTO `vote`(`employee_id`, `post_id`, `comment_id`, `delta`) " +
		"VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE delta = delta + ?"
	res, err := s.InsertBySql(sqlQuery, model.EmployeeID, model.PostID,
		model.CommentID, model.Delta, model.Delta).Exec()
	if err != nil {
		return err
	}

	return database.SomeAffected(res)
}
