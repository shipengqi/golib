package log

import (
	"io"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	closer  io.Closer
	log     *zap.Logger
	sugared *zap.SugaredLogger
}

// New creates a new Logger.
func New(opts *Options) *Logger {
	if errs := opts.Validate(); len(errs) > 0 {
		panic(errs)
	}
	var cores []zapcore.Core

	encoderConfig := zapcore.EncoderConfig{
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if opts.TimeEncoder != nil {
		encoderConfig.EncodeTime = opts.TimeEncoder
	}
	if !opts.DisableConsole {
		var consoleLevel Level
		err := consoleLevel.Set(strings.ToLower(opts.ConsoleLevel))
		if err != nil {
			consoleLevel = InfoLevel
		}
		consoleEncCfg := encoderConfig
		if !opts.DisableConsoleLevel {
			consoleEncCfg.LevelKey = "level"
		}
		if !opts.DisableConsoleTime {
			consoleEncCfg.TimeKey = "time"
		}
		if !opts.DisableConsoleCaller {
			consoleEncCfg.CallerKey = "caller"
		}
		if !opts.DisableConsoleColor {
			consoleEncCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		consoleLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= consoleLevel
		})
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncCfg)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), consoleLevelEnabler))
	}

	var (
		syncer zapcore.WriteSyncer
		closer io.Closer
	)
	if !opts.DisableFile {
		var fileLevel Level
		if opts.FileLevel == "" {
			opts.FileLevel = InfoLevel.String()
		}
		err := fileLevel.Set(strings.ToLower(opts.FileLevel))
		if err != nil {
			fileLevel = InfoLevel
		}
		// Add level key for file log by default
		encoderConfig.LevelKey = "level"
		if !opts.DisableFileTime {
			encoderConfig.TimeKey = "time"
		}
		if !opts.DisableFileCaller {
			encoderConfig.CallerKey = "caller"
		}
		fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		if !opts.DisableFileJson {
			fileEncoder = zapcore.NewJSONEncoder(encoderConfig)
		}

		fileLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= fileLevel
		})
		syncer, closer = rollingFileEncoder(opts)
		cores = append(cores, zapcore.NewCore(fileEncoder, syncer, fileLevelEnabler))
	}
	core := zapcore.NewTee(cores...)
	unsugared := zap.New(core)
	return &Logger{
		log:     unsugared,
		sugared: unsugared.Sugar(),
		closer:  closer,
	}
}

func (l *Logger) DebugLogger() DebugLogger {
	return l
}

func (l *Logger) InfoLogger() InfoLogger {
	return l
}

func (l *Logger) Debugt(msg string, fields ...Field) {
	l.log.Debug(msg, fields...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugared.Debugf(template, args...)
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.sugared.Debugw(msg, keysAndValues...)
}

func (l *Logger) Infot(msg string, fields ...Field) {
	l.log.Info(msg, fields...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugared.Infof(template, args...)
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.sugared.Infow(msg, keysAndValues...)
}

func (l *Logger) Warnt(msg string, fields ...Field) {
	l.log.Warn(msg, fields...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugared.Warnf(template, args...)
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.sugared.Warnw(msg, keysAndValues...)
}

func (l *Logger) Errort(msg string, fields ...Field) {
	l.log.Error(msg, fields...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugared.Errorf(template, args...)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.sugared.Errorw(msg, keysAndValues...)
}

func (l *Logger) Panict(msg string, fields ...Field) {
	l.log.Panic(msg, fields...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.sugared.Panicf(template, args...)
}

func (l *Logger) Panic(msg string, keysAndValues ...interface{}) {
	l.sugared.Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatalt(msg string, fields ...Field) {
	l.log.Fatal(msg, fields...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugared.Fatalf(template, args...)
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.sugared.Fatalw(msg, keysAndValues...)
}

func (l *Logger) AtLevelt(level Level, msg string, fields ...Field) {
	switch level {
	case DebugLevel:
		l.Debugt(msg, fields...)
	case PanicLevel:
		l.Panict(msg, fields...)
	case ErrorLevel:
		l.Errort(msg, fields...)
	case WarnLevel:
		l.Warnt(msg, fields...)
	case InfoLevel:
		l.Infot(msg, fields...)
	case FatalLevel:
		l.Fatalt(msg, fields...)
	default:
		l.Warnt("unknown level", Any("level", level))
		l.Warnt(msg, fields...)
	}
}

func (l *Logger) AtLevel(level Level, msg string, keysAndValues ...interface{}) {
	switch level {
	case DebugLevel:
		l.Debug(msg, keysAndValues...)
	case PanicLevel:
		l.Panic(msg, keysAndValues...)
	case ErrorLevel:
		l.Error(msg, keysAndValues...)
	case WarnLevel:
		l.Warn(msg, keysAndValues...)
	case InfoLevel:
		l.Info(msg, keysAndValues...)
	case FatalLevel:
		l.Fatal(msg, keysAndValues...)
	default:
		l.Warnt("unknown level", Any("level", level))
		l.Warn(msg, keysAndValues...)
	}
}

func (l *Logger) AtLevelf(level Level, msg string, args ...interface{}) {
	switch level {
	case DebugLevel:
		l.Debugf(msg, args...)
	case PanicLevel:
		l.Panicf(msg, args...)
	case ErrorLevel:
		l.Errorf(msg, args...)
	case WarnLevel:
		l.Warnf(msg, args...)
	case InfoLevel:
		l.Infof(msg, args...)
	case FatalLevel:
		l.Fatalf(msg, args...)
	default:
		l.Warnt("unknown level", Any("level", level))
		l.Warnf(msg, args...)
	}
}

func (l *Logger) WithValues(fields ...Field) *Logger {
	newl := l.log.With(fields...)
	return &Logger{
		log:     newl,
		sugared: newl.Sugar(),
	}
}

func (l *Logger) Flush() error {
	return l.log.Sync()
}

func (l *Logger) Close() error {
	// https://github.com/uber-go/zap/issues/772
	_ = l.Flush()

	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}
