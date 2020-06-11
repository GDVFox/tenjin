package args

import (
	"net/http"

	"github.com/GDVFox/tenjin/utils/server"
)

// SkillCreate represents arguments for skill creation
type SkillCreate struct {
	Name       string                 `json:"name"`
	Difficulty *SkillDifficultyFilter `json:"difficulty"`
}

// Validate checks argument correct
func (a *SkillCreate) Validate() error {
	if a.Name == "" {
		return server.NewHTTPError(http.StatusBadRequest, "skill name can not be empty")
	}

	return nil
}

// SkillUpdate represents skill update args
type SkillUpdate struct {
	Type       SkillUpdateType        `json:"type"`
	Name       string                 `json:"name"`
	Difficulty *SkillDifficultyFilter `json:"difficulty"`
}

// Validate checks argument correct
func (a *SkillUpdate) Validate() error {
	if a.Name == "" {
		return server.BadRequestHTTPError
	}

	return nil
}
