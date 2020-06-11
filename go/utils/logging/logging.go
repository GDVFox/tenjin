package logging

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger struct for logging
type Logger struct {
	*zap.SugaredLogger
}

// NewFileLogger creates a new file logger
func NewFileLogger(cfg *Config) (*Logger, error) {
	lvl := zap.NewAtomicLevel()
	if err := lvl.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, errors.Wrap(err, "can not set logging level")
	}

	logger, err := zap.Config{
		Encoding:    "json",
		Level:       lvl,
		OutputPaths: []string{cfg.File},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}.Build()
	if err != nil {
		return nil, errors.Wrap(err, "can init logger")
	}

	return &Logger{SugaredLogger: logger.Sugar()}, nil
}
