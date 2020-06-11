package create

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/utils/server"
)

// Arguments represents arguments for person creation
type Arguments struct {
	args.PostCreate
	TaskID int64 `json:"task_id"`
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	if a.TaskID == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "task_id can not be 0")
	}

	return a.PostCreate.Validate()
}

// NewArguments creates new CreatePersonArguments instanse
func NewArguments() server.Arguments {
	return &Arguments{}
}
