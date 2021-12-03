package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Default_Logger(t *testing.T) {
	tmp := os.Stdout
	defer func() {
		os.Stdout = tmp
	}()

	str := "Hello, world!"
	tmpl := "Hello, %s"
	printItems := []printItem{
		{level: DebugLevel, tmpl: str, expected: ""},
		{level: InfoLevel, tmpl: str, expected: str},
		{level: WarnLevel, tmpl: str, expected: str},
		{level: ErrorLevel, tmpl: str, expected: str},
	}
	printfItems := []printItem{
		{level: DebugLevel, tmpl: tmpl, expected: "Hello, test1", value: []string{"test1"}},
		{level: InfoLevel, tmpl: tmpl, expected: "Hello, test2", value: []string{"test2"}},
		{level: WarnLevel, tmpl: tmpl, expected: "Hello, test3", value: []string{"test3"}},
		{level: ErrorLevel, tmpl: tmpl, expected: "Hello, test4", value: []string{"test4"}},
	}
	opts := NewOptions()
	Configure(opts)
	testPrintf(printfItems, t)

	opts.ConsoleLevel = DebugLevel.String()
	Configure(opts)
	testPrint(printItems, t)

	opts.DisableConsoleColor = true
	Configure(opts)
	testPrint(printItems, t)
}

type printItem struct {
	level    Level
	expected string
	tmpl     string
	value    []string
}

func testPrintf(items []printItem, t *testing.T) {
	var r, w *os.File
	for _, v := range items {
		r, w, _ = os.Pipe()
		os.Stdout = w
		AtLevelf(v.level, v.tmpl, v.value)
		stdout, _ := ioutil.ReadAll(r)
		assert.Contains(t, v.expected, string(stdout))
	}
}

func testPrint(items []printItem, t *testing.T) {
	var r, w *os.File
	for _, v := range items {
		r, w, _ = os.Pipe()
		os.Stdout = w
		AtLevel(v.level, v.tmpl, v.value)
		stdout, _ := ioutil.ReadAll(r)
		assert.Contains(t, v.expected, string(stdout))
		_ = w.Close()
	}
}
