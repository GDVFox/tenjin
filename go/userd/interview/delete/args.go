package delete

import (
	"net/http"

	"github.com/GDVFox/tenjin/utils/server"
)

// Arguments represents update args
type Arguments struct {
	ID int64 `json:"id"`
}

// Validate represents arguments for task update
func (a *Arguments) Validate() error {
	if a.ID == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "id can not be 0")
	}

	return nil
}

// NewArguments creates new update arguments
func NewArguments() server.Arguments {
	return &Arguments{}
}
