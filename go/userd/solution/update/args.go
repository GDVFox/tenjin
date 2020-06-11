package update

import (
	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/utils/server"
)

// Arguments represents update args
type Arguments struct {
	args.PostUpdate
	IsApproved *bool `json:"is_approved"`
}

// Validate represents arguments for task update
func (a *Arguments) Validate() error {
	if a.IsApproved == nil && a.PostUpdate.Text == nil && len(a.Attachments) == 0 {
		return server.BadRequestHTTPError
	}

	return a.PostUpdate.Validate()
}

// IsSolutionUpdated return true if solution entity has to be updated
func (a *Arguments) IsSolutionUpdated() bool {
	return a.IsApproved != nil
}

// NewArguments creates new update arguments
func NewArguments() server.Arguments {
	return &Arguments{}
}
