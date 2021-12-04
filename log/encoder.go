package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// FilenameEncoder log filename encoder,
// return the fullname of the log file.
type FilenameEncoder func() string

// DefaultFilenameEncoder return <processname>-<date>.log.
func DefaultFilenameEncoder() string {
	return fmt.Sprintf("%s-%s.log", filepath.Base(os.Args[0]), time.Now().Format("20060102"))
}

func DefaultTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func rollingfileEncoder(opts *Options) zapcore.WriteSyncer {
	encoded := opts.FilenameEncoder()
	f := filepath.Join(opts.Output, encoded)
	if opts.DisableRotate {
		fd, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0644)
		if err != nil {
			panic(err)
		}
		return zapcore.AddSync(fd)
	}

	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   f,
		MaxSize:    opts.MaxSize,
		MaxAge:     opts.MaxAge,
		MaxBackups: opts.MaxBackups,
	})
}
