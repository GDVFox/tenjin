package update

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/utils/server"
)

// Arguments represents update args
type Arguments struct {
	args.PostUpdate
	Title  *string             `json:"text"`
	Skills []*args.SkillUpdate `json:"skills"`
}

// Validate represents arguments for task update
func (a *Arguments) Validate() error {
	if a.Title == nil && a.PostUpdate.Text == nil &&
		len(a.Skills) == 0 && len(a.Attachments) == 0 {
		return server.EmptyArgumentsHTTPError
	}

	if a.Title != nil && *a.Title == "" {
		return server.NewHTTPError(http.StatusBadRequest, "title can not be empty")
	}

	for _, s := range a.Skills {
		if err := s.Validate(); err != nil {
			return err
		}
	}

	return a.PostUpdate.Validate()
}

// IsTaskUpdated return true if task entity has to be updated
func (a *Arguments) IsTaskUpdated() bool {
	return a.Title != nil
}

// IsSkillUpdated return true if some skills has to be updated
func (a *Arguments) IsSkillUpdated() bool {
	return len(a.Skills) != 0
}

// NewArguments creates new update arguments
func NewArguments() server.Arguments {
	return &Arguments{}
}
