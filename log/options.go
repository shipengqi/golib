package log

import (
	"errors"
)

// Options Configuration for logging.
type Options struct {
	// DisableConsole whether to log to console
	DisableConsole bool `json:"disable-console" mapstructure:"disable-console"`
	// DisableConsoleColor force disabling colors.
	DisableConsoleColor bool `json:"disable-console-color" mapstructure:"disable-console-color"`
	// DisableConsoleTime whether to add a time
	DisableConsoleTime bool `json:"disable-console-time" mapstructure:"disable-console-time"`
	// DisableConsoleLevel whether to add a level
	DisableConsoleLevel bool `json:"disable-console-level" mapstructure:"disable-console-level"`
	// DisableConsoleCaller whether to log caller info
	DisableConsoleCaller bool `json:"disable-console-caller" mapstructure:"disable-console-caller"`

	// DisableFile whether to log to file
	DisableFile bool `json:"disable-file" mapstructure:"disable-file"`
	// DisableFileJson whether to enable json format for log file
	DisableFileJson bool `json:"disable-file-json" mapstructure:"disable-file-json"`
	// DisableFileTime whether to add a time
	DisableFileTime bool `json:"disable-file-time" mapstructure:"disable-file-time"`
	// DisableFileCaller whether to log caller info
	DisableFileCaller bool `json:"disable-file-caller" mapstructure:"disable-file-caller"`

	// DisableRotate whether to enable log file rotate
	DisableRotate bool `json:"disable-rotate" mapstructure:"disable-rotate"`
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `json:"max-size" mapstructure:"max-size"`
	// MaxBackups the max number of rolled files to keep
	MaxBackups int `json:"max-backups" mapstructure:"max-backups"`
	// MaxAge the max age in days to keep a logfile
	MaxAge int `json:"max-age" mapstructure:"max-age"`

	// ConsoleLevel sets the standard logger level
	ConsoleLevel string `json:"console-level" mapstructure:"console-level"`
	// FileLevel sets the file logger level.
	FileLevel string `json:"file-level" mapstructure:"file-level"`

	// Output directory for logging when DisableFile is false
	Output string `json:"output" mapstructure:"output"`
	// FilenameEncoder log filename encoder
	FilenameEncoder FilenameEncoder
	// TimeEncoder time encoder
	TimeEncoder TimeEncoder
}

// NewOptions creates an Options with default parameters.
func NewOptions() *Options {
	return &Options{
		DisableFile:          true,
		DisableRotate:        true,
		DisableConsoleCaller: true,
		DisableFileCaller:    true,
		ConsoleLevel:         InfoLevel.String(),
		FilenameEncoder:      DefaultFilenameEncoder,
		TimeEncoder:          DefaultTimeEncoder,
	}
}

// Validate validates the options fields.
func (o *Options) Validate() []error {
	var errs []error
	var level Level

	if o.ConsoleLevel != "" {
		if err := level.UnmarshalText([]byte(o.ConsoleLevel)); err != nil {
			errs = append(errs, err)
		}
	}

	if o.FileLevel != "" {
		if err := level.UnmarshalText([]byte(o.FileLevel)); err != nil {
			errs = append(errs, err)
		}
	}

	if o.DisableConsole && o.DisableFile {
		errs = append(errs, errors.New("no enabled logger"))
	}

	if !o.DisableFile && o.Output == "" {
		errs = append(errs, errors.New("no log output"))
	}
	return errs
}
