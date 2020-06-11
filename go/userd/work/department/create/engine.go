package create

import (
	"net/http"

	"github.com/GDVFox/tenjin/userd/database"
	db "github.com/GDVFox/tenjin/utils/database"
	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/GDVFox/tenjin/utils/server"
	"github.com/pkg/errors"
)

// Process creates department record in database
func Process(r *http.Request, sess *db.Session, a server.Arguments, logger *logging.Logger) (interface{}, error) {
	arg := a.(*Arguments)

	if err := database.CreateDepartment(sess, arg.DepartmentName); err != nil {
		logger.Errorf("can not create department: %s", err)
		return nil, errors.Wrap(err, "can not create department")
	}

	return map[string]string{"result": "OK"}, nil
}
