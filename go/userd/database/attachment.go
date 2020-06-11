package database

import (
	"github.com/GDVFox/tenjin/utils/database"
	"github.com/gocraft/dbr/v2"
)

// AttachmentStatus is alias for status field in table
type AttachmentStatus int32

// AttachmentStatus is mapping of attachment.status column
const (
	AttachmentStatusActive AttachmentStatus = 1 + iota
	AttachmentStatusDeleted
)

// AttachmentModel represents post or comment attachment
type AttachmentModel struct {
	ID        int64         `db:"id"`
	URI       string        `db:"uri"`
	PostID    dbr.NullInt64 `db:"post_id"`
	CommentID dbr.NullInt64 `db:"comment_id"`
}

// CreateAttachments creates attachment rows
func CreateAttachments(s dbr.SessionRunner, attchs []*AttachmentModel) error {
	q := s.InsertInto("attachment").Columns("uri", "post_id", "comment_id")
	for _, a := range attchs {
		q.Record(a)
	}

	res, err := q.Exec()
	if err != nil {
		return err
	}

	return database.AssertAffected(res, int64(len(attchs)))
}

// ReadAttachments loads attachments for posts with postIDs and comments with commentIDs
func ReadAttachments(s dbr.SessionRunner, postIDs, commentIDs []int64) ([]*AttachmentModel, error) {
	q := s.Select(
		"id", "uri", "CAST(status AS unsigned) AS status",
		"post_id", "comment_id",
	).From("attachment")
	q.Where(dbr.Neq("status", AttachmentStatusDeleted))
	if len(postIDs) != 0 || len(commentIDs) != 0 {
		q.Where(dbr.Or(dbr.Eq("post_id", postIDs), dbr.Eq("comment_id", commentIDs)))
	}

	res := make([]*AttachmentModel, 0, len(postIDs)+len(commentIDs))
	_, err := q.Load(&res)
	return res, err
}

// DeleteAttachments remove attachments
func DeleteAttachments(s dbr.SessionRunner, postIDs, commentIDs int64, ids []int64) error {
	if len(ids) == 0 && postIDs == 0 && commentIDs == 0 {
		return ErrNoKeysSpecified
	}

	q := s.Update("attachment").
		Set("status", AttachmentStatusDeleted).
		Where(dbr.Neq("status", AttachmentStatusDeleted))
	if len(ids) != 0 {
		q.Where(dbr.Eq("id", ids))
	}
	if postIDs != 0 {
		q.Where(dbr.Eq("post_id", postIDs))
	}
	if commentIDs != 0 {
		q.Where(dbr.Eq("comment_id", commentIDs))
	}

	res, err := q.Exec()
	if err != nil {
		return err
	}

	return database.SomeAffected(res)
}
