package base

import (
	"github.com/GDVFox/tenjin/userd/database"
)

// AuthorItem represents short info about author
type AuthorItem struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// TaskReplyItem represents reply item
type TaskReplyItem struct {
	ID     int64      `json:"id"`
	Title  string     `json:"title"`
	Text   string     `json:"text"`
	Author AuthorItem `json:"author"`
	Skills []string   `json:"require_skills,omitempty"`
}

type replyBuilder struct {
	tasks           []*TaskReplyItem
	tasksSet        map[int64]struct{}
	tasksAuthorsSet map[int64]struct{}
	authorsSet      map[int64]*database.PersonModel
	taskSkills      map[int64][]string
}

func newReplyBulder() *replyBuilder {
	return &replyBuilder{
		tasks:           make([]*TaskReplyItem, 0),
		tasksSet:        make(map[int64]struct{}),
		tasksAuthorsSet: make(map[int64]struct{}),
		authorsSet:      make(map[int64]*database.PersonModel),
		taskSkills:      make(map[int64][]string),
	}
}

func (b *replyBuilder) consumeTasks(tt []*database.TaskModel) {
	for _, t := range tt {
		b.tasks = append(b.tasks, &TaskReplyItem{
			ID:    t.ID,
			Title: t.Title,
			Text:  t.Text,
			Author: AuthorItem{
				ID: t.EmployeeID,
			},
		})
		b.tasksSet[t.ID] = struct{}{}
		b.tasksAuthorsSet[t.EmployeeID] = struct{}{}
	}
}

func (b *replyBuilder) getTasks() []int64 {
	result := make([]int64, 0, len(b.tasksSet))
	for t := range b.tasksSet {
		result = append(result, t)
	}
	return result
}

func (b *replyBuilder) getAuthors() []int64 {
	result := make([]int64, 0, len(b.tasksAuthorsSet))
	for a := range b.tasksAuthorsSet {
		result = append(result, a)
	}
	return result
}

func (b *replyBuilder) consumeAuthors(pp []*database.PersonModel) {
	for _, p := range pp {
		b.authorsSet[p.ID] = p
	}
}

func (b *replyBuilder) consumeSkillRequirements(ss []*database.SkillRequirement) {
	for _, s := range ss {
		if _, ok := b.taskSkills[s.TaskID]; !ok {
			b.taskSkills[s.TaskID] = []string{s.SkillName}
			continue
		}

		b.taskSkills[s.TaskID] = append(b.taskSkills[s.TaskID], s.SkillName)
	}
}

func (b *replyBuilder) reply() []*TaskReplyItem {
	for _, t := range b.tasks {
		author := b.authorsSet[t.Author.ID]
		t.Author.FirstName = author.FirstName
		t.Author.LastName = author.LastName

		t.Skills = b.taskSkills[t.ID]
	}

	return b.tasks
}
