package container

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	EnvExecPid = "_CONTAINER_PID"
	EnvExecCmd = "_CONTAINER_CMD"
)

func ExecInContainer(containerId, command string) (string, string, error) {
	return "", "", nil
}

func GetPidWithContainerId(containerId string) (string, error) {
	return "", nil
}

func GetEnvWithPid(pid string) ([]string, error) {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	envs := strings.Split(string(content), "\u0000")
	return envs, nil
}
