package log

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	r, w, _ := os.Pipe()
	tmp := os.Stdout
	defer func() {
		os.Stdout = tmp
	}()
	os.Stdout = w

	name := "world"
	str := "Hello, world!"
	opts := NewOptions()
	Configure(opts)
	Debugf("Hello, %s", name+"1")
	Infof("Hello, %s", name+"2")
	Warnf("Hello, %s", name+"3")
	Errorf("Hello, %s", name+"4")

	opts.ConsoleLevel = DebugLevel.String()
	Configure(opts)
	Debug(str)
	Info(str)
	Warn(str)
	Error(str)

	opts.DisableConsoleColor = true
	Configure(opts)
	Debug(str)
	Info(str)
	Warn(str)
	Error(str)

	expected := []string{
		"\x1b[34mINFO\x1b[0m\tHello, world2",
		"\x1b[33mWARN\x1b[0m\tHello, world3",
		"\x1b[31mERROR\x1b[0m\tHello, world4",
		"\x1b[35mDEBUG\x1b[0m\tHello, world!",
		"\x1b[34mINFO\x1b[0m\tHello, world!",
		"\x1b[33mWARN\x1b[0m\tHello, world!",
		"\x1b[31mERROR\x1b[0m\tHello, world!",
		"debug\tHello, world!",
		"info\tHello, world!",
		"warn\tHello, world!",
		"error\tHello, world!",
	}
	_ = w.Close()
	stdout, _ := ioutil.ReadAll(r)
	reader := bytes.NewReader(stdout)
	scanner := bufio.NewScanner(reader)
	for _, v := range expected {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		assert.Contains(t, line, v)
	}
}

func TestLoggerPanic(t *testing.T) {
	str := "test panic"
	opts := NewOptions()
	Configure(opts)
	defer func() {
		if err := recover(); err != nil {
			assert.Equal(t, err, str)
		} else {
			t.Fatal("no panic")
		}
	}()
	Panic(str)
}

func TestWithValues(t *testing.T) {
	r, w, _ := os.Pipe()
	tmp := os.Stdout
	defer func() {
		os.Stdout = tmp
	}()
	os.Stdout = w
	opts := NewOptions()
	Configure(opts)

	logger := WithValues(String("test key", "test value"))
	logger.Info("Hello, world!")

	_ = w.Close()
	stdout, _ := ioutil.ReadAll(r)
	assert.Contains(t, string(stdout), "Hello, world!\t{\"test key\": \"test value\"}")
}

func TestDefaultLoggerWithoutTime(t *testing.T) {
	r, w, _ := os.Pipe()
	tmp := os.Stdout
	defer func() {
		os.Stdout = tmp
	}()
	os.Stdout = w
	opts := NewOptions()
	opts.DisableConsoleTime = true
	Configure(opts)

	Info("Hello, world!")
	_ = w.Close()
	stdout, _ := ioutil.ReadAll(r)
	assert.Equal(t, "\u001B[34mINFO\u001B[0m\tHello, world!\n", string(stdout))
}
