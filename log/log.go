package log

import (
	"io"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Logger is an interface the defines logging methods.
type Logger interface {
	Errorf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Debugf(string, ...interface{})
	Fatalf(string, ...interface{})
}

// LoggerWithInterceptor is an interface that uses a logging interface
// and configures unary server interceptors for
// gRPC logging.
type LoggerWithInterceptor interface {
	Logger
	WithServerInterceptors() []grpc.UnaryServerInterceptor
}

// LogrusLogger is a wrapper for Logrus to implement the Logger interface.
type LogrusLogger struct {
	Logger *logrus.Logger
}

// NewLogrusLogger creates a new LogrusLogger with the given inputs.
func NewLogrusLogger(out io.Writer, level string) (*LogrusLogger, error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	slogr := &LogrusLogger{
		Logger: &logrus.Logger{
			Out:       out,
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     lvl,
		},
	}

	return slogr, nil
}

func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.Logger.Logf(logrus.ErrorLevel, format, args...)
}

func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.Logger.Logf(logrus.InfoLevel, format, args...)
}

func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.Logger.Logf(logrus.WarnLevel, format, args...)
}

func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.Logger.Logf(logrus.DebugLevel, format, args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.Logger.Logf(logrus.FatalLevel, format, args...)
}

func (l *LogrusLogger) WithServerInterceptors() []grpc.UnaryServerInterceptor {
	// Logrus entry is used, allowing pre-definition of certain fields by the user.
	logrusEntry := logrus.NewEntry(l.Logger)
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
	}
	// Make sure that log statements internal to gRPC library are logged using the logrus Logger as well.
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)
	return []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
	}
}
