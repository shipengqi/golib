package log

import "errors"

// Options Configuration for logging
type Options struct {
	// DisableConsole whether to log to console
	DisableConsole bool `json:"disable-console" mapstructure:"disable-console"`
	// DisableConsoleColor force disabling colors.
	DisableConsoleColor bool `json:"disable-console-color" mapstructure:"disable-console-color"`
	// DisableFile whether to log to file
	DisableFile bool `json:"disable-file" mapstructure:"disable-file"`

	// DisableFileJson whether to enable json format for log file
	DisableFileJson bool `json:"disable-file-json" mapstructure:"disable-file-json"`

	// DisableCaller whether to log caller info
	DisableCaller bool `json:"disable-caller" mapstructure:"disable-caller"`

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

	// Output directory for logging when FileEnabled is true
	Output string `json:"output" mapstructure:"output"`
	// FilenameEncoder log filename encoder
	FilenameEncoder FilenameEncoder `json:"filename-encoder" mapstructure:"filename-encoder"`
}

// NewOptions creates a Options with default parameters.
func NewOptions() *Options {
	return &Options{
		DisableFile:     true,
		DisableRotate:   true,
		ConsoleLevel:    InfoLevel.String(),
		FilenameEncoder: DefaultFilenameEncoder,
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