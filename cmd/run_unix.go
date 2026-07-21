//go:build linux || darwin || freebsd || openbsd || netbsd

package cmd

import "syscall"

func killProcessTree(pid int) error {
	return syscall.Kill(-pid, syscall.SIGKILL)
}