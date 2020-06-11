package requirement

import (
	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

// CreateRequirementChecks creates requirement_check rows in database
func CreateRequirementChecks(s dbr.SessionRunner, inverviewID, vacancyID int64, arg []*args.RequirementCheckCreateArgument) error {
	params := make([]*args.RequirementCheckUpdateArgument, 0, len(arg))
	for _, a := range arg {
		p := &args.RequirementCheckUpdateArgument{
			Type:      args.AddUpdateType,
			TaskID:    a.TaskID,
			SkillName: a.SkillName,
			Comment:   a.Comment,
			Score:     &a.Score,
		}

		params = append(params, p)
	}

	return UpdateRequirementChecks(s, inverviewID, vacancyID, params)
}

// UpdateRequirementChecks creates and deletes requirement_check rows in database
func UpdateRequirementChecks(s dbr.SessionRunner, inverviewID, vacancyID int64, arg []*args.RequirementCheckUpdateArgument) error {
	skillNames := make([]string, 0, len(arg))
	for _, s := range arg {
		skillNames = append(skillNames, s.SkillName)
	}

	skills, err := database.ReadSkill(s, skillNames)
	if err != nil {
		return err
	}

	// short path
	if len(skills) != len(skillNames) {
		return dbr.ErrNotFound
	}

	knownSkills := make(map[string]int64, len(skills))
	for _, s := range skills {
		knownSkills[s.Name] = s.ID
	}

	createReqs := make([]*database.RequirementCheckModel, 0)
	deleteReqs := make([]*database.RequirementCheckModel, 0)
	for _, s := range arg {
		req := &database.RequirementCheckModel{
			VacancyID:   vacancyID,
			InterviewID: inverviewID,
			TaskID:      s.TaskID,
			SkillID:     knownSkills[s.SkillName],
		}
		if s.Type == args.AddUpdateType {
			req.Comment = db.NewNullString(s.Comment)
			req.Score = *s.Score
			createReqs = append(createReqs, req)
		} else {
			deleteReqs = append(deleteReqs, req)
		}
	}

	if len(createReqs) != 0 {
		if err := database.CreateRequirementCheck(s, createReqs); err != nil {
			return err
		}
	}

	if len(deleteReqs) != 0 {
		if err := database.DeleteRequirementCheck(s, deleteReqs); err != nil {
			return err
		}
	}

	return nil
}
