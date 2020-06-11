package create

import (
	"net/http"

	"github.com/GDVFox/tenjin/utils/server"
)

// Arguments for department requests
type Arguments struct {
	DepartmentName string `json:"name"`
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	if a.DepartmentName == "" {
		return server.NewHTTPError(http.StatusBadRequest, "name can not be empty")
	}

	return nil
}

// NewArguments creates new CreatePersonArguments instanse
func NewArguments() server.Arguments {
	return &Arguments{}
}
