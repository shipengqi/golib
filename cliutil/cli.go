package cliutil

import (
	"os"
	"strings"
)

const defaultPlaceholder = "<>"

// RetrieveFlagFromCLI returns value of the given flag from os.Args.
func RetrieveFlagFromCLI(long string, short string) (value string, ok bool) {
	args := os.Args[1:]
	if len(args) < 1 {
		return
	}
	if len(short) == 0 {
		short = defaultPlaceholder // placeholder
	}
	return RetrieveFlag(args, long, short)
}

// RetrieveFlag returns value of the given flag from args.
func RetrieveFlag(args []string, long, short string) (value string, ok bool) {
	var index int
	for k := range args {
		if args[k] == long || args[k] == short {
			index = k + 1
			ok = true
			break
		}
	}
	if index == 0 {
		return
	}
	if len(args) < index+1 {
		return
	}
	if strings.HasPrefix(args[index], "-") {
		return "", ok
	}
	return args[index], ok
}
