// Package log is a structured logger for Go, based on https://github.com/uber-go/zap.
package log

import "go.uber.org/zap"

var (
	defaultLogger *Logger
	// EncodedFilename filename for logging when DisableFile is false.
	EncodedFilename string
)

func init() {
	defaultLogger = New(NewOptions())
}

type DebugLogger interface {
	Debugt(msg string, fields ...Field)
	Debugf(template string, args ...interface{})
	Debug(msg string, keysAndValues ...interface{})
}

type InfoLogger interface {
	DebugLogger

	Infot(msg string, fields ...Field)
	Infof(template string, args ...interface{})
	Info(msg string, keysAndValues ...interface{})
}

type Interface interface {
	InfoLogger

	Warnt(msg string, fields ...Field)
	Warnf(template string, args ...interface{})
	Warn(msg string, keysAndValues ...interface{})

	Errort(msg string, fields ...Field)
	Errorf(template string, args ...interface{})
	Error(msg string, keysAndValues ...interface{})

	Panict(msg string, fields ...Field)
	Panicf(template string, args ...interface{})
	Panic(msg string, keysAndValues ...interface{})

	Fatalt(msg string, fields ...Field)
	Fatalf(template string, args ...interface{})
	Fatal(msg string, keysAndValues ...interface{})

	AtLevelt(level Level, msg string, fields ...Field)
	AtLevelf(level Level, template string, args ...interface{})
	AtLevel(level Level, msg string, keysAndValues ...interface{})

	// WithValues creates a child logger and adds some Field of
	// context to this logger.
	WithValues(fields ...Field) *Logger

	// Flush calls the underlying Core's Sync method, flushing any buffered
	// log entries. Applications should take care to call Sync before exiting.
	Flush() error

	// Close implements io.Closer, and closes the current logfile.
	Close() error
}

// Configure sets up the logging framework.
func Configure(opts *Options) {
	l := New(opts)
	defaultLogger = l
	zap.RedirectStdLog(defaultLogger.log)
}

// Debugt logs a message at DebugLevel.
func Debugt(msg string, fields ...Field) {
	defaultLogger.Debugt(msg, fields...)
}

// Debugf logs a message at DebugLevel.
func Debugf(template string, args ...interface{}) {
	defaultLogger.Debugf(template, args...)
}

// Debug logs a message at DebugLevel.
func Debug(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debug(msg, keysAndValues...)
}

// Infot logs a message at InfoLevel.
func Infot(msg string, fields ...Field) {
	defaultLogger.Infot(msg, fields...)
}

// Infof logs a message at InfoLevel.
func Infof(template string, args ...interface{}) {
	defaultLogger.Infof(template, args...)
}

// Info logs a message at InfoLevel.
func Info(msg string, keysAndValues ...interface{}) {
	defaultLogger.Info(msg, keysAndValues...)
}

// Warnt logs a message at WarnLevel.
func Warnt(msg string, fields ...Field) {
	defaultLogger.Warnt(msg, fields...)
}

// Warnf logs a message at WarnLevel.
func Warnf(template string, args ...interface{}) {
	defaultLogger.Warnf(template, args...)
}

// Warn logs a message at WarnLevel.
func Warn(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warn(msg, keysAndValues...)
}

// Errort logs a message at ErrorLevel.
func Errort(msg string, fields ...Field) {
	defaultLogger.Errort(msg, fields...)
}

// Errorf logs a message at ErrorLevel.
func Errorf(template string, args ...interface{}) {
	defaultLogger.Errorf(template, args...)
}

// Error logs a message at ErrorLevel.
func Error(msg string, keysAndValues ...interface{}) {
	defaultLogger.Error(msg, keysAndValues...)
}

// Panict logs a message at PanicLevel, then panics.
func Panict(msg string, fields ...Field) {
	defaultLogger.Panict(msg, fields...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	defaultLogger.Panicf(template, args...)
}

// Panic logs a message with some additional context, then panics.
func Panic(msg string, keysAndValues ...interface{}) {
	defaultLogger.Panic(msg, keysAndValues...)
}

// Fatalt logs a message at FatalLevel, then calls os.Exit(1).
func Fatalt(msg string, fields ...Field) {
	defaultLogger.Fatalt(msg, fields...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	defaultLogger.Fatalf(template, args...)
}

// Fatal logs a message with some additional context, then calls os.Exit.
func Fatal(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatal(msg, keysAndValues...)
}

// AtLevelt logs a message at Level.
func AtLevelt(level Level, msg string, fields ...Field) {
	defaultLogger.AtLevelt(level, msg, fields...)
}

// AtLevelf logs a message at Level.
func AtLevelf(level Level, msg string, args ...interface{}) {
	defaultLogger.AtLevelf(level, msg, args...)
}

// AtLevel logs a message at Level.
func AtLevel(level Level, msg string, keysAndValues ...interface{}) {
	defaultLogger.AtLevel(level, msg, keysAndValues...)
}

// WithValues creates a child logger and adds some Field of
// context to this logger.
func WithValues(fields ...Field) *Logger {
	return defaultLogger.WithValues(fields...)
}

// Flush calls the underlying Core's Sync method, flushing any buffered
// log entries. Applications should take care to call Sync before exiting.
func Flush() error { return defaultLogger.Flush() }

// Close implements io.Closer, and closes the current logfile of default logger.
func Close() error { return defaultLogger.Close() }
