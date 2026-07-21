//go:build windows

package cmd

import (
	"os/exec"
	"strconv"
)

func killProcessTree(pid int) error {
	cmd := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}