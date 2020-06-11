package database

import (
	"github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

// RequirementCheckModel represents requirement_check row in database
type RequirementCheckModel struct {
	VacancyID   int64          `db:"vacancy_id"`
	InterviewID int64          `db:"interview_id"`
	TaskID      int64          `db:"task_id"`
	SkillID     int64          `db:"skill_id"`
	Comment     dbr.NullString `db:"comment"`
	Score       int64          `db:"score"`

	SkillName string `db:"skill_name"`
}

// CreateRequirementCheck creates requirement_check rows in database
func CreateRequirementCheck(s dbr.SessionRunner, models []*RequirementCheckModel) error {
	q := s.InsertInto("requirement_check").
		Columns("vacancy_id", "interview_id", "task_id", "skill_id", "comment", "score")
	for _, m := range models {
		q.Record(m)
	}

	res, err := q.Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, int64(len(models)))
}

// ReadRequirementCheck loads checks data from database
func ReadRequirementCheck(s dbr.SessionRunner, interviewID int64) ([]*RequirementCheckModel, error) {
	q := s.Select(
		"c.task_id", "c.comment",
		"c.score", "s.name AS skill_name",
	).From(dbr.I("requirement_check").As("c"))
	q.Join(dbr.I("skill").As("s"), "c.skill_id=s.id")
	q.Where(dbr.Eq("c.interview_id", interviewID))

	res := make([]*RequirementCheckModel, 0)
	_, err := q.Load(&res)
	return res, err
}

// DeleteRequirementCheck deletes requirement_check
func DeleteRequirementCheck(s dbr.SessionRunner, models []*RequirementCheckModel) error {
	sqlQuery := "DELETE FROM `requirement_check` WHERE "
	values := make([]interface{}, 0, len(models)*4)
	for i, m := range models {
		sqlQuery += "(vacancy_id = ? AND interview_id = ? AND task_id = ? AND skill_id = ?)"
		values = append(values, m.VacancyID, m.InterviewID, m.TaskID, m.SkillID)
		if i != len(models)-1 {
			sqlQuery += " OR "
		}
	}

	res, err := s.DeleteBySql(sqlQuery, values...).Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, int64(len(models)))
}
