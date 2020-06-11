package args

import (
	"net/http"

	"github.com/GDVFox/tenjin/utils/server"
)

// PostCreate common arguments for post creation
type PostCreate struct {
	Text        string              `json:"text"`
	AuthorID    int64               `json:"author_id"`
	Attachments []*AttachmentCreate `json:"attachments"`
}

// Validate checks argument correct
func (a *PostCreate) Validate() error {
	if a.Text == "" {
		return server.NewHTTPError(http.StatusBadRequest, "post text can not be empty")
	}

	if a.AuthorID == 0 {
		return server.NewHTTPError(http.StatusBadRequest, "post author_id can not be empty")
	}

	for _, a := range a.Attachments {
		if err := a.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// PostUpdate  common arguments for post update
type PostUpdate struct {
	ID          int64               `json:"id"`
	Text        *string             `json:"text"`
	Attachments []*AttachmentUpdate `json:"attachments"`
}

// Validate checks argument correct
func (a *PostUpdate) Validate() error {
	if a.Text != nil && *a.Text == "" {
		return server.NewHTTPError(http.StatusBadRequest, "title can not be empty")
	}

	for _, a := range a.Attachments {
		if err := a.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// IsAttachmentsUpdated return true if attachments has to be updated
func (a *PostUpdate) IsAttachmentsUpdated() bool {
	return len(a.Attachments) > 0
}

// IsTextUpdated return true if text entity has to be updated
func (a *PostUpdate) IsTextUpdated() bool {
	return a.Text != nil
}

// IsPostUpdated return true if post entity has to be updated
func (a *PostUpdate) IsPostUpdated() bool {
	return a.IsTextUpdated() || a.IsAttachmentsUpdated()
}
