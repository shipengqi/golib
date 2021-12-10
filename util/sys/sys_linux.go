package sys

import (
	"bytes"
	"os/exec"
)

func fqdn() (name string, err error) {
	cmd := exec.Command("hostname", "-f")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return
	}
	name = out.String()
	name = name[:len(name)-1] // removing EOL
	return
}
