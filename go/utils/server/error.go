package server

import (
	"fmt"
	"net/http"

	"github.com/GDVFox/tenjin/utils/database"
	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
)

// HTTPError represents http error
type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP [%d]: %s", e.Code, e.Message)
}

// NewHTTPError creates new http error
func NewHTTPError(code int, msg string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: msg,
	}
}

// WrapHTTPError creates new HTTPError from another
func WrapHTTPError(err *HTTPError, msg string) *HTTPError {
	return &HTTPError{
		Code:    err.Code,
		Message: err.Message + ": " + msg,
	}
}

var (
	// UnknownHTTPError represents internal unknown error
	UnknownHTTPError = NewHTTPError(http.StatusInternalServerError, "Something went wrong!")
	// BadRequestHTTPError represents bad arguments
	BadRequestHTTPError = NewHTTPError(http.StatusBadRequest, "Wrong arguments!")
	// DuplicationHTTPError represents that this entity already exisits
	DuplicationHTTPError = NewHTTPError(http.StatusConflict, "This entity already exists!")
	// NotFoundHTTPError represents not found error
	NotFoundHTTPError = NewHTTPError(http.StatusNotFound, "Entity not found!")
	// NoRowsAffectedHTTPError represents not found error
	NoRowsAffectedHTTPError = NewHTTPError(http.StatusNotFound, "Nothing affected!")
	// EmptyArgumentsHTTPError represent error when request is empty
	EmptyArgumentsHTTPError = NewHTTPError(http.StatusBadRequest, "Empty arguments!")
)

// WriteError write error in response
func WriteError(w http.ResponseWriter, err error) error {
	if httpErr, ok := err.(*HTTPError); ok {
		return WriteApplicationJSON(w, httpErr.Code, err)
	}

	httpErr := errorToHTTP(err)
	return WriteApplicationJSON(w, httpErr.Code, httpErr)
}

func errorToHTTP(err error) *HTTPError {
	if err == dbr.ErrNotFound {
		return NotFoundHTTPError
	}

	if err == database.ErrNoRowsAffected {
		return NoRowsAffectedHTTPError
	}

	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		// unique key check failed
		if mysqlError.Number == 1062 {
			return DuplicationHTTPError
		}

		// check failed or foreign key check failed or
		// custom trigger check failed
		if mysqlError.Number == 3819 || mysqlError.Number == 1452 ||
			mysqlError.Number == 1644 {
			return BadRequestHTTPError
		}
	}

	return UnknownHTTPError
}
