package database

import (
	"database/sql"
	"errors"

	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/gocraft/dbr/v2"
)

var (
	// ErrNoRowsAffected request done nothing
	ErrNoRowsAffected = errors.New("no rows affected")
)

const (
	driver = "mysql"
)

var (
	connection *dbr.Connection
)

//Session represents database session
type Session struct {
	*dbr.Session
}

// Open opens new database connection against cfg
func Open(cfg *Config) error {
	var err error
	connection, err = dbr.Open(driver, cfg.DSN(), nil)
	return err
}

// NewSession creates a new sessions
func NewSession(logger *logging.Logger) *Session {
	return &Session{Session: connection.NewSession(NewLogger(logger))}
}

// AssertAffected retuns ErrNoRowsAffected if res.affected != expected
func AssertAffected(res sql.Result, expected int64) error {
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != expected {
		return ErrNoRowsAffected
	}
	return nil
}

// SomeAffected retuns ErrNoRowsAffected if res.affected == 0
func SomeAffected(res sql.Result) error {
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNoRowsAffected
	}
	return nil
}
