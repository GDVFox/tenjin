package full

import (
	"time"

	"github.com/GDVFox/tenjin/userd/database"
	"github.com/GDVFox/tenjin/utils/functions"
)

// AuthorItem represents short info about author
type AuthorItem struct {
	ID        int64  `json:"id"`
	PhotoURI  string `json:"photo_uri,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// AttachmentItem represents attachment
type AttachmentItem struct {
	ID  int64  `json:"id"`
	URI string `json:"uri"`
}

// CommentItem represents comment
type CommentItem struct {
	ID          int64             `json:"id"`
	Text        *string           `json:"text,omitempty"`
	CreatedAt   *time.Time        `json:"created_at,omitempty"`
	UpdatedAt   *time.Time        `json:"updated_at,omitempty"`
	Deleted     bool              `json:"deleted"`
	Author      *AuthorItem       `json:"author,omitempty"`
	Childs      []*CommentItem    `json:"replies,omitempty"`
	Attachments []*AttachmentItem `json:"attachments,omitempty"`
}

// SolutionItem represents soltion info for task
type SolutionItem struct {
	ID          int64             `json:"id"`
	Text        string            `json:"text"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	IsApproved  bool              `json:"is_approved"`
	Author      *AuthorItem       `json:"author"`
	Comments    []*CommentItem    `json:"comments,omitempty"`
	Attachments []*AttachmentItem `json:"attachments,omitempty"`
}

// SkillItem represents skill requirement
type SkillItem struct {
	Name       string                   `json:"name"`
	Difficulty database.SkillDifficulty `json:"difficulty"`
}

