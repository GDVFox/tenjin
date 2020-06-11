package attachment

import (
	"database/sql"

	"github.com/GDVFox/tenjin/userd/args"
	"github.com/GDVFox/tenjin/userd/database"
	"github.com/gocraft/dbr/v2"
)

// CreateAttachments from request arguments
func CreateAttachments(s dbr.SessionRunner, postID, commentID int64, arg []*args.AttachmentCreate) error {
	attchs := make([]*database.AttachmentModel, 0, len(arg))
	for _, a := range arg {
		da := &database.AttachmentModel{URI: a.URI}

		if postID != 0 {
			da.PostID = dbr.NullInt64{NullInt64: sql.NullInt64{Int64: postID, Valid: true}}
		} else if commentID != 0 {
			da.CommentID = dbr.NullInt64{NullInt64: sql.NullInt64{Int64: commentID, Valid: true}}
		}

		attchs = append(attchs, da)
	}

	return database.CreateAttachments(s, attchs)
}

// UpdateAttachments updates attachments from request
func UpdateAttachments(s dbr.SessionRunner, postID, commentID int64, arg []*args.AttachmentUpdate) error {
	createAttachments := make([]*args.AttachmentUpdate, 0)
	deleteAttachments := make([]*args.AttachmentUpdate, 0)

	for _, s := range arg {
		switch s.Type {
		case args.AddUpdateType:
			createAttachments = append(createAttachments, s)
		case args.DeleteUpdateType:
			deleteAttachments = append(deleteAttachments, s)
		}
	}

	if len(createAttachments) > 0 {
		attchsToCreate := make([]*args.AttachmentCreate, 0, len(createAttachments))
		for _, a := range createAttachments {
			da := &args.AttachmentCreate{URI: *a.URI}
			attchsToCreate = append(attchsToCreate, da)
		}

		if err := CreateAttachments(s, postID, commentID, attchsToCreate); err != nil {
			return err
		}
	}

	if len(deleteAttachments) > 0 {
		attchs := make([]int64, 0, len(deleteAttachments))
		for _, a := range deleteAttachments {
			attchs = append(attchs, a.ID)
		}

		if err := database.DeleteAttachments(s, postID, commentID, attchs); err != nil {
			return err
		}
	}

	return nil
}
