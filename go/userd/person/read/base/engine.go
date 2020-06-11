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

	b := newReplyBulder(len(arg.ID))
	g, _ := errgroup.WithContext(r.Context())
	g.Go(func() error {
		persons, err := database.ReadPerson(sess, arg.ID, arg.EmployeesOnly)
		if err != nil {
			return errors.Wrap(err, "can not read person")
		}

		b.consumePersons(persons)
		return nil
	})

	if arg.EmployeesOnly {
		g.Go(func() error {
			history, err := database.ReadWorkHistory(sess, arg.ID)
			if err != nil {
				return errors.Wrap(err, "can not read work history")
			}

			b.consumeWorkHistory(history)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		logger.Errorf("can not read persons: %s", err)
		return nil, err
	}

	return b.reply(), nil
}
