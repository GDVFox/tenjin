package update

import (
	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/utils/server"
)

// Arguments represents update args
type Arguments struct {
	args.PostUpdate
}

// Validate represents arguments for task update
func (a *Arguments) Validate() error {
	if a.PostUpdate.Text == nil && len(a.PostUpdate.Attachments) == 0 {
		return server.EmptyArgumentsHTTPError
	}

	return a.PostUpdate.Validate()
}

// NewArguments creates new update arguments
func NewArguments() server.Arguments {
	return &Arguments{}
}
