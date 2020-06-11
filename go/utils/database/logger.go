package database

import (
	"time"

	"github.com/GDVFox/tenjin/utils/logging"
	"github.com/gocraft/dbr/v2"
)

// Logger implements dbr.EventReceiver interface
type Logger struct {
	*dbr.NullEventReceiver
	logger *logging.Logger
}

// TimingKv receives the time an event took to happen along with optional key/value data
func (l *Logger) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	l.logger.Infof("Query %s timing: %s", kvs["sql"], time.Duration(nanoseconds).String())
}

// EventErrKv receives a notification of an error if one occurs along with optional key/value data
func (l *Logger) EventErrKv(eventName string, err error, kvs map[string]string) error {
	l.logger.Errorf("Query %s failed: %s", kvs["sql"], err.Error())
	return err
}

// Timing receives the time an event took to happen
func (l *Logger) Timing(eventName string, nanoseconds int64) {
	l.logger.Infof("%s timing: %s", eventName, time.Duration(nanoseconds).String())
}

// EventErr receives a notification of an error if one occurs
func (l *Logger) EventErr(eventName string, err error) error {
	l.logger.Errorf("%s failed: %s", eventName, err)
	return err
}

// NewLogger creates a logger, which will log each query.
func NewLogger(logger *logging.Logger) *Logger {
	return &Logger{logger: logger}
}
