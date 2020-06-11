package read

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

	b := newReplyBuilder()

	g, _ := errgroup.WithContext(r.Context())
	g.Go(func() error {
		vacancies, err := database.ReadInterview(sess, arg.ID)
		if err != nil {
			return errors.Wrap(err, "can not read interview")
		}
		b.consumeInterview(vacancies[0])
		return nil
	})
	g.Go(func() error {
		checks, err := database.ReadRequirementCheck(sess, arg.ID)
		if err != nil {
			return errors.Wrap(err, "can not read checks")
		}
		b.consumeChecks(checks)
		return nil
	})
	if err := g.Wait(); err != nil {
		logger.Errorf("can not load interview: %s", err)
		return nil, err
	}

	return b.reply(), nil
}