// TaskReplyItem represents reply item
type TaskReplyItem struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Text        string            `json:"text"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Author      *AuthorItem       `json:"author"`
	Solutions   []*SolutionItem   `json:"solutions,omitempty"`
	Skills      []*SkillItem      `json:"require_skills,omitempty"`
	Comments    []*CommentItem    `json:"comments,omitempty"`
	Attachments []*AttachmentItem `json:"attachments,omitempty"`
}

type replyBuilder struct {
	task *TaskReplyItem

	solutions []*SolutionItem

	skills  []*SkillItem
	postIDs map[int64]struct{}

	comments               []*CommentItem
	rootComments           map[int64][]int
	parentIDToChildIndexes map[int64][]int

	authors     []*AuthorItem
	authorIndex map[int64]int

	attachmentsPosts    map[int64][]*AttachmentItem
	attachmentsComments map[int64][]*AttachmentItem
}

func newReplyBulder() *replyBuilder {
	return &replyBuilder{
		authorIndex: make(map[int64]int),
		postIDs:     make(map[int64]struct{}),
	}
}

func (b *replyBuilder) getPostIDs() []int64 {
	return functions.SetToSliceInt64(b.postIDs)
}

func (b *replyBuilder) getAuthors() []int64 {
	result := make([]int64, 0, len(b.authorIndex))
	for id := range b.authorIndex {
		result = append(result, id)
	}

	return result
}

func (b *replyBuilder) consumeTasks(t *database.TaskModel) {
	b.task = &TaskReplyItem{
		ID:        t.ID,
		Title:     t.Title,
		Text:      t.Text,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		Author: &AuthorItem{
			ID: t.EmployeeID,
		},
	}
	b.authorIndex[t.EmployeeID] = -1
	b.postIDs[t.ID] = struct{}{}
}

func (b *replyBuilder) consumeSolutions(ss []*database.SolutionModel) {
	b.solutions = make([]*SolutionItem, 0, len(ss))
	for _, s := range ss {
		b.solutions = append(b.solutions, &SolutionItem{
			ID:         s.ID,
			Text:       s.Text,
			CreatedAt:  s.CreatedAt,
			UpdatedAt:  s.UpdatedAt,
			IsApproved: s.IsApproved,
			Author: &AuthorItem{
				ID: s.EmployeeID,
			},
		})
		b.postIDs[s.ID] = struct{}{}
		b.authorIndex[s.EmployeeID] = -1
	}
}

func (b *replyBuilder) consumeSkillRequirements(ss []*database.SkillRequirement) {
	b.skills = make([]*SkillItem, 0, len(ss))
	for _, s := range ss {
		b.skills = append(b.skills, &SkillItem{
			Name:       s.SkillName,
			Difficulty: *s.Difficulty,
		})
	}
}

func (b *replyBuilder) consumeComments(cc []*database.CommentModel) {
	b.comments = make([]*CommentItem, 0, len(cc))
	b.rootComments = make(map[int64][]int)
	b.parentIDToChildIndexes = make(map[int64][]int)

	rootSet := make(map[int64]struct{})
	for i, c := range cc {
		if c.Status == database.CommentStatusDeleted {
			b.comments = append(b.comments, &CommentItem{
				ID:      c.ID,
				Deleted: true,
			})
		} else {
			b.comments = append(b.comments, &CommentItem{
				ID:        c.ID,
				Text:      &c.Text,
				CreatedAt: &c.CreatedAt,
				UpdatedAt: &c.UpdatedAt,
				Deleted:   false,
				Author: &AuthorItem{
					ID: c.EmployeeID,
				},
			})
			b.authorIndex[c.EmployeeID] = -1
		}

		if c.Parent.Valid {
			b.parentIDToChildIndexes[c.Parent.Int64] =
				append(b.parentIDToChildIndexes[c.Parent.Int64], i)
		} else {
			b.rootComments[c.PostID] = append(b.rootComments[c.PostID], i)
			rootSet[c.ID] = struct{}{}
		}
	}
}

func (b *replyBuilder) consumeAttachments(aa []*database.AttachmentModel) {
	b.attachmentsPosts = make(map[int64][]*AttachmentItem, len(aa))
	b.attachmentsComments = make(map[int64][]*AttachmentItem, len(aa))

	for _, a := range aa {
		attch := &AttachmentItem{
			ID:  a.ID,
			URI: a.URI,
		}

		if a.PostID.Valid {
			b.attachmentsPosts[a.PostID.Int64] =
				append(b.attachmentsPosts[a.PostID.Int64], attch)
		} else if a.CommentID.Valid {
			b.attachmentsComments[a.CommentID.Int64] =
				append(b.attachmentsComments[a.CommentID.Int64], attch)
		}
	}
}

func (b *replyBuilder) getCommentsIDs() []int64 {
	result := make([]int64, 0, len(b.comments))
	for _, c := range b.comments {
		if !c.Deleted {
			result = append(result, c.ID)
		}
	}

	return result
}

func (b *replyBuilder) consumeAuthors(pp []*database.PersonModel) {
	for i, p := range pp {
		b.authorIndex[p.ID] = i
	}

	b.task.Author.PhotoURI = pp[b.authorIndex[b.task.Author.ID]].PhotoURI.String
	b.task.Author.FirstName = pp[b.authorIndex[b.task.Author.ID]].FirstName
	b.task.Author.LastName = pp[b.authorIndex[b.task.Author.ID]].LastName

	for _, s := range b.solutions {
		s.Author.PhotoURI = pp[b.authorIndex[s.Author.ID]].PhotoURI.String
		s.Author.FirstName = pp[b.authorIndex[s.Author.ID]].FirstName
		s.Author.LastName = pp[b.authorIndex[s.Author.ID]].LastName
	}

	for _, c := range b.comments {
		if c.Author == nil {
			continue
		}
		c.Author.PhotoURI = pp[b.authorIndex[c.Author.ID]].PhotoURI.String
		c.Author.FirstName = pp[b.authorIndex[c.Author.ID]].FirstName
		c.Author.LastName = pp[b.authorIndex[c.Author.ID]].LastName
	}
}

func (b *replyBuilder) traverseComments(root *CommentItem) {
	childs := b.parentIDToChildIndexes[root.ID]
	if len(childs) == 0 {
		return
	}

	root.Childs = make([]*CommentItem, 0, len(childs))
	for _, indx := range childs {
		nextRoot := b.comments[indx]
		nextRoot.Attachments = b.attachmentsComments[nextRoot.ID]

		b.traverseComments(nextRoot)
		root.Childs = append(root.Childs, nextRoot)
	}
}

func (b *replyBuilder) buildRootComments(postID int64) []*CommentItem {
	roots, ok := b.rootComments[postID]
	if !ok {
		return nil
	}

	result := make([]*CommentItem, 0, len(roots))
	for _, indx := range roots {
		root := b.comments[indx]
		root.Attachments = b.attachmentsComments[root.ID]

		b.traverseComments(root)
		result = append(result, root)
	}

	return result
}

func (b *replyBuilder) reply() *TaskReplyItem {
	b.task.Comments = b.buildRootComments(b.task.ID)
	b.task.Attachments = b.attachmentsPosts[b.task.ID]

	for _, s := range b.solutions {
		s.Comments = b.buildRootComments(s.ID)
		s.Attachments = b.attachmentsPosts[s.ID]
	}

	b.task.Solutions = b.solutions
	b.task.Skills = b.skills
	return b.task
}
