package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Interface interface {
	Debugt(msg string, fields ...zapcore.Field)
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Debugs(args ...interface{})

	Infot(msg string, fields ...zapcore.Field)
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Infos(args ...interface{})

	Warnt(msg string, fields ...zapcore.Field)
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Warns(args ...interface{})

	Errort(msg string, fields ...zapcore.Field)
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Errors(args ...interface{})

	Panict(msg string, fields ...zapcore.Field)
	Panicf(template string, args ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Panic(msg string, keysAndValues ...interface{})
	Panics(args ...interface{})

	Fatalt(msg string, fields ...zapcore.Field)
	Fatalf(template string, args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	Fatals(args ...interface{})

	AtLevel(level zapcore.Level, msg string, fields ...zapcore.Field) *Logger
}