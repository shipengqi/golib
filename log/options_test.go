package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Default_Options_Validate(t *testing.T) {
	opts := NewOptions()
	errs := opts.Validate()
	expected := 0

	assert.Equal(t, expected, len(errs))
}

func TestOptions_Validate(t *testing.T) {
	opts := NewOptions()
	opts.DisableConsole = true
	opts.DisableFile = true
	errs := opts.Validate()

	expected := `[no enabled logger]`
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))

	opts.DisableFile = false
	opts.Output = ""
	errs = opts.Validate()

	expected = `[no log output]`
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))

	opts.ConsoleLevel = "failed"
	opts.FileLevel = "failed"
	errs = opts.Validate()

	expected = `[unrecognized level: "failed" unrecognized level: "failed" no log output]`
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))
}
