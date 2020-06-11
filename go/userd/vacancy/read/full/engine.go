package full

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/gocraft/dbr/v2"
	"golang.org/x/sync/errgroup"
)

// Process reads person record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	b := newReplyBuilder()

	g, _ := errgroup.WithContext(r.Context())
	g.Go(func() error {
		vacancies, err := database.ReadVacancy(sess, []int64{arg.ID}, nil, nil)
		if err != nil {
			return err
		}
		if len(vacancies) == 0 {
			return dbr.ErrNotFound
		}
		b.consumeVacancy(vacancies[0])
		return nil
	})
	g.Go(func() error {
		skills, err := database.ReadVacancyRequirements(sess, []int64{arg.ID})
		if err != nil {
			return nil
		}
		b.consumeSkillRequirements(skills)
		return nil
	})
	g.Go(func() error {
		interview, err := database.ReadInterviewByVacancy(sess, arg.ID)
		if err != nil {
			return err
		}
		b.consumeInterviews(interview)
		return nil
	})
	if err := g.Wait(); err != nil {
		logger.Errorf("can not load vacancy: %s", err)
		return nil, err
	}

	return b.reply(), nil
}
