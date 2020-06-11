package create

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/utils/server"
)

// Arguments for department requests
type Arguments struct {
	args.PostCreate
	PostID int64  `json:"post_id"`
	Parent *int64 `json:"parent"`
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	if a.PostID == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "post_id can not be 0")
	}

	if a.Parent != nil && *a.Parent == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "parent can not be 0")
	}

	return a.PostCreate.Validate()
}

// NewArguments creates new CreatePersonArguments instanse
func NewArguments() server.Arguments {
	return &Arguments{}
}
