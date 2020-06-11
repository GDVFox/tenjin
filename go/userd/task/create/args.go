package create

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/utils/server"
)

// Arguments represents arguments for task creation
type Arguments struct {
	args.PostCreate
	Title  string              `json:"title"`
	Skills []*args.SkillCreate `json:"skills"`
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	if a.Title == "" {
		return server.NewHTTPError(http.StatusBadRequest, "title can not be empty")
	}

	for _, s := range a.Skills {
		if err := s.Validate(); err != nil {
			return err
		}
	}

	return a.PostCreate.Validate()
}

// NewArguments creates new args instanse
func NewArguments() server.Arguments {
	return &Arguments{}
}
