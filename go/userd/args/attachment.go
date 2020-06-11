package args

import (
	"net/http"

	"github.com/GDVFox/tenjin/utils/server"
)

// AttachmentCreate common arguments for attachment creation
type AttachmentCreate struct {
	URI string `json:"uri"`
}

// Validate checks argument correct
func (a *AttachmentCreate) Validate() error {
	if a.URI == "" {
		return server.NewHTTPError(http.StatusBadRequest, "attachment uri can not be empty")
	}

	return nil
}

// AttachmentUpdate common arguments for attachment update
type AttachmentUpdate struct {
	Type AttachmentUpdateType `json:"type"`
	ID   int64                `json:"id"`
	URI  *string              `json:"uri"`
}

// Validate checks argument correct
func (a *AttachmentUpdate) Validate() error {
	if a.Type == AddUpdateType &&
		(a.URI == nil || *a.URI == "") {
		return server.NewHTTPError(http.StatusBadRequest, "uri can not be empty with type=add")
	}

	return nil
}
