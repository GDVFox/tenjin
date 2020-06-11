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

// Process reads full information about person
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	b := newReplyBulder()
	g, _ := errgroup.WithContext(r.Context())

	g.Go(func() error {
		employee, err := database.ReadPersonFull(sess, []int64{arg.ID})
		if err != nil {
			return errors.Wrap(err, "can not read full person")
		}
		if len(employee) == 0 {
			return server.NotFoundHTTPError
		}

		b.consumePerson(employee[0])
		return nil
	})

	g.Go(func() error {
		history, err := database.ReadWorkHistory(sess, []int64{arg.ID})
		if err != nil {
			return errors.Wrap(err, "can not read work history")
		}

		b.consumeWorkHistory(history)
		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Errorf("can not read persons: %s", err)
		return nil, err
	}

	return b.reply(), nil
}
