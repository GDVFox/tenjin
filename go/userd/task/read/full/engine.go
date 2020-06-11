package full

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Process reads full task record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	b := newReplyBulder()

	baseGroup, _ := errgroup.WithContext(r.Context())
	baseGroup.Go(func() error {
		tasks, err := database.ReadTaskFull(sess, arg.ID)
		if err != nil {
			return errors.Wrap(err, "can not read task")
		}
		b.consumeTasks(tasks)
		return nil
	})

	baseGroup.Go(func() error {
		solutions, err := database.ReadSolution(sess, []int64{arg.ID}, nil,
			false, database.ApprovedNewestSolutionOrder)
		if err != nil {
			return errors.Wrap(err, "can not read solutions")
		}
		b.consumeSolutions(solutions)
		return nil
	})

	baseGroup.Go(func() error {
		skills, err := database.ReadTaskSkillRequirements(sess, []int64{arg.ID})
		if err != nil {
			return errors.Wrap(err, "can not skill reqs")
		}
		b.consumeSkillRequirements(skills)
		return nil
	})

	if err := baseGroup.Wait(); err != nil {
		logger.Errorf("can not read task info: %s", err)
		return nil, err
	}

	comments, err := database.ReadComments(sess, 0, b.getPostIDs())
	if err != nil {
		logger.Errorf("can not read comment info: %s", err)
		return nil, errors.Wrap(err, "can not read comments")
	}
	b.consumeComments(comments)

	additionGroup, _ := errgroup.WithContext(r.Context())
	additionGroup.Go(func() error {
		authors, err := database.ReadPerson(sess, b.getAuthors(), true)
		if err != nil {
			return errors.Wrap(err, "can not read authors")
		}
		b.consumeAuthors(authors)
		return nil
	})

	additionGroup.Go(func() error {
		attchs, err := database.ReadAttachments(sess, b.getPostIDs(), b.getCommentsIDs())
		if err != nil {
			return errors.Wrap(err, "can not read attachments")
		}
		b.consumeAttachments(attchs)
		return nil
	})

	if err := additionGroup.Wait(); err != nil {
		logger.Errorf("can not read task additional info: %s", err)
		return nil, err
	}

	return b.reply(), nil
}
