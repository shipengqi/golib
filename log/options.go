package log

// Options Configuration for logging
type Options struct {
	// ConsoleEnabled the max size in MB of the logfile before it's rolled
	ConsoleEnabled      bool `json:"console-enabled" mapstructure:"console-enabled"`
	// ConsoleColorEnabled the max size in MB of the logfile before it's rolled
	ConsoleColorEnabled bool `json:"console-color-enabled" mapstructure:"console-color-enabled"`
	// ConsoleJson the max size in MB of the logfile before it's rolled
	ConsoleJson         bool `json:"console-json" mapstructure:"console-json"`
	// Filejson the max size in MB of the logfile before it's rolled
	FileJson            bool `json:"file-json" mapstructure:"file-json"`
	// FileEnabled the max size in MB of the logfile before it's rolled
	FileEnabled         bool `json:"level" mapstructure:"level"`
	// CallerEnabled the max size in MB of the logfile before it's rolled
	CallerEnabled       bool `json:"level" mapstructure:"level"`
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `json:"level" mapstructure:"level"`
	// MaxBackups the max number of rolled files to keep
	MaxBackups int `json:"level" mapstructure:"level"`
	// MaxAge the max age in days to keep a logfile
	MaxAge       int `json:"level" mapstructure:"level"`
	// CallerEnabled the max size in MB of the logfile before it's rolled
	ConsoleLevel string `json:"level" mapstructure:"level"`
	FileLevel    string `json:"level" mapstructure:"level"`
	// Directory to log when FileEnabled is true
	Directory string `json:"level" mapstructure:"level"`
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string `json:"level" mapstructure:"level"`
}

// NewOptions creates a Options with default parameters.
func NewOptions() *Options {
	return &Options{}
}