package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Logger is an interface the defines logging methods.
type Logger interface {
	Errorf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Debugf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type standardLogger struct {
	logger *logrus.Logger
}

// NewLogger returns a new Logger.
func NewLogger(out io.Writer, level string) (Logger, error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	slogr := &standardLogger{
		logger: &logrus.Logger{
			Out:       out,
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     lvl,
		},
	}

	return slogr, nil
}

func (l *standardLogger) Errorf(format string, args ...interface{}) {
	l.logger.Logf(logrus.ErrorLevel, format, args...)
}

func (l *standardLogger) Infof(format string, args ...interface{}) {
	l.logger.Logf(logrus.InfoLevel, format, args...)
}

func (l *standardLogger) Warnf(format string, args ...interface{}) {
	l.logger.Logf(logrus.WarnLevel, format, args...)
}

func (l *standardLogger) Debugf(format string, args ...interface{}) {
	l.logger.Logf(logrus.DebugLevel, format, args...)
}

func (l *standardLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Logf(logrus.FatalLevel, format, args...)
}
