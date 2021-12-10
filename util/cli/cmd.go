package cli

import "strings"

// Cmd represents an external command being prepared or run.
//
// A Cmd cannot be reused after calling its Run, Output or CombinedOutput
// methods.
type Cmd struct {
	// Path is the path of the command to run.
	//
	// This is the only field that must be set to a non-zero
	// value. If Path is relative, it is evaluated relative
	// to Dir.
	Path string

	// Args holds command line arguments.
	// If the Args field is empty or nil, Run uses {Path}.
	//
	// In typical use, both Path and Args are set by calling Command.
	Args []string
}

func (c *Cmd) String() string {
	b := new(strings.Builder)
	b.WriteString(c.Path)
	for _, a := range c.Args {
		b.WriteByte(' ')
		b.WriteString(a)
	}
	return b.String()
}
