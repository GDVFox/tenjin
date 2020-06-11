package base

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Process reads person record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	b := newReplyBulder()
	tasks, err := database.ReadTask(sess, arg.ID, arg.EmployeeID)
	if err != nil {
		logger.Errorf("can not read tasks: %s", err)
		return nil, err
	}
	b.consumeTasks(tasks)

	additionalGroup, _ := errgroup.WithContext(r.Context())
	additionalGroup.Go(func() error {
		persons, err := database.ReadPerson(sess, b.getAuthors(), true)
		if err != nil {
			return errors.Wrap(err, "can not read persons")
		}
		b.consumeAuthors(persons)
		return nil
	})
	additionalGroup.Go(func() error {
		skills, err := database.ReadTaskSkillRequirements(sess, b.getTasks())
		if err != nil {
			return errors.Wrap(err, "can not read skill req")
		}
		b.consumeSkillRequirements(skills)
		return nil
	})
	if err := additionalGroup.Wait(); err != nil {
		logger.Errorf("can not read task list: %s", err)
		return nil, err
	}

	return b.reply(), nil
}
