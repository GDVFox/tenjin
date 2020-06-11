package args

import (
	"net/http"

	"github.com/GDVFox/tenjin/utils/server"
)

// RequirementCheckCreateArgument represents create for requirement
type RequirementCheckCreateArgument struct {
	TaskID    int64   `json:"task_id"`
	SkillName string  `json:"skill_name"`
	Comment   *string `json:"comment"`
	Score     int64   `json:"score"`
}

// Validate represents arguments for check update
func (a *RequirementCheckCreateArgument) Validate() error {
	if a.TaskID == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "task_id can not be 0")
	}

	if a.SkillName == "" {
		return server.NewHTTPError(http.StatusBadRequest, "skill_name can not be empty")
	}

	if a.Comment != nil && *a.Comment == "" {
		return server.NewHTTPError(http.StatusBadRequest, "comment can not be empty")
	}

	return nil
}

// RequirementCheckUpdateArgument represents update for requirement
type RequirementCheckUpdateArgument struct {
	Type      RequirementCheckUpdateType `json:"type"`
	TaskID    int64                      `json:"task_id"`
	SkillName string                     `json:"skill_name"`
	Comment   *string                    `json:"comment"`
	Score     *int64                     `json:"score"`
}

// Validate represents arguments for check update
func (a *RequirementCheckUpdateArgument) Validate() error {
	if a.TaskID == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "task_id can not be 0")
	}

	if a.SkillName == "" {
		return server.NewHTTPError(http.StatusBadRequest, "skill_name can not be empty")
	}

	if a.Type == AddUpdateType && a.Score == nil {
		return server.NewHTTPError(http.StatusBadRequest, "can not add without score")
	}

	if a.Comment != nil && *a.Comment == "" {
		return server.NewHTTPError(http.StatusBadRequest, "comment can not be empty")
	}

	return nil
}
