package sysutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Kill kills process by pid
func Kill(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}

// PID get process ID
func PID() int {
	return os.Getpid()
}

// GetProcessByPid 根据PID获取进程的执行路径.
func GetProcessByPid(pid int) string {
	exe := fmt.Sprintf("/proc/%d/exe", pid)
	path, _ := os.Readlink(exe)
	return path
}

func fqdn() (name string, err error) {
	cmd := exec.Command("hostname", "-f")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return
	}
	name = out.String()
	// removing EOL
	name = name[:len(name)-1]
	return
}
