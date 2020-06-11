package database

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

// SkillDifficulty is alias for difficulty field in table
type SkillDifficulty int32

// SkillDifficulty is mapping of task_skill_requirement.difficulty column
const (
	LowSkillDifficulty SkillDifficulty = 1 + iota
	MediumSkillDifficulty
	HardSkillDifficulty
	UnsolvableSkillDifficulty
)

var skillDifficulties = map[SkillDifficulty][]byte{
	LowSkillDifficulty:        []byte("\"low\""),
	MediumSkillDifficulty:     []byte("\"medium\""),
	HardSkillDifficulty:       []byte("\"hard\""),
	UnsolvableSkillDifficulty: []byte("\"unsolvable\""),
}

// MarshalJSON returns status as JSON
func (d SkillDifficulty) MarshalJSON() ([]byte, error) {
	b, ok := skillDifficulties[d]
	if !ok {
		return nil, &json.MarshalerError{Type: reflect.TypeOf(d), Err: errors.New("unknown difficulty")}
	}

	return b, nil
}

// SkillModel represents skill rows in skill table
type SkillModel struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

// CreateSkills creates new skill with name
func CreateSkills(s dbr.SessionRunner, skills []string, ignore bool) error {
	q := s.InsertInto("skill").Columns("name")
	for _, s := range skills {
		q.Values(s)
	}

	if ignore {
		q.Ignore()
	}

	res, err := q.Exec()
	if err != nil {
		return err
	}

	if ignore {
		return nil
	}

	return database.SomeAffected(res)
}

// ReadSkill loads skills with names from skill table
func ReadSkill(s dbr.SessionRunner, names []string) ([]*SkillModel, error) {
	q := s.Select("id", "name").From("skill")
	if len(names) != 0 {
		q.Where(dbr.Eq("name", names))
	}

	res := make([]*SkillModel, 0, len(names))
	_, err := q.Load(&res)
	return res, err
}

// SkillRequirement represents task or vacanct skill requirement
type SkillRequirement struct {
	TaskID    int64 `db:"task_id"`
	VacancyID int64 `db:"vacancy_id"`

	SkillID    int64            `db:"skill_id"`
	SkillName  string           `db:"skill_name"`
	Difficulty *SkillDifficulty `db:"difficulty"`
}

// CreateTaskSkillRequirement creates connection between task and skills
func CreateTaskSkillRequirement(s dbr.SessionRunner, skills []*SkillRequirement) error {
	return createSkillRequirement(s, "task_skill_requirement", "task_id", skills)
}

// CreateVacancySkillRequirement creates connection between task and skills
func CreateVacancySkillRequirement(s dbr.SessionRunner, skills []*SkillRequirement) error {
	return createSkillRequirement(s, "vacancy_skill_requirement", "vacancy_id", skills)
}

// ReadTaskSkillRequirements reads skills required for task
func ReadTaskSkillRequirements(s dbr.SessionRunner, taskIDs []int64) ([]*SkillRequirement, error) {
	return readSkillRequirements(s, "task_skill_requirement", "task_id", taskIDs)
}

// ReadVacancyRequirements reads skills required for vacancy
func ReadVacancyRequirements(s dbr.SessionRunner, vacancyIDs []int64) ([]*SkillRequirement, error) {
	return readSkillRequirements(s, "vacancy_skill_requirement", "vacancy_id", vacancyIDs)
}

// DeleteTaskSkillRequirements reads skills required for task
func DeleteTaskSkillRequirements(s dbr.SessionRunner, skillIDs []int64) error {
	return deleteSkillRequirements(s, "task_skill_requirement", skillIDs)
}

// DeleteVacancySkillRequirements reads skills required for task
func DeleteVacancySkillRequirements(s dbr.SessionRunner, skillIDs []int64) error {
	return deleteSkillRequirements(s, "vacancy_skill_requirement", skillIDs)
}

func createSkillRequirement(s dbr.SessionRunner, table, foreign string, skills []*SkillRequirement) error {
	sqlQuery := "INSERT INTO " + table + " (" + foreign + ", skill_id, difficulty) VALUES "
	values := make([]interface{}, 0, len(skills))
	for i, s := range skills {
		if s.Difficulty == nil {
			sqlQuery += "(?, ?, DEFAULT)"
		} else {
			sqlQuery += "(?, ?, ?)"
		}

		if i != len(skills)-1 {
			sqlQuery += ", "
		}

		if foreign == "task_id" {
			values = append(values, s.TaskID)
		} else if foreign == "vacancy_id" {
			values = append(values, s.VacancyID)
		}
		values = append(values, s.SkillID)
		if s.Difficulty != nil {
			values = append(values, *s.Difficulty)
		}
	}

	res, err := s.InsertBySql(sqlQuery, values...).Exec()
	if err != nil {
		return err
	}

	return database.SomeAffected(res)
}

func readSkillRequirements(s dbr.SessionRunner, table, foreign string, foreignIDs []int64) ([]*SkillRequirement, error) {
	q := s.Select(
		"r."+foreign+" AS "+foreign,
		"s.name AS skill_name",
		"CAST(r.difficulty as unsigned) AS difficulty",
	).From(dbr.I(table).As("r"))
	q.Join(dbr.I("skill").As("s"), "r.skill_id=s.id")
	q.Where(dbr.Eq(foreign, foreignIDs))

	res := make([]*SkillRequirement, 0)
	_, err := q.Load(&res)
	return res, err
}

func deleteSkillRequirements(s dbr.SessionRunner, table string, skillIDs []int64) error {
	q := s.DeleteFrom(table)
	if len(skillIDs) != 0 {
		q.Where(dbr.Eq("skill_id", skillIDs))
	}

	res, err := q.Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, int64(len(skillIDs)))
}
