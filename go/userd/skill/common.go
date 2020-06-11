package skill

import (
	"errors"

	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/userd/database"
	"github.com/gocraft/dbr/v2"
)

// RequirementType represents task or vacancy skill type
type RequirementType int32

// RequirementType values
const (
	TaskRequirementType = iota
	VacancyRequirementType
)

var (
	// ErrUnknownRequirementType represents error when RequirementType
	// is not TaskRequirementType or VacancyRequirementType
	ErrUnknownRequirementType = errors.New("unknown RequirementType")
)

// CreateSkills creates skill rows in database
func CreateSkills(s dbr.SessionRunner, t RequirementType, entityID int64, arg []*args.SkillCreate) error {
	skillNames := make([]string, 0, len(arg))
	for _, s := range arg {
		skillNames = append(skillNames, s.Name)
	}

	if err := database.CreateSkills(s, skillNames, true); err != nil {
		return err
	}

	skills, err := database.ReadSkill(s, skillNames)
	if err != nil {
		return err
	}

	knownSkills := make(map[string]int64, len(skills))
	for _, s := range skills {
		knownSkills[s.Name] = s.ID
	}

	reqs := make([]*database.SkillRequirement, 0, len(arg))
	for _, s := range arg {
		req := &database.SkillRequirement{
			SkillID: knownSkills[s.Name],
		}
		if t == TaskRequirementType {
			req.TaskID = entityID
		} else if t == VacancyRequirementType {
			req.VacancyID = entityID
		}
		if s.Difficulty != nil {
			req.Difficulty = &s.Difficulty.Value
		}

		reqs = append(reqs, req)
	}

	switch t {
	case TaskRequirementType:
		return database.CreateTaskSkillRequirement(s, reqs)
	case VacancyRequirementType:
		return database.CreateVacancySkillRequirement(s, reqs)
	default:
		return ErrUnknownRequirementType
	}
}

// UpdateSkills upated skill rows in database
func UpdateSkills(s dbr.SessionRunner, t RequirementType, entityID int64, arg []*args.SkillUpdate) error {
	createSkills := make([]*args.SkillUpdate, 0)
	deleteSkills := make([]*args.SkillUpdate, 0)
	for _, s := range arg {
		switch s.Type {
		case args.AddUpdateType:
			createSkills = append(createSkills, s)
		case args.DeleteUpdateType:
			deleteSkills = append(deleteSkills, s)
		}
	}

	if len(createSkills) != 0 {
		skillsToCreate := make([]*args.SkillCreate, 0, len(createSkills))
		for _, s := range createSkills {
			cr := &args.SkillCreate{
				Name:       s.Name,
				Difficulty: s.Difficulty,
			}
			skillsToCreate = append(skillsToCreate, cr)
		}

		if err := CreateSkills(s, t, entityID, skillsToCreate); err != nil {
			return err
		}
	}

	if len(deleteSkills) != 0 {
		skillsToDelete := make([]string, 0, len(deleteSkills))
		for _, s := range deleteSkills {
			skillsToDelete = append(skillsToDelete, s.Name)
		}

		if err := DeleteSkills(s, t, skillsToDelete); err != nil {
			return err
		}
	}

	return nil
}

// DeleteSkills deletes skill connection between post and skills
func DeleteSkills(s dbr.SessionRunner, t RequirementType, names []string) error {
	skills, err := database.ReadSkill(s, names)
	if err != nil {
		return err
	}

	if len(skills) != len(names) {
		return dbr.ErrNotFound
	}

	skillIDs := make([]int64, 0, len(skills))
	for _, s := range skills {
		skillIDs = append(skillIDs, s.ID)
	}

	switch t {
	case TaskRequirementType:
		return database.DeleteTaskSkillRequirements(s, skillIDs)
	case VacancyRequirementType:
		return database.DeleteVacancySkillRequirements(s, skillIDs)
	default:
		return ErrUnknownRequirementType
	}
}
