package create

import (
	"net/http"

	"github.com/GDVFox/tenjin/utils/server"
)

// Arguments represents arguments for vote creation
type Arguments struct {
	PostID     *int64 `json:"post_id"`
	CommentID  *int64 `json:"comment_id"`
	EmployeeID int64  `json:"employee_id"`
	Delta      int    `json:"delta"`
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	if a.PostID == nil && a.CommentID == nil {
		return server.EmptyArgumentsHTTPError
	}

	if a.PostID != nil && *a.PostID == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "post_id can not be 0")
	}

	if a.CommentID != nil && *a.CommentID == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "comment_id can not be 0")
	}

	if a.EmployeeID == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "employee_id can not be 0")
	}

	if a.Delta != -1 && a.Delta != 1 {
		return server.NewHTTPError(http.StatusBadRequest, "delta must be (-1 or 1)")
	}

	return nil
}

// NewArguments creates new vote args instanse
func NewArguments() server.Arguments {
	return &Arguments{}
}
